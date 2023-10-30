package jsonstore

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"refactoring/internal/app/models"
	"refactoring/internal/app/store"
	"strconv"
)

type UserRepository struct {
	store *UserStore
}

func (r *UserRepository) SearchUsers() (map[string]models.User, error) {
	store, err := getStoreFromJSONFile(r.store.userStoreUrl)

	if err != nil {
		return nil, err
	}

	return store.List, nil
}

func (r *UserRepository) CreateUser(u *models.User) (string, error) {
	s, err := getStoreFromJSONFile(r.store.userStoreUrl)

	if err != nil {
		return "", err
	}

	s.Increment++

	id := strconv.Itoa(s.Increment)
	s.List[id] = *u

	b, _ := json.Marshal(&s)
	_ = ioutil.WriteFile(r.store.userStoreUrl, b, fs.ModePerm)
	return id, nil
}

func (r *UserRepository) GetUser(id string) (*models.User, error) {
	s, err := getStoreFromJSONFile(r.store.userStoreUrl)

	if err != nil {
		return nil, err
	}

	user, ok := s.List[id]

	if !ok {
		return nil, store.ErrUserNotFound
	}

	return &user, nil
}

func (r *UserRepository) UpdateUser(id string, u *models.User) error {
	s, err := getStoreFromJSONFile(r.store.userStoreUrl)

	if err != nil {
		return err
	}

	user := s.List[id]
	user.DisplayName = u.DisplayName
	s.List[id] = user

	err = writeStoreIntoJSONFile(s, r.store.userStoreUrl)

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) DeleteUser(id string) error {
	s, err := getStoreFromJSONFile(r.store.userStoreUrl)

	if err != nil {
		return err
	}

	delete(s.List, id)

	err = writeStoreIntoJSONFile(s, r.store.userStoreUrl)

	if err != nil {
		return err
	}

	return nil
}
