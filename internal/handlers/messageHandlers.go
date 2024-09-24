package handlers

import (
	"context"
	"errors"
	"myProject/internal/messagesService"
	"myProject/internal/web/messages"
)

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
        message := messages.Message{
            Id:      &msg.ID,
            Message: &msg.Text,
        }
        response = append(response, message)
    }

    // Возвращаем респонс и nil!
    return response, nil
}

func (h *Handler) PostMessages(_ context.Context, request messages.PostMessagesRequestObject) (messages.PostMessagesResponseObject, error) {
    // Распаковываем тело запроса напрямую, без декодера!
    messageRequest := request.Body
    // Создаем сообщение
    messageToCreate := messagesService.Message{Text: *messageRequest.Message}
    createdMessage, err := h.Service.CreateMessage(messageToCreate)

    if err != nil {
        return nil, err
    }
    // Создаем структуру респонс
    response := messages.PostMessages201JSONResponse{
        Id:      &createdMessage.ID,
        Message: &createdMessage.Text,
    }
    // Возвращаем респонс!
    return response, nil
}


// Обработчик PATCH запроса для обновления сообщения
func (h *Handler) PatchMessage(_ context.Context, request messages.PatchMessagesRequestObject) (messages.PatchMessagesResponseObject, error) {
    // Проверяем, что Body не nil и Message не nil
    if request.Body == nil || request.Body.Message == nil {
        return nil, errors.New("body or message cannot be nil")
    }

    // Обновляем сообщение, используя метод из MessageService
    updatedMessage, err := h.Service.UpdateMessageByID(int(*request.Body.Id), messagesService.Message{Text: *request.Body.Message}) 
    if err != nil {
        return nil, err
    }

    // Создаем структуру респонс
    response := messages.PatchMessages200JSONResponse{
        Id:      &updatedMessage.ID,
        Message: &updatedMessage.Text,
    }
    // Возвращаем респонс и nil
    return response, nil
}

// Обработчик DELETE запроса для удаления сообщения
func (h *Handler) DeleteMessages(_ context.Context, request messages.DeleteMessagesRequestObject) (messages.DeleteMessagesResponseObject, error) {
    err := h.Service.DeleteMessageByID(int(request.Id))
    if err != nil {
        return nil, err
    }
    // Возвращаем пустой ответ, так как сообщение успешно удалено
    return messages.DeleteMessages204Response{}, nil
}

