package userService

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// UserRepositoryInterface интерфейс для работы с пользователями
type UserRepositoryInterface interface {
    CreateUser(user DBUser) (DBUser, error)

    GetAllUsers() ([]DBUser, error)

    UpdateUser(id uint, user DBUser) (DBUser, error)

    DeleteUser(id uint) error
}
// UserRepository - структура, поддерживающая интерфейс, т.к. поддерживает все методы интерфейса
type UserRepository struct {
    DB *gorm.DB
}

// NewUserRepository - создает новый экземпляр UserRepository
func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{DB: db} // инициализация поля DB структуры UserRepository переданным аргументом db
}

// CreateUser создает нового пользователя
func (repo *UserRepository) CreateUser(user DBUser) (DBUser, error) {
    // Проверка на существование пользователя с таким же email
    var existingUser DBUser
    if err := repo.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
        return DBUser{}, errors.New("пользователь с таким email уже существует")
    }

    // Создание нового пользователя
    if err := repo.DB.Create(&user).Error; err != nil {
        return DBUser{}, err
    }
    return user, nil
}

// GetAllUsers возвращает всех пользователей
func (repo *UserRepository) GetAllUsers() ([]DBUser, error) {
    var users []DBUser
    if err := repo.DB.Find(&users).Error; err != nil {
        return nil, err
    }
    return users, nil
}

// UpdateUser - реализация метода для обновления информации о пользователе по ID
func (repo *UserRepository) UpdateUser(id uint, user DBUser) (DBUser, error) {
    var existingUser DBUser

    // Проверяем, существует ли пользователь с заданным ID
    if err := repo.DB.First(&existingUser, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return DBUser{}, fmt.Errorf("user with id %d not found", id)
        }
        return DBUser{}, err // Возвращаем ошибку, если произошла другая ошибка
    }

    // Обновляем поля существующего пользователя
    existingUser.Email = user.Email
    existingUser.Password = user.Password

    // Сохраняем изменения в базе данных
    if err := repo.DB.Save(&existingUser).Error; err != nil {
        return DBUser{}, err // Возвращаем ошибку, если сохранение не удалось
    }

    return existingUser, nil // Возвращаем обновленного пользователя
}

// DeleteUser удаляет пользователя по ID
func (repo *UserRepository) DeleteUser(id uint) error {
    var user DBUser
    // Проверяем, существует ли пользователь
    if err := repo.DB.First(&user, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return fmt.Errorf("user with id %d not found", id) // Пользователь не найден
        }
        return err // Возвращаем другую ошибку
    }

    // Удаляем пользователя
    return repo.DB.Delete(&user).Error
}