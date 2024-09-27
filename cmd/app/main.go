package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"myProject/internal/database"
	"myProject/internal/handlers"
	"myProject/internal/messagesService"
	"myProject/internal/userService"
	"myProject/internal/web/messages"
	"myProject/internal/web/users"
)

func main() {
    database.InitDB()

    // Проверяем ошибку при миграции для сообщений
    if err := database.DB.AutoMigrate(&messagesService.DBMessage{}); err != nil {
        log.Fatalf("failed to migrate database for messages: %v", err)
    } else {
        log.Println("Messages table migrated successfully")
    }
    
    if err := database.DB.AutoMigrate(&userService.DBUser{}); err != nil {
        log.Fatalf("failed to migrate database for users: %v", err)
    } else {
        log.Println("Users table migrated successfully")
    }

    // Разделение логики на слои для сообщений
    messagesRepo := messagesService.NewMessageRepository(database.DB)
    messagesService := messagesService.NewService(messagesRepo)
    messagesHandler := handlers.NewHandler(messagesService)

    // Разделение логики на слои для пользователей
    userRepo := userService.NewUserRepository(database.DB)
    userService := userService.NewUserService(userRepo)
    userHandler := handlers.NewUserHandler(userService)

    // Инициализируем echo - фреймворк
    e := echo.New()
    e.Use(middleware.Logger())//логирование (отслеживание работы приложения)
    e.Use(middleware.Recover())//восстановление после ошибок

    // Обработчики для сообщений
    e.POST("/api/messages", func(c echo.Context) error {
        var request messages.PostMessagesRequestObject

        // Десериализация данных из тела запроса к структуре request
        if err := c.Bind(&request); err != nil {
            return c.String(http.StatusBadRequest, err.Error())
        }

        // Проверяем, что Body не nil
        if request.Body == nil || request.Body.Message == nil {
            return c.String(http.StatusBadRequest, "body or message cannot be nil")
        }

        // Создание сообщения
        response, err := messagesHandler.PostMessages(c.Request().Context(), request)
        if err != nil {
            return c.String(http.StatusInternalServerError, err.Error())
        }

        // Если всё успешно, то выдаём JSON клиенту
        return c.JSON(http.StatusCreated, response)
    })

    e.GET("/api/messages", func(c echo.Context) error {
        var request messages.GetMessagesRequestObject

        // Извлечение сообщений из БД
        response, err := messagesHandler.GetMessages(c.Request().Context(), request)
        if err != nil {
            return c.String(http.StatusInternalServerError, err.Error())
        }

        // Если всё успешно, то выдаём JSON клиенту
        return c.JSON(http.StatusOK, response)
    })

    // PATCH запрос для обновления сообщения
    e.PATCH("/api/messages/:id", func(c echo.Context) error {
        messageID := c.Param("id") // Получаем ID сообщения из URL
        var request messages.PatchMessagesRequestObject

        // Десериализация данных из тела запроса
        if err := c.Bind(&request); err != nil {
            return c.String(http.StatusBadRequest, err.Error())
        }

        // Преобразуем messageID из строки в uint
        id, err := strconv.ParseUint(messageID, 10, 32)
        if err != nil {
            return c.String(http.StatusBadRequest, "invalid message ID")
        }

        // Устанавливаем ID в Body запроса
        request.Body.Id = new(uint) // Создаем указатель на uint
        *request.Body.Id = uint(id)  // Присваиваем значение

        // Обновление сообщения
        response, err := messagesHandler.PatchMessage(c.Request().Context(), request)
        if err != nil {
            return c.String(http.StatusInternalServerError, err.Error())
        }

        return c.JSON(http.StatusOK, response)
    })

    e.DELETE("/api/messages/:id", func(c echo.Context) error {
        messageID := c.Param("id") // Получаем ID сообщения из URL

        // Преобразуем messageID из строки в uint
        id, err := strconv.ParseUint(messageID, 10, 32)
        if err != nil {
            return c.String(http.StatusBadRequest, "invalid message ID")
        }

        // Создаем объект запроса
        request := messages.DeleteMessagesRequestObject{
            Id: uint(id), // Присваиваем значение ID
        }

        // Вызов функции для удаления сообщения
        _ , err = messagesHandler.DeleteMessages(c.Request().Context(), request)
        if err != nil {
            return c.String(http.StatusInternalServerError, err.Error())
        }

        return c.String(http.StatusOK, "Запись успешно удалена!")
    })

    // Обработчики для пользователей
    e.POST("/api/users", func(c echo.Context) error {
        var userRequest users.CreateUserRequest
        if err := c.Bind(&userRequest); err != nil {
            return c.String(http.StatusBadRequest, err.Error())
        }
    
        // Проверяем, что тело запроса не нулевое
        if userRequest.Body == nil || userRequest.Body.Name == "" || userRequest.Body.Email == "" || userRequest.Body.Password == "" {
            return c.String(http.StatusBadRequest, "body, name, email, or password cannot be nil or empty")
        }
    
        // Создаем объект PostUsersRequestObject
        postUsersRequest := users.PostUsersRequestObject{
            Body: userRequest.Body, // Передаем тело из userRequest
        }
    
        // Вызываем метод PostUsers с правильным типом аргумента
        newUserResponse, err := userHandler.PostUsers(c.Request().Context(), postUsersRequest)
        if err != nil {
            return c.String(http.StatusInternalServerError, err.Error())
        }
        return c.JSON(http.StatusCreated, newUserResponse)
    })

    e.GET("/api/users", func(c echo.Context) error {
        var request users.GetUsersRequestObject

        // Вызываем метод GetUsers с правильными аргументами
        response, err := userHandler.GetUsers(c.Request().Context(), request)
        if err != nil {
            return c.String(http.StatusInternalServerError, err.Error())
        }

        // Возвращаем успешный ответ
        return response.VisitGetUsersResponse(c.Response())
    })
    // PATCH запрос для обновления пользователя
// PATCH запрос для обновления пользователя
e.PATCH("/api/users/:id", func(c echo.Context) error {
    userID := c.Param("id") // Получаем ID пользователя из URL
    var request users.PatchUsersRequestObject // Создаем объект запроса

    // Привязываем тело запроса
    if err := c.Bind(&request); err != nil {
        return c.String(http.StatusBadRequest, err.Error())
    }

    // Преобразуем userID из строки в uint
    id, err := strconv.ParseUint(userID, 10, 32)
    if err != nil {
        return c.String(http.StatusBadRequest, "invalid user ID")
    }

     // Устанавливаем ID в Body запроса
     request.Body.ID = uint(id) // Присваиваем значение напрямую

    // Вызываем метод PatchUser для обновления пользователя
    updatedUserResponse, err := userHandler.PatchUsers(c.Request().Context(), request)
    if err != nil {
        return c.String(http.StatusInternalServerError, err.Error())
    }

    // Возвращаем успешный ответ
    return c.JSON(http.StatusOK, updatedUserResponse)
})

// DELETE запрос для удаления пользователя
e.DELETE("/api/users/:id", func(c echo.Context) error {
    userID := c.Param("id") // Получаем ID пользователя из URL

    // Преобразуем userID из строки в uint
    id, err := strconv.ParseUint(userID, 10, 32)
    if err != nil {
        return c.String(http.StatusBadRequest, "invalid user ID")
    }

    // Создаем объект запроса
    request := users.DeleteUsersRequestObject{
        Id: uint(id), // Присваиваем значение ID
    }

    // Вызываем метод DeleteUser для удаления пользователя
    if _ , err := userHandler.DeleteUsers(c.Request().Context(), request); err != nil {
        return c.String(http.StatusInternalServerError, err.Error())
    }

    return c.String(http.StatusOK, "Запись успешно удалена!")
})

// Запуск сервера на порту 8080
if err := e.Start(":8080"); err != nil {
    log.Fatalf("failed to start server: %v", err)
 }

}
    