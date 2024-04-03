package service

import (
	"errors"
	"testing"

	"github.com/kaiohenricunha/go-music-k8s/backend/internal/dao/mocks"
	"github.com/kaiohenricunha/go-music-k8s/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestRegisterUser(t *testing.T) {
	mockDAO := &mocks.MusicDAO{}
	userService := NewUserService(mockDAO)

	password, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	testUser := &model.User{Username: "testUser", Password: string(password)}

	// Mock setup corrected to accurately reflect the logic flow
	// Scenario 1: User does not exist and is created successfully
	mockDAO.On("GetUserByUsername", "testUser").Return(nil, nil) // Simulate user not existing
	mockDAO.On("CreateUser", mock.AnythingOfType("*model.User")).Return(nil)

	err := userService.RegisterUser(testUser)
	assert.NoError(t, err)
	mockDAO.AssertExpectations(t)

	// Reset mock expectations between scenarios
	mockDAO.ExpectedCalls = nil
	mockDAO.Calls = nil

	// Scenario 2: User already exists
	mockDAO.On("GetUserByUsername", "testUser").Return(testUser, nil) // Simulate user existing

	err = userService.RegisterUser(testUser)
	assert.Equal(t, ErrUsernameOrEmailTaken, err)
	mockDAO.AssertExpectations(t)
}

func TestValidateUser(t *testing.T) {
	mockDAO := &mocks.MusicDAO{}
	userService := NewUserService(mockDAO)

	password := "password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := &model.User{Username: "testUser", Password: string(hashedPassword)}

	// Scenario 1: Successfully validate user
	mockDAO.On("GetUserByUsername", "testUser").Return(user, nil)

	userID, valid := userService.ValidateUser("testUser", password)
	assert.True(t, valid)
	assert.Equal(t, user.ID, userID)
	mockDAO.AssertExpectations(t)

	// Scenario 2: Invalid username
	mockDAO.On("GetUserByUsername", "invalidUser").Return(nil, ErrUserNotFound)

	_, valid = userService.ValidateUser("invalidUser", password)
	assert.False(t, valid)
	mockDAO.AssertExpectations(t)

	// Scenario 3: Invalid password
	_, valid = userService.ValidateUser("testUser", "wrongPassword")
	assert.False(t, valid)
	mockDAO.AssertExpectations(t)
}

func TestGetAllUsers(t *testing.T) {
	mockDAO := &mocks.MusicDAO{}
	userService := NewUserService(mockDAO)

	users := []model.User{
		{Username: "user1"},
		{Username: "user2"},
	}

	// Scenario 1: Successfully retrieve all users
	mockDAO.On("GetAllUsers").Return(users, nil).Once()

	result, err := userService.GetAllUsers()
	assert.NoError(t, err, "Expected no error for GetAllUsers")
	assert.Len(t, result, 2, "Expected result to have length 2")
	assert.Equal(t, users, result, "Expected returned users to match the mock users")
	mockDAO.AssertExpectations(t)

	// Reset mock expectations to test the failure scenario
	mockDAO.ExpectedCalls = nil
	mockDAO.Calls = nil

	// Scenario 2: Error retrieving users
	mockDAO.On("GetAllUsers").Return(nil, errors.New("error")).Once()

	_, err = userService.GetAllUsers()
	assert.Error(t, err, "Expected an error for GetAllUsers")
	mockDAO.AssertExpectations(t)
}

func TestGetUserByUsername(t *testing.T) {
	mockDAO := &mocks.MusicDAO{}
	userService := NewUserService(mockDAO)

	testUser := &model.User{Username: "testUser", Password: "hashedpassword"}

	// Scenario 1: Successfully retrieve a user by username
	mockDAO.On("GetUserByUsername", "testUser").Return(testUser, nil) // Simulates user exists
	err := userService.RegisterUser(testUser)
	assert.Equal(t, ErrUsernameOrEmailTaken, err) // Verifies the correct error is returned

	// Scenario 2: User not found
	mockDAO.On("GetUserByUsername", "nonExistingUser").Return(nil, ErrUserNotFound)

	_, err = userService.GetUserByUsername("nonExistingUser")
	assert.Error(t, err)
	assert.Equal(t, ErrUserNotFound, err)
	mockDAO.AssertExpectations(t)
}
