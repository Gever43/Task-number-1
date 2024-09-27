package handlers

import (
	"context"
	"errors"
	"fmt"
	"myProject/internal/userService"
	"myProject/internal/web/users"
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
            ID:       user.ID,
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
        ID:       createdUser.ID,
        Name:     createdUser.Name,
        Email:    createdUser.Email,
    }

    return response, nil
}

// Обработчик обновления пользователя
func (h *UserHandler) PatchUsers(ctx context.Context, request users.PatchUsersRequestObject) (users.PatchUsersResponseObject, error) {
    if request.Body == nil {
        return nil, errors.New("body cannot be nil")
    }

    // Преобразование int в uint
    id := uint(request.Id) 

    // Создание объекта пользователя для обновления
    userToUpdate := userService.DBUser{
        Name:     request.Body.Name,
        Email:    request.Body.Email,
        Password: request.Body.Password,
    }

    // Обновление пользователя
    updatedUser, err := h.Service.UpdateUserByID(id, userToUpdate)
    if err != nil {
        return nil, fmt.Errorf("failed to update user: %w", err) // Более информативное сообщение об ошибке
    }

    response := users.PatchUsers200JSONResponse{
        ID:       updatedUser.ID,
        Name:     updatedUser.Name,
        Email:    updatedUser.Email,
    }

    return response, nil
}

// Обработчик удаления пользователя
func (h *UserHandler) DeleteUsers(ctx context.Context, request users.DeleteUsersRequestObject) (users.DeleteUsersResponseObject, error) {
    id := uint(request.Id) // Преобразование int в uint

    // Удаление пользователя
    err := h.Service.DeleteUserByID(id)
    if err != nil {
        return nil, fmt.Errorf("failed to delete user: %w", err) // Более информативное сообщение об ошибке
    }

    return users.DeleteUsers204Response{}, nil
}