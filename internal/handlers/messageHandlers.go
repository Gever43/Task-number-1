package handlers

import (
	"context"
	"errors"
	"fmt"
	"myProject/internal/messagesService"
	"myProject/internal/web/messages"

	"gorm.io/gorm"
)

//создание структуры для "спуска" логики на уровень сервиса
type Handler struct {
    Service *messagesService.MessageService
}

// NewHandler создает новый экземпляр Handler
func NewHandler(service *messagesService.MessageService) *Handler {
    return &Handler{
        Service: service,
    }
}

func (h *Handler) GetMessages(_ context.Context, _ messages.GetMessagesRequestObject) (messages.GetMessagesResponseObject, error) {
    // Получение всех сообщений из сервиса
    allMessages, err := h.Service.GetAllMessages()
    if err != nil {
        return nil, err
    }

    // Создаем переменную респон типа 200 джейсон респонс
    response := messages.GetMessages200JSONResponse{}

    // Заполняем слайс response всеми сообщениями из БД
    for _, msg := range allMessages {
        message := messages.Message {
            Id:      &msg.ID,
            Message: msg.Text,
        }
        response = append(response, message)
    }

    // Возвращаем респонс и nil!
    return response, nil
}

func (h *Handler) PostMessages(_ context.Context, request messages.PostMessagesRequestObject) (messages.PostMessagesResponseObject, error) {
    // Проверяем, что Body не nil
    if request.Body == nil {
        return nil, errors.New("body cannot be nil")
    }

    // Распаковываем тело запроса напрямую, без декодера!
    messageRequest := request.Body
    // Создаем сообщение
    messageToCreate := messagesService.DBMessage {Text: messageRequest.Message}
    // Передаём на другой уровень
    createdMessage, err := h.Service.CreateMessage(messageToCreate)

    if err != nil {
        return nil, err
    }
    // Создаем структуру респонс
    response := messages.PostMessages201JSONResponse{
        Id:      &createdMessage.ID,
        Message: createdMessage.Text,
    }
    // Возвращаем респонс!
    return response, nil
}


func (h *Handler) PatchMessage(_ context.Context, request messages.PatchMessagesIdRequestObject) (messages.PatchMessagesId200JSONResponse, error) {
    // Проверяем, что Body не nil и Message не пустое
    if request.Body == nil || request.Body.Message == "" {
        return messages.PatchMessagesId200JSONResponse{}, errors.New("body or message cannot be empty")
    }

    // Обновляем сообщение, используя метод из MessageService
    updatedMessage, err := h.Service.UpdateMessageByID(int(*request.Body.Id), messagesService.DBMessage{Text: request.Body.Message})
    if err != nil {
        return messages.PatchMessagesId200JSONResponse{}, err
    }

    // Создаем структуру респонс
    response := messages.PatchMessagesId200JSONResponse{
        Id:      &updatedMessage.ID,
        Message: updatedMessage.Text,
    }
    
    // Возвращаем респонс и nil
    return response, nil
}

// Обработчик DELETE запроса для удаления сообщения
func (h *Handler) DeleteMessages(ctx context.Context, request messages.DeleteMessagesIdRequestObject) (messages.DeleteMessagesIdResponseObject, error) {
    // Удаление сообщения
    err := h.Service.DeleteMessageByID(int(request.Id))
    if err != nil {
        // Проверяем, если ошибка связана с отсутствием сообщения
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, fmt.Errorf("message with id %d not found", request.Id) // Сообщение об ошибке
        }
        return nil, fmt.Errorf("failed to delete message: %w", err) // Более информативное сообщение об ошибке
    }

    // Возвращаем пустой ответ, так как сообщение успешно удалено
    return messages.DeleteMessagesId204Response{}, nil
}

