package store

import (
	"refactoring/internal/app/models"
)

type UserRepository interface {
	SearchUsers() (map[string]models.User, error)
	CreateUser(*models.User) (string, error)
	GetUser(id string) (*models.User, error)
	UpdateUser(id string, u *models.User) error
	DeleteUser(id string) error
}
