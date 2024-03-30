package service

import (
	"errors"

	"github.com/kaiohenricunha/go-music-k8s/backend/internal/dao"
	"github.com/kaiohenricunha/go-music-k8s/backend/internal/model"
	"golang.org/x/crypto/bcrypt"
)

// ErrUserNotFound is returned when a user is not found.
var ErrUserNotFound = errors.New("user not found")

// UserService outlines the interface for user-related operations.
type UserService interface {
	ValidateUser(username, password string) bool
	RegisterUser(user *model.User) error
	UpdateUser(user *model.User) error
	DeleteUser(userID uint) error
	GetAllUsers() ([]model.User, error)
	FindUserByUsername(username string) (*model.User, error)
}

// userService is the concrete implementation of UserService.
type userService struct {
	userDAO dao.MusicDAO
}

type UserDTO struct {
	ID       uint   `json:"id,omitempty"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// NewUserService creates a new UserService with the given MusicDAO.
func NewUserService(userDAO dao.MusicDAO) UserService {
	return &userService{userDAO: userDAO}
}

// ValidateUser checks if the username and password are correct.
func (us *userService) ValidateUser(username, password string) bool {
	user, err := us.userDAO.FindByUsername(username)
	if err != nil || user == nil {
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

// RegisterUser handles registering a new user with hashed password.
func (us *userService) RegisterUser(user *model.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return us.userDAO.CreateUser(user)
}

// UpdateUser updates an existing user's information.
func (us *userService) UpdateUser(user *model.User) error {
	// TODO: Hash password if it's being changed.
	return us.userDAO.UpdateUser(user)
}

// DeleteUser removes a user by their ID.
func (us *userService) DeleteUser(userID uint) error {
	return us.userDAO.DeleteUser(userID)
}

// GetAllUsers retrieves all users.
func (us *userService) GetAllUsers() ([]model.User, error) {
	return us.userDAO.GetAllUsers()
}

// FindUserByUsername retrieves a user by their username.
func (us *userService) FindUserByUsername(username string) (*model.User, error) {
	user, err := us.userDAO.FindByUsername(username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}
