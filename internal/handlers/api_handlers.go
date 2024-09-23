package handlers

import (
	"encoding/json"
	"myProject/internal/messagesService"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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

// GetMessagesHandler - обработчик для получения всех сообщений
func (h *Handler) GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
    messages, err := h.Service.GetAllMessages()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(messages)
}

// PostMessageHandler - обработчик для создания нового сообщения
func (h *Handler) PostMessageHandler(w http.ResponseWriter, r *http.Request) {
    var message messagesService.Message
    err := json.NewDecoder(r.Body).Decode(&message)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    createdMessage, err := h.Service.CreateMessage(message)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(createdMessage)
}

// PatchMessageHandler - обработчик для обновления сообщения по ID
func (h *Handler) PatchMessageHandler(w http.ResponseWriter, r *http.Request) {
    // Извлечение ID из URL
    vars := mux.Vars(r)
    idStr := vars["id"]

    // Преобразуем ID из строки в целое число
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    var message messagesService.Message
    err = json.NewDecoder(r.Body).Decode(&message)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    updatedMessage, err := h.Service.UpdateMessageByID(id, message)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(updatedMessage)
}

// DeleteMessageHandler - обработчик для удаления сообщения по ID
func (h *Handler) DeleteMessageHandler(w http.ResponseWriter, r *http.Request) {
    // Извлечение ID из URL
    vars := mux.Vars(r)
    idStr := vars["id"]

    // Преобразуем ID из строки в целое число
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    // Вызываем метод сервиса для удаления сообщения
    err = h.Service.DeleteMessageByID(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent) // Возвращаем статус 204 No Content при успешном удалении
}
