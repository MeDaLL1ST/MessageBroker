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
	Lock    *sync.Mutex
	RLock   *sync.RWMutex
}

func (s *Store) GetUpdates(key string) <-chan Update {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	if _, ok := s.Updates[key]; !ok {
		s.Uses[key] = 0
		s.Updates[key] = make(chan Update)
	}
	return s.Updates[key]
}

func (s *Store) QClear(key string) {
	s.Lock.Lock()
	delete(s.Updates, key)
	delete(s.Uses, key)
	s.Lock.Unlock()
}

func (s *Store) RPush(key string, value string) {
	s.Lock.Lock()
	if updates, ok := s.Updates[key]; ok {
		updates <- Update{Key: key, Value: value}
	}
	s.Lock.Unlock()
}

func (s *Store) IncUses(key string) {
	s.Lock.Lock()
	s.Uses[key]++
	s.Lock.Unlock()
}

func (s *Store) DecUses(key string) {
	s.Lock.Lock()
	s.Uses[key]--
	s.Lock.Unlock()
}
