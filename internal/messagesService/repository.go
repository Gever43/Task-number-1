package messagesService

import (
	"gorm.io/gorm"
)

// MessageRepository - интерфейс для работы с сообщениями
type MessageRepository interface {

    // CreateMessage - Создает новое сообщение и возвращает его
    CreateMessage(message DBMessage ) (DBMessage , error)

    // GetAllMessages - Возвращает все сообщения из БД
    GetAllMessages() ([]DBMessage , error)

    // UpdateMessageByID - Обновляет сообщение по ID и возвращает его
    UpdateMessageByID(id int, message DBMessage ) ( DBMessage , error)
	
    // DeleteMessageByID - Удаляет сообщение по ID
    DeleteMessageByID(id int) error
}

// messageRepository - структура, поддерживающая интерфейс, т.к. поддерживает все методы интерфейса
type messageRepository struct {
    db *gorm.DB
}

// newMessageRepository - создает новый экземпляр messageRepository
func NewMessageRepository(db *gorm.DB) *messageRepository {
    return &messageRepository{db: db} //инициализация поля db структуры messageRepository перереданным аргументом db
}

// CreateMessage - реализация метода для создания сообщения
func (r *messageRepository) CreateMessage(message DBMessage ) (DBMessage , error) {
    result := r.db.Create(&message)
    if result.Error != nil {
        return DBMessage {}, result.Error
    }
    return message, nil
}

// GetAllMessages - реализация метода для получения всех сообщений
func (r *messageRepository) GetAllMessages() ([]DBMessage , error) {
    var messages []DBMessage 
    err := r.db.Find(&messages).Error
    return messages, err
}

// UpdateMessageByID - реализация метода для обновления сообщения по ID
func (r *messageRepository) UpdateMessageByID(id int, message DBMessage ) (DBMessage , error) {
    var existingMessage DBMessage 
    if err := r.db.First(&existingMessage, id).Error; err != nil {
        return DBMessage {}, err
    }

    // Обновляем поля существующего сообщения
    existingMessage.Text = message.Text
    if err := r.db.Save(&existingMessage).Error; err != nil {
        return DBMessage {}, err
    }

    return existingMessage, nil
}

// DeleteMessageByID - реализация метода для удаления сообщения по ID
func (r *messageRepository) DeleteMessageByID(id int) error {
    if err := r.db.Delete(&DBMessage {}, id).Error; err != nil {
        return err
    }
    return nil
}