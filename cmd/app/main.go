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
	"myProject/internal/web/messages"
)

func main() {
    database.InitDB()
    // Проверяем ошибку при миграции
    if err := database.DB.AutoMigrate(&messagesService.Message{}); err != nil {
        log.Fatalf("failed to migrate database: %v", err)
    }
    
    repo := messagesService.NewMessageRepository(database.DB)
    service := messagesService.NewService(repo)

    handler := handlers.NewHandler(service)

    // Инициализируем echo
    e := echo.New()
    // используем Logger и Recover
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())

    // Маршруты
    e.POST("/api/messages", func(c echo.Context) error {
        var request messages.PostMessagesRequestObject
        // Привязываем тело запроса
        if err := c.Bind(&request); err != nil {
            return c.String(http.StatusBadRequest, err.Error())
        }
    
        // Проверяем, что Body не nil
        if request.Body == nil || request.Body.Message == nil {
            return c.String(http.StatusBadRequest, "body or message cannot be nil")
        }
    
        response, err := handler.PostMessages(c.Request().Context(), request)
        if err != nil {
            return c.String(http.StatusInternalServerError, err.Error())
        }
    
        return c.JSON(http.StatusCreated, response)
    })

    e.GET("/api/messages", func(c echo.Context) error {
        response, err := handler.GetMessages(c.Request().Context(), messages.GetMessagesRequestObject{})
        if err != nil {
            return c.String(http.StatusInternalServerError, err.Error())
        }

        return c.JSON(http.StatusOK, response)
    })
    
    // PATCH запрос для обновления сообщения
    e.PATCH("/api/messages/:id", func(c echo.Context) error {
        messageID := c.Param("id") // Получаем ID сообщения из URL
        var request messages.PatchMessagesRequestObject // Изменяем тип на PatchMessagesRequestObject

        // Привязываем тело запроса
        if err := c.Bind(&request); err != nil {
            return c.String(http.StatusBadRequest, err.Error())
        }

        // Проверяем, что Body не nil
        if request.Body == nil || request.Body.Message == nil {
            return c.String(http.StatusBadRequest, "body or message cannot be nil")
        }

        // Преобразуем messageID из строки в uint
        id, err := strconv.ParseUint(messageID, 10, 32) // Преобразуем в uint
        if err != nil {
            return c.String(http.StatusBadRequest, "invalid message ID")
        }

        // Устанавливаем ID в Body запроса
        request.Body.Id = new(uint) // Создаем указатель на uint
        *request.Body.Id = uint(id)  // Присваиваем значение

        // Вызов функции для обновления сообщения
        response, err := handler.PatchMessage(c.Request().Context(), request)
        if err != nil {
            return c.String(http.StatusInternalServerError, err.Error())
        }

        return c.JSON(http.StatusOK, response)
    })

    // DELETE запрос для удаления сообщения
    e.DELETE("/api/messages/:id", func(c echo.Context) error {
        messageID := c.Param("id") // Получаем ID сообщения из URL

        // Преобразуем messageID из строки в uint
        id, err := strconv.ParseUint(messageID, 10, 32) // Преобразуем в uint
        if err != nil {
            return c.String(http.StatusBadRequest, "invalid message ID")
        }

        // Создаем объект запроса
        request := messages.DeleteMessagesRequestObject{
            Id: uint(id), // Присваиваем значение ID
        }

        // Вызов функции для удаления сообщения
        _, err = handler.DeleteMessages(c.Request().Context(), request)
        if err != nil {
            return c.String(http.StatusInternalServerError, err.Error())
        }

        return c.NoContent(http.StatusNoContent) // Возвращаем статус 204 No Content
    }) 

    // Запуск сервера на порту 8080
    if err := e.Start(":8080"); err != nil {
        log.Fatalf("failed to start server: %v", err)
    }
}