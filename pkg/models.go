package pkg

import "sync"

type Item struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type SubscribeRequest struct {
	Key string `json:"key"`
}

type Update struct {
	Key   string
	Value string
}

type Store struct {
	Uses    map[string]int
	Updates map[string]chan Update
	RLock   *sync.RWMutex
}

func (s *Store) GetUpdates(key string) <-chan Update {
	s.RLock.Lock()
	defer s.RLock.Unlock()

	if _, ok := s.Updates[key]; !ok {
		s.Uses[key] = 0
		s.Updates[key] = make(chan Update)
	}
	return s.Updates[key]
}

func (s *Store) QClear(key string) {
	s.RLock.Lock()
	delete(s.Updates, key)
	delete(s.Uses, key)
	s.RLock.Unlock()
}

func (s *Store) RPush(key string, value string) {
	s.RLock.RLock()
	if updates, ok := s.Updates[key]; ok {
		updates <- Update{Key: key, Value: value}
	}
	s.RLock.RUnlock()
}

func (s *Store) IncUses(key string) {
	s.RLock.Lock()
	s.Uses[key]++
	s.RLock.Unlock()
}

func (s *Store) DecUses(key string) {
	s.RLock.Lock()
	s.Uses[key]--
	s.RLock.Unlock()
}
