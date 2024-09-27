package userService

import (
	"gorm.io/gorm"
)

// User определяет структуру пользователя для работы с БД
type DBUser struct {
    gorm.Model
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"password"`
}