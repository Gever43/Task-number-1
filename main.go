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

// Обработчик PATCH запросов для обновления сообщения по ID
func patchHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r) // Получаем переменные маршрута (создаем мап)
    id := vars["id"]    // Извлекаем ID из мата по ключу

    var reqBody requestBody // Структура для хранения данных запроса
    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&reqBody)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Обновляем сообщение
    var message Message
    if err := DB.First(&message, id).Error; err != nil {   //извлечение из БД первого значения по заданным условиям (куда запишем, условие)
        http.Error(w, "Message not found", http.StatusNotFound)
        return
    }

    message.Text = reqBody.Message // Обновляем текст сообщения
    if err := DB.Save(&message).Error; err != nil {//сохраняем куда
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    fmt.Fprintf(w, "Message updated: %s", message.Text)
}

// Обработчик DELETE запросов для удаления сообщения по ID
func deleteHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r) // Получаем переменные маршрута
    id := vars["id"]    // Извлекаем ID

    if err := DB.Delete(&Message{}, id).Error; err != nil {//удаляем откуда (где искать значения), искать по какому значению
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    fmt.Fprintf(w, "Message deleted with ID: %s", id)
}
// Основная функция
func main() {
    InitDB() //установка соединения с БД
    DB.AutoMigrate(&Message{}) //изменение БД в соответствии моделью (структурой Message)

    router := mux.NewRouter()
    router.HandleFunc("/api/messages", postHandler).Methods("POST")
    router.HandleFunc("/api/messages", getHandler).Methods("GET")
    router.HandleFunc("/api/messages/{id}", patchHandler).Methods("PATCH") 
    router.HandleFunc("/api/messages/{id}", deleteHandler).Methods("DELETE") 
    http.ListenAndServe(":8080", router)
}