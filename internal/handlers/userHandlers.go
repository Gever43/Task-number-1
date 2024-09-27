package handlers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"myProject/internal/userService"
	"myProject/internal/web/users"

	"gorm.io/gorm"
)

// UserHandler для работы с пользователями
type UserHandler struct {
    Service *userService.UserService
}

// NewUserHandler создает новый экземпляр UserHandler
func NewUserHandler(service *userService.UserService) *UserHandler {
    return &UserHandler{
        Service: service,
    }
}

// Обработчик получения пользователей
func (h *UserHandler) GetUsers(ctx context.Context, request users.GetUsersRequestObject) (users.GetUsersResponseObject, error) {
    allUsers, err := h.Service.GetAllUsers()
    if err != nil {
        return nil, err
    }

    response := users.GetUsers200JSONResponse{}

    // Заполняем ответ
    for _, user := range allUsers {
        userResponse := users.User{
            Id:       &user.ID,
            Name:     user.Name,
            Email:    user.Email,
            Password: user.Password, 
        }
        response = append(response, userResponse)
    }

    return response, nil
}

// Обработчик создания пользователя
func (h *UserHandler) PostUsers(ctx context.Context, request users.PostUsersRequestObject) (users.PostUsersResponseObject, error) {
    if request.Body == nil {
        return nil, errors.New("body cannot be nil")
    }

    userToCreate := userService.DBUser{
        Name:     request.Body.Name,
        Email:    request.Body.Email,
        Password: request.Body.Password,
    }

    createdUser, err := h.Service.CreateUser(userToCreate.Name, userToCreate.Email, userToCreate.Password)
    if err != nil {
        return nil, err
    }

    response := users.PostUsers201JSONResponse{
        Id:       &createdUser.ID,
        Name:     createdUser.Name,
        Email:    createdUser.Email,
    }

    return response, nil
}

// Обработчик обновления пользователя
func (h *UserHandler) PatchUsers(ctx context.Context, request users.PatchUsersIdRequestObject) (users.PatchUsersId200JSONResponse, error) {
    log.Printf("Received request to update user with ID: %d", request.Body.Id)
    log.Printf("Request body: %+v", request.Body)

    // Проверяем, что Body не nil и Name не пустое
    if request.Body == nil || request.Body.Name == "" {
        return users.PatchUsersId200JSONResponse{}, errors.New("body or name cannot be empty")
    }

    // Преобразование int в uint
    id := uint(*request.Body.Id)

    // Создание объекта пользователя для обновления
    userToUpdate := userService.DBUser{
        Name:     request.Body.Name,
        Email:    request.Body.Email,
        Password: request.Body.Password,
    }

    // Обновление пользователя
    updatedUser, err := h.Service.UpdateUserByID(id, userToUpdate)
    if err != nil {
        log.Printf("Error updating user: %v", err)
        return users.PatchUsersId200JSONResponse{}, fmt.Errorf("failed to update user: %w", err)
    }

    // Создаем структуру респонс
    response := users.PatchUsersId200JSONResponse{
        Id:      &updatedUser.ID, 
        Name:    updatedUser.Name,
        Email:   updatedUser.Email,
    }

    return response, nil
}

// Обработчик удаления пользователя
func (h *UserHandler) DeleteUsers(ctx context.Context, request users.DeleteUsersIdRequestObject) (users.DeleteUsersIdResponseObject, error) {
    id := uint(request.Id) // Преобразование int в uint

    // Удаление пользователя
    err := h.Service.DeleteUserByID(id)
    if err != nil {
        // Проверяем, если ошибка связана с отсутствием пользователя
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, fmt.Errorf("user with id %d not found", id) // Сообщение об ошибке
        }
        return nil, fmt.Errorf("failed to delete user: %w", err) // Более информативное сообщение об ошибке
    }

    return users.DeleteUsersId204Response{}, nil
}