package services

import (
	"errors"
	"sync"

	"github.com/igralkin/go-highload/models"
)

type UserService struct {
	mu     sync.RWMutex
	users  map[int]models.User
	nextID int
}

func NewUserService() *UserService {
	return &UserService{
		users:  make(map[int]models.User),
		nextID: 1,
	}
}

func (s *UserService) Create(name, email string) models.User {
	s.mu.Lock()
	defer s.mu.Unlock()

	user := models.User{
		ID:    s.nextID,
		Name:  name,
		Email: email,
	}
	s.users[user.ID] = user
	s.nextID++
	return user
}

func (s *UserService) GetAll() []models.User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]models.User, 0, len(s.users))
	for _, u := range s.users {
		result = append(result, u)
	}
	return result
}

func (s *UserService) GetByID(id int) (models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.users[id]
	if !ok {
		return models.User{}, errors.New("user not found")
	}
	return user, nil
}

func (s *UserService) Update(id int, name, email string) (models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, ok := s.users[id]
	if !ok {
		return models.User{}, errors.New("user not found")
	}

	user.Name = name
	user.Email = email
	s.users[id] = user
	return user, nil
}

func (s *UserService) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.users[id]; !ok {
		return errors.New("user not found")
	}
	delete(s.users, id)
	return nil
}
