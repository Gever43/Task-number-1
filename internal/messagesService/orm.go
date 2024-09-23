package messagesService

import (
	"gorm.io/gorm"
)

// Определяем структуру Message для работы с БД
type Message struct {
    gorm.Model // Встраиваем модель GORM для автоматического добавления ID, Timestamp и других полей
    Text       string `json:"text"` // Наш сервер будет ожидать JSON с полем text
}