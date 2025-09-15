// repository/task_repository.go
package repository

import (
	"gorm.io/gorm"
	"todo-list/backend/internal/models"
)

// NewTaskRepository создает новый репозиторий задач (для совместимости с start.go)
func NewTaskRepository(db *gorm.DB) *Repository {
	return NewRepository(db)
}

// TodoRepository интерфейс для работы с задачами
type TodoRepository interface {
	Create(todo *models.Todo) error
	GetByID(id uint) (*models.Todo, error)
	GetAll() ([]models.Todo, error)
	GetByUserID(userID uint) ([]models.Todo, error)
	Update(todo *models.Todo) error
	Delete(id uint) error
	GetByStatus(completed bool) ([]models.Todo, error)
}

// UserRepository интерфейс для работы с пользователями
type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
}

// CategoryRepository интерфейс для работы с категориями
type CategoryRepository interface {
	Create(category *models.Category) error
	GetByID(id uint) (*models.Category, error)
	GetAll() ([]models.Category, error)
	Update(category *models.Category) error
	Delete(id uint) error
}

// Repository объединяет все репозитории
type Repository struct {
	Todo     TodoRepository
	User     UserRepository
	Category CategoryRepository
}

// todoRepo реализация TodoRepository
type todoRepo struct {
	db *gorm.DB
}

// userRepo реализация UserRepository
type userRepo struct {
	db *gorm.DB
}

// categoryRepo реализация CategoryRepository
type categoryRepo struct {
	db *gorm.DB
}

// NewRepository создает новый экземпляр Repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Todo:     &todoRepo{db: db},
		User:     &userRepo{db: db},
		Category: &categoryRepo{db: db},
	}
}

// Реализация TodoRepository
func (r *todoRepo) Create(todo *models.Todo) error {
	return r.db.Create(todo).Error
}

func (r *todoRepo) GetByID(id uint) (*models.Todo, error) {
	var todo models.Todo
	err := r.db.Preload("User").Preload("Category").First(&todo, id).Error
	return &todo, err
}

func (r *todoRepo) GetAll() ([]models.Todo, error) {
	var todos []models.Todo
	err := r.db.Preload("User").Preload("Category").Find(&todos).Error
	return todos, err
}

func (r *todoRepo) GetByUserID(userID uint) ([]models.Todo, error) {
	var todos []models.Todo
	err := r.db.Where("user_id = ?", userID).Preload("Category").Find(&todos).Error
	return todos, err
}

func (r *todoRepo) Update(todo *models.Todo) error {
	return r.db.Save(todo).Error
}

func (r *todoRepo) Delete(id uint) error {
	return r.db.Delete(&models.Todo{}, id).Error
}

func (r *todoRepo) GetByStatus(completed bool) ([]models.Todo, error) {
	var todos []models.Todo
	err := r.db.Where("completed = ?", completed).Preload("User").Preload("Category").Find(&todos).Error
	return todos, err
}

// Реализация UserRepository
func (r *userRepo) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepo) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Todos").First(&user, id).Error
	return &user, err
}

func (r *userRepo) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *userRepo) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *userRepo) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepo) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

// Реализация CategoryRepository
func (r *categoryRepo) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepo) GetByID(id uint) (*models.Category, error) {
	var category models.Category
	err := r.db.Preload("Todos").First(&category, id).Error
	return &category, err
}

func (r *categoryRepo) GetAll() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Find(&categories).Error
	return categories, err
}

func (r *categoryRepo) Update(category *models.Category) error {
	return r.db.Save(category).Error
}

func (r *categoryRepo) Delete(id uint) error {
	return r.db.Delete(&models.Category{}, id).Error
}
