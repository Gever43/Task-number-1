package main

import (
	"log"
	"net/http"

	"myProject/internal/database"
	"myProject/internal/handlers"
	"myProject/internal/messagesService"

	"github.com/gorilla/mux"
)

func main() {
    // Инициализация базы данных
    database.InitDB()
    database.DB.AutoMigrate(&messagesService.Message{}) // Автоматическая миграция структуры Message (обновление)

    // Создание репозитория и сервиса
    repo := messagesService.NewMessageRepository(database.DB)
    service := messagesService.NewService(repo)

    // Создание хендлера
    handler := handlers.NewHandler(service)

    // Создание маршрутизатора
    router := mux.NewRouter()
    router.HandleFunc("/api/messages", handler.GetMessagesHandler).Methods("GET")
    router.HandleFunc("/api/messages", handler.PostMessageHandler).Methods("POST")
    router.HandleFunc("/api/messages/{id}", handler.PatchMessageHandler).Methods("PATCH")
    router.HandleFunc("/api/messages/{id}", handler.DeleteMessageHandler).Methods("DELETE")

    // Запуск сервера
    log.Println("Server is running on port 8080...")
    if err := http.ListenAndServe(":8080", router); err != nil {
        log.Fatal(err)
    }
}