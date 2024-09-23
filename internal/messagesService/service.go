package messagesService

type MessageService struct {
	repo MessageRepository
}

// NewService - создает новый экземпляр MessageService
func NewService(repo MessageRepository) *MessageService {
	return &MessageService{repo: repo}
}

// CreateMessage - создает новое сообщение
func (s *MessageService) CreateMessage(message Message) (Message, error) {
	return s.repo.CreateMessage(message)
}

// GetAllMessages - возвращает все сообщения
func (s *MessageService) GetAllMessages() ([]Message, error) {
	return s.repo.GetAllMessages()
}

// UpdateMessageByID - обновляет сообщение по ID
func (s *MessageService) UpdateMessageByID(id int, message Message) (Message, error) {
	return s.repo.UpdateMessageByID(id, message)
}

// DeleteMessageByID - удаляет сообщение по ID
func (s *MessageService) DeleteMessageByID(id int) error {
	return s.repo.DeleteMessageByID(id)
}