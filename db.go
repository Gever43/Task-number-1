package main

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// переменная, через которую мы будем работать с БД
var DB *gorm.DB                                                                                                             //данный тип нужен, чтобы взаимодействовать с БД

func InitDB() {                                                                                                             //здесь происходит открытие БД
// в dsn вводим данные, которые мы указали при создании контейнера
    dsn := "host=localhost user=postgres password=yourpassword dbname=postgres port=5432 sslmode=disable"                   //строка подключения (параметры подключения)
    var err error
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})                                                                 //попытка открыть соединение с БД (какую БД/конфигурация подключения - тут стандарт)
    if err != nil {
        log.Fatal("Failed to connect to database: ", err)                                                                   //обработка ошибок и закрытие
    }
}