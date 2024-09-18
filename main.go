// /Задание 2
package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Структура для запроса
type requestBody struct {
    Message string `json:"message"` 
}

// Обработчик POST запросов для создания сообщения
func postHandler(w http.ResponseWriter, r *http.Request) {
    var reqBody requestBody
    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&reqBody)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    message := Message{Text: reqBody.Message} //инициализация переменной полем структуры
    if err := DB.Create(&message).Error; err != nil { //через DB подключаемся к БД и создаём методом Create запись в БД
        http.Error(w, err.Error(), http.StatusInternalServerError) //обработка ошибок
        return
    }
    fmt.Fprintf(w, "Message received: %s", message.Text)
}

// Обработчик GET запросов для получения всех сообщений
func getHandler(w http.ResponseWriter, r *http.Request) {
    var messages []Message //слайс для массивов ответов
    if err := DB.Find(&messages).Error; err != nil { //Находим все записи в таблице и записываем в слайс
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
//установка заголовка
    w.Header().Set("Content-Type", "application/json") //SET - устанавливает значение заголовка. Content-Type-заголовок тип содержимого. Аpplication/json - значение заголовка (говорит, что тип будет JSON)
    json.NewEncoder(w).Encode(messages) //кодирование слайса messages в JSON и отправка как ответ
}

// Основная функция
func main() {
    InitDB() //установка соединения с БД
    DB.AutoMigrate(&Message{}) //изменение БД в соответствии моделью (структурой Message)

    router := mux.NewRouter()
    router.HandleFunc("/api/messages", postHandler).Methods("POST")
    router.HandleFunc("/api/messages", getHandler).Methods("GET")
    http.ListenAndServe(":8080", router)
}