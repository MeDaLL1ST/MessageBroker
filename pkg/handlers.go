package pkg

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"selfmq/metrics"
	"sync"

	"github.com/gorilla/websocket"
)

var store = &Store{
	Uses:    make(map[string]int),
	Updates: make(map[string]chan Update),
	Lock:    &sync.Mutex{},
	RLock:   &sync.RWMutex{},
}

func AddHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) == 0 {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return
	}

	if !validateToken(authHeader) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	go metrics.Incr()
	var item Item
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	go store.RPush(item.Key, item.Value)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(item)

}

func ListHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) == 0 {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return
	}

	if !validateToken(authHeader) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	go metrics.IncrInf()
	// Получаем все ключи из мапы Uses
	keys := make([]string, 0, len(store.Uses))
	store.RLock.RLock()
	for k := range store.Uses {
		keys = append(keys, k)
	}
	store.RLock.RUnlock()

	// Создаем JSON-массив с ключами
	response := struct {
		Keys []string `json:"keys"`
	}{
		Keys: keys,
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func InfoHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) == 0 {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return
	}

	if !validateToken(authHeader) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	go metrics.IncrInf()
	var item SubscribeRequest
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	store.RLock.RLock()
	subs := store.Uses[item.Key]
	store.RLock.RUnlock()

	response := struct {
		Subs int `json:"subs"`
	}{
		Subs: subs,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

func SubscribeHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) == 0 {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return
	}

	if !validateToken(authHeader) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	go metrics.IncrC()
	defer metrics.DecrC()
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	var currentCancel context.CancelFunc

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				//log.Println("Connection closed by client")
			} else {
				log.Println("Read error:", err)
			}
			break
		}

		var req SubscribeRequest
		err = json.Unmarshal(message, &req)
		if err != nil {
			log.Printf("Error unmarshalling request: %v", err)
			continue
		}

		// Завершение предыдущей подписки, если она существует
		if currentCancel != nil {
			currentCancel()
		}

		// Создание нового контекста и новой подписки
		var ctx context.Context
		ctx, currentCancel = context.WithCancel(context.Background())
		go subscribeToKey(ctx, conn, req.Key)
	}

	// Закрытие текущей подписки при завершении соединения
	if currentCancel != nil {
		currentCancel()
	}
}

func subscribeToKey(ctx context.Context, conn *websocket.Conn, key string) {
	updates := store.GetUpdates(key)
	store.IncUses(key)
	for {
		select {
		case <-ctx.Done():
			store.DecUses(key)
			if store.Uses[key] == 0 {
				store.QClear(key)
			}
			return
		case update := <-updates:
			err := conn.WriteMessage(websocket.TextMessage, []byte(update.Value))
			if err != nil {
				return
			}
		}
	}
}
