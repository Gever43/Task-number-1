package userService

// UserService представляет сервис для работы с пользователями
type UserService struct {
    Repo UserRepositoryInterface
}

// NewService - создает новый экземпляр UserService
func NewUserService(repo UserRepositoryInterface) *UserService {
    return &UserService{Repo: repo}
}

// CreateUser создает нового пользователя
func (s *UserService) CreateUser(name string, email string, password string) (DBUser, error) {
    user := DBUser{
        Name:     name,
        Email:    email,
        Password: password,
    }
    return s.Repo.CreateUser(user)
}

// GetAllUsers возвращает всех пользователей
func (s *UserService) GetAllUsers() ([]DBUser, error) {
    return s.Repo.GetAllUsers()
}

// UpdateUserByID обновляет пользователя по ID
func (s *UserService) UpdateUserByID(id uint, user DBUser) (DBUser, error) {
    return s.Repo.UpdateUser(id, user)
}

// DeleteUserByID удаляет пользователя по ID
func (s *UserService) DeleteUserByID(id uint) error {
    return s.Repo.DeleteUser(id)
}