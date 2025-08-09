package websocket

import (
	"smart-dialog-ai/internal/model"
	"smart-dialog-ai/internal/repository"
	"sync"

	"gorm.io/gorm"
)

type ChatHistoryManager struct {
	mu      sync.RWMutex
	history []model.Message
}

func NewChatHistoryManager() *ChatHistoryManager {
	return &ChatHistoryManager{
		history: make([]model.Message, 0),
	}
}

func (m *ChatHistoryManager) LoadHistory(db *gorm.DB , userID string) {
	// 这里你调用 repository.LoadHistory
	historyFromDB := repository.LoadHistory(db, userID)
	m.mu.Lock()
	m.history = historyFromDB
	m.mu.Unlock()
}

func (m *ChatHistoryManager) AppendMessage(msg model.Message) {
	m.mu.Lock()
	m.history = append(m.history, msg)
	m.mu.Unlock()
}

func (m *ChatHistoryManager) GetHistoryCopy() []model.Message {
	m.mu.RLock()
	defer m.mu.RUnlock()
	copied := make([]model.Message, len(m.history))
	copy(copied, m.history)
	return copied
}
