///ПРИМЕР РАБОТЫ
/*
package main

import (
	"fmt"
	"net/http"












	"github.com/gorilla/mux"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Hello, World!")
}

func main() {
    router := mux.NewRouter()
    // наше приложение будет слушать запросы на localhost:8080/api/hello
    router.HandleFunc("/api/hello", HelloHandler).Methods("GET")
    http.ListenAndServe(":8080", router)
}
*/
package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

var message string                                                                                          //хранит сообщение от клиента
                                                                                                            
type requestBody struct {                                                                                   // определяем структуру requestBody, которая содержит одно поле Message
    Message string `json:"message"`                                                                         // добавляем тег для сериализации. В JSON Message будет message
}
                                                                                                            // Обработчики запросов Сперва POST
func postHandler(w http.ResponseWriter, r *http.Request) {                                                  // w - позволяет формировать ответ клиенту, r сслыка на структура запроса
    var reqBody requestBody                                                                                 // переменная для хранения данных, полученных от клиента
    decoder := json.NewDecoder(r.Body)                                                                      // декодер для десериализации. Из JSON в структуры. Body это поле структуры Request
    err := decoder.Decode(&reqBody)                                                                         // декодирование. Метод Decode принимает данные, которые необходимо преобразовать
    if err != nil {                                                                                         // обработка ошибки: встроенная функция позволяет выводить ошибку в правильном формате
        http.Error(w, err.Error(), http.StatusBadRequest)                                                   // 3 аргумента - формировщик ответа / тип ошибки / код состояния HTTP
        return                                                                                              // возвращение значения
    }
    message = reqBody.Message                                                                               // запись в переменную через поле структуры
    fmt.Fprintf(w, "Message received: %s", message)                                                         // ответ (куда/ что / и ещё переменная, которую передаём )
}

func getHandler(w http.ResponseWriter, r *http.Request) {                                                   // обработчик GET-запросов.
    fmt.Fprintf(w, "Hello, %s!", message)                                                                   // выводит приветствие и сообщение с глобальной переменнной
}

func main() {
    router := mux.NewRouter()                                                                               // создание маршрутизатора

    // Регистрация маршрутов                                                                                (путь/функция-обработчик/метод ограничитель - это здесь)

    router.HandleFunc("/post", postHandler).Methods("POST")                                                 // если запрос на /post, то вызывается функция postHandler. Methods("POST") только для обработки POST запросов
    router.HandleFunc("/get", getHandler).Methods("GET")                                                    // регистрация GET запроса и проверка на метод GET, чтобы другие не обрабатывать

    fmt.Println("Server is listening on port 8080...")
    http.ListenAndServe(":8080", router)                                                                    // запуск HTTP-сервера, позволяя обрабатывать запросы и связывать их с обработчиками
}