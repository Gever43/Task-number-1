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
    if err := database.DB.AutoMigrate(&messagesService.DBMessage{}); err != nil {
        log.Fatalf("failed to migrate database: %v", err)
    }

    //разделение логики на слои
    repo := messagesService.NewMessageRepository(database.DB)//слой доступа к данным. Предоставляет интерфейс для взаимодействия с источником данных (БД)
    service := messagesService.NewService(repo)//сервисный слой - бизнес-логика приложения. Взаимодействует с обработчиками и взаимодействует с репозиториями для доступа к данным
    handler := handlers.NewHandler(service)//слой обработки запросов
    

    // Инициализируем echo - фреймворк, упрощающий создание веб-приложений
    e := echo.New()
    // используем Logger и Recover
    e.Use(middleware.Logger())//логирование (отслеживание работы приложения)
    e.Use(middleware.Recover())//восстановление после ошибок

    // Обработчики
    e.POST("/api/messages", func(c echo.Context) error {
        var request messages.PostMessagesRequestObject//хранение данных тела запроса
        // Десериализация (призявка) данных из тела запроса к структуре request
        if err := c.Bind(&request); err != nil {
            return c.String(http.StatusBadRequest, err.Error())//код HTTP и сообщение об ошибке
        }
    
        // Проверяем, что Body не nil (тело запроса и сообщение не пустые)
        if request.Body == nil || request.Body.Message == nil {
            return c.String(http.StatusBadRequest, "body or message cannot be nil")
        }
        //создание сообщения (передаём контекст - всю информацию о сообщении)
        response, err := handler.PostMessages(c.Request().Context(), request)
        if err != nil {
            return c.String(http.StatusInternalServerError, err.Error())
        }
        //если всё успешно, то выдаём JSON клиенту
        return c.JSON(http.StatusCreated, response)
    })

    e.GET("/api/messages", func(c echo.Context) error {
        var request messages.GetMessagesRequestObject//(можно так не делать а сразу вниз messages.GetMessagesRequestObject) - хранение данных тела запроса
        //- извлечение сообщений из БД
        response, err := handler.GetMessages(c.Request().Context(), request)
        if err != nil {
            return c.String(http.StatusInternalServerError, err.Error())
        }
        //если всё успешно,то выдаём JSON клиенту
        return c.JSON(http.StatusOK, response)
    })
    
    // PATCH запрос для обновления сообщения
    e.PATCH("/api/messages/:id", func(c echo.Context) error {
        messageID := c.Param("id") // Получаем ID сообщения из URL
        var request messages.PatchMessagesRequestObject //хранение данных тела запроса

        // Десериализация к структуре request
        if err := c.Bind(&request); err != nil {
            return c.String(http.StatusBadRequest, err.Error())
        }

        // Проверяем, что Body не nil
        if request.Body == nil || request.Body.Message == nil {
            return c.String(http.StatusBadRequest, "body or message cannot be nil")
        }

        // Преобразуем messageID из строки в uint 64. 10 - основание системы счисления, 32 - битность, которой ограничено значение. Формат uint 64 всё равно сохранится
        id, err := strconv.ParseUint(messageID, 10, 32) // Преобразуем в uint 32
        if err != nil {
            return c.String(http.StatusBadRequest, "invalid message ID")
        }

        // Устанавливаем ID в Body запроса
        request.Body.Id = new(uint) // Создаем указатель на uint
        *request.Body.Id = uint(id)  // Присваиваем значение, меняя тип переменной с uint 64 на uint.

        // Обновление сообщения
        response, err := handler.PatchMessage(c.Request().Context(), request)
        if err != nil {
            return c.String(http.StatusInternalServerError, err.Error())
        }

        return c.JSON(http.StatusOK, response)
    })

    // DELETE запрос для удаления сообщения
    e.DELETE("/api/messages/:id", func(c echo.Context) error {
        messageID := c.Param("id") // Получаем ID сообщения из URL

        // Преобразуем messageID из строки в uint 64
        id, err := strconv.ParseUint(messageID, 10, 32) // Преобразуем в uint 64
        if err != nil {
            return c.String(http.StatusBadRequest, "invalid message ID")
        }

        //хранение данных тела запроса
        var request = messages.DeleteMessagesRequestObject{
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