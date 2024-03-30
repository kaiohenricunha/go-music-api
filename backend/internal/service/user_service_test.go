package service

import (
	"errors"
	"testing"

	"github.com/kaiohenricunha/go-music-k8s/backend/internal/dao"
	"github.com/kaiohenricunha/go-music-k8s/backend/internal/dao/mocks"
	"github.com/kaiohenricunha/go-music-k8s/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestValidateUser(t *testing.T) {
	mockDAO := new(mocks.MusicDAO)
	userService := NewUserService(mockDAO)

	username := "testuser"
	password := "password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	mockUser := &model.User{Username: username, Password: string(hashedPassword)}
	id := uint(1)
	mockUser.ID = id

	// Expectation: FindByUsername is called and returns mockUser, nil
	mockDAO.On("FindByUsername", username).Return(mockUser, nil)

	userID, valid := userService.ValidateUser(username, password)
	assert.True(t, valid)
	assert.Equal(t, mockUser.ID, userID)

	mockDAO.AssertExpectations(t)
}

func TestValidateUser_UserNotFound(t *testing.T) {
	mockDAO := new(mocks.MusicDAO)
	userService := NewUserService(mockDAO)

	username := "testuser"
	password := "password"

	// Expectation: FindByUsername is called and returns nil, nil
	mockDAO.On("FindByUsername", username).Return(nil, nil)

	userID, valid := userService.ValidateUser(username, password)
	assert.False(t, valid)
	assert.Equal(t, uint(0), userID)

	mockDAO.AssertExpectations(t)
}

func TestValidateUser_PasswordMismatch(t *testing.T) {
	mockDAO := new(mocks.MusicDAO)
	userService := NewUserService(mockDAO)

	username := "testuser"
	password := "password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	mockUser := &model.User{Username: username, Password: string(hashedPassword)}
	id := uint(1)
	mockUser.ID = id

	// Expectation: FindByUsername is called and returns mockUser, nil
	mockDAO.On("FindByUsername", username).Return(mockUser, nil)

	userID, valid := userService.ValidateUser(username, "wrongpassword")
	assert.False(t, valid)
	assert.Equal(t, uint(0), userID)

	mockDAO.AssertExpectations(t)
}

func TestRegisterUser(t *testing.T) {
	mockDAO := new(mocks.MusicDAO)
	userService := NewUserService(mockDAO)

	username := "testuser"
	password := "password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	mockUser := &model.User{Username: username, Password: string(hashedPassword)}

	// Expectation: FindByUsername is called and returns nil, nil
	mockDAO.On("FindByUsername", username).Return(nil, nil)
	// Expectation: CreateUser is called with mockUser and returns nil
	mockDAO.On("CreateUser", mockUser).Return(nil)

	err := userService.RegisterUser(mockUser)
	assert.Nil(t, err)

	mockDAO.AssertExpectations(t)
}

func TestRegisterUser_UsernameTaken(t *testing.T) {
	mockDAO := new(mocks.MusicDAO)
	userService := NewUserService(mockDAO)

	username := "testuser"
	password := "password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	mockUser := &model.User{Username: username, Password: string(hashedPassword)}

	// Expectation: FindByUsername is called and returns mockUser, nil
	mockDAO.On("FindByUsername", username).Return(mockUser, nil)

	err := userService.RegisterUser(mockUser)
	assert.Equal(t, ErrUsernameTaken, err)

	mockDAO.AssertExpectations(t)
}

func TestUpdateUser(t *testing.T) {
	mockDAO := new(mocks.MusicDAO)
	userService := NewUserService(mockDAO)

	username := "testuser"
	password := "password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	mockUser := &model.User{Username: username, Password: string(hashedPassword)}

	// Expectation: UpdateUser is called with mockUser and returns nil
	mockDAO.On("UpdateUser", mockUser).Return(nil)

	err := userService.UpdateUser(mockUser)
	assert.Nil(t, err)

	mockDAO.AssertExpectations(t)
}

func TestDeleteUser(t *testing.T) {
	mockDAO := new(mocks.MusicDAO)
	userService := NewUserService(mockDAO)

	userID := uint(1)

	// Expectation: DeleteUser is called with userID and returns nil
	mockDAO.On("DeleteUser", userID).Return(nil)

	err := userService.DeleteUser(userID)
	assert.Nil(t, err)

	mockDAO.AssertExpectations(t)
}

func TestGetAllUsers(t *testing.T) {
	mockDAO := new(mocks.MusicDAO)
	userService := NewUserService(mockDAO)

	mockUsers := []model.User{
		{Username: "testuser1"},
		{Username: "testuser2"},
	}

	// Expectation: GetAllUsers is called and returns mockUsers, nil
	mockDAO.On("GetAllUsers").Return(mockUsers, nil)

	users, err := userService.GetAllUsers()
	assert.Nil(t, err)
	assert.Equal(t, mockUsers, users)

	mockDAO.AssertExpectations(t)
}

func TestFindUserByUsername(t *testing.T) {
	mockDAO := new(mocks.MusicDAO)
	userService := NewUserService(mockDAO)

	username := "testuser"
	mockUser := &model.User{Username: username}

	// Expectation: FindByUsername is called with username and returns mockUser, nil
	mockDAO.On("FindByUsername", username).Return(mockUser, nil)

	user, err := userService.FindUserByUsername(username)
	assert.Nil(t, err)
	assert.Equal(t, mockUser, user)

	mockDAO.AssertExpectations(t)
}

func TestFindUserByUsername_UserNotFound(t *testing.T) {
	mockDAO := new(mocks.MusicDAO)
	userService := NewUserService(mockDAO)

	username := "testuser"

	// Expectation: FindByUsername is called with username and returns nil, ErrUserNotFound
	mockDAO.On("FindByUsername", username).Return(nil, dao.ErrUserNotFound)

	user, err := userService.FindUserByUsername(username)
	assert.Nil(t, user)
	assert.Equal(t, dao.ErrUserNotFound, err)

	mockDAO.AssertExpectations(t)
}

func TestFindUserByUsername_Error(t *testing.T) {
	mockDAO := new(mocks.MusicDAO)
	userService := NewUserService(mockDAO)

	username := "testuser"

	// Expectation: FindByUsername is called with username and returns nil, errors.New("error")
	mockDAO.On("FindByUsername", username).Return(nil, errors.New("error"))

	user, err := userService.FindUserByUsername(username)
	assert.Nil(t, user)
	assert.NotNil(t, err)

	mockDAO.AssertExpectations(t)
}
