package service

import (
	"finfluence-chat/internal/model"
	"sync"
)

type HistoryStore interface {
	Get(userID string) []model.Message
	Add(userID string, msg model.Message)
}

type InMemoryStore struct {
	HistoryMap map[string][]model.Message
	Capacity   int
	rwMutex    sync.RWMutex
}

func NewInMemoryStore(Capacity int) *InMemoryStore {
	return &InMemoryStore{
		HistoryMap: make(map[string][]model.Message),
		Capacity:   Capacity,
		rwMutex:    sync.RWMutex{},
	}
}

func (inMemoryStore *InMemoryStore) Get(userID string) []model.Message {
	inMemoryStore.rwMutex.RLock()
	defer inMemoryStore.rwMutex.RUnlock()
	return append([]model.Message(nil), inMemoryStore.HistoryMap[userID]...)
}

func (inMemoryStore *InMemoryStore) Add(userID string, msg model.Message) {
	inMemoryStore.rwMutex.Lock()
	defer inMemoryStore.rwMutex.Unlock()
	inMemoryStore.HistoryMap[userID] = append(inMemoryStore.HistoryMap[userID], msg)
	if len(inMemoryStore.HistoryMap[userID]) > inMemoryStore.Capacity {
		inMemoryStore.HistoryMap[userID] = inMemoryStore.HistoryMap[userID][len(inMemoryStore.HistoryMap[userID])-inMemoryStore.Capacity:]
	}
}
