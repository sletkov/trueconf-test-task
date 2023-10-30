package jsonstore

import (
	"refactoring/internal/app/models"
	"refactoring/internal/app/store"
)

type UserList map[string]models.User

type UserStore struct {
	Increment      int      `json:"increment"`
	List           UserList `json:"list"`
	userRepository *UserRepository
	userStoreUrl   string
}

func New(userStoreUrl string) *UserStore {
	return &UserStore{
		userStoreUrl: userStoreUrl,
	}
}

func (s *UserStore) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}
