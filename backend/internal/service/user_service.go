package service

import (
	"errors"
	"log"

	"github.com/kaiohenricunha/go-music-k8s/backend/internal/dao"
	"github.com/kaiohenricunha/go-music-k8s/backend/internal/model"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrUserNotFound is returned when a user is not found.
	ErrUserNotFound = errors.New("user not found")

	// ErrUsernameTaken is returned when the username is already used by another user.
	ErrUsernameTaken = errors.New("username already taken")

	// ErrInvalidCredentials is returned when the username or password is incorrect.
	ErrInvalidCredentials = errors.New("invalid username or password")
)

// UserService outlines the interface for user-related operations.
type UserService interface {
	ValidateUser(username, password string) (uint, bool)
	RegisterUser(user *model.User) error
	UpdateUser(user *model.User) error
	DeleteUser(userID uint) error
	GetAllUsers() ([]model.User, error)
	FindUserByUsername(username string) (*model.User, error)
}

type userService struct {
	userDAO dao.MusicDAO
}

func NewUserService(userDAO dao.MusicDAO) UserService {
	return &userService{userDAO: userDAO}
}

// ValidateUser checks if the username and password are correct.
func (us *userService) ValidateUser(username, password string) (uint, bool) {
	log.Printf("Validating user: %s\n", username)
	user, err := us.userDAO.FindByUsername(username)
	if err != nil || user == nil {
		log.Printf("User not found: %s\n", username)
		return 0, false
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Printf("Password comparison failed for user: %s\n", username)
		return 0, false
	}

	return user.ID, true
}

// RegisterUser handles registering a new user with hashed password.
func (us *userService) RegisterUser(user *model.User) error {
	// Check if username already exists
	existingUser, _ := us.FindUserByUsername(user.Username)
	if existingUser != nil {
		return ErrUsernameTaken
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	// Create the user
	return us.userDAO.CreateUser(user)
}

// UpdateUser updates an existing user's information.
func (us *userService) UpdateUser(user *model.User) error {
	// If a new password is provided, hash it before saving.
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
	}

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
