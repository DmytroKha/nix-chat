// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	"errors"
	database "github.com/DmytroKha/nix-chat/internal/infra/database"
	mock "github.com/stretchr/testify/mock"

	requests "github.com/DmytroKha/nix-chat/internal/infra/http/requests"
)

// AuthService is an autogenerated mock type for the AuthService type
type AuthService struct {
	mock.Mock
}

// GenerateJwt provides a mock function with given fields: user
func (_m *AuthService) GenerateJwt(user database.User) (string, error) {
	ret := _m.Called(user)

	var r0 string
	if rf, ok := ret.Get(0).(func(database.User) string); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(database.User) error); ok {
		r1 = rf(user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Login provides a mock function with given fields: user
func (_m *AuthService) Login(user requests.UserLoginRequest) (database.User, string, error) {
	if user.Name == "userabsent" {
		return database.User{}, "", errors.New("User doesn`t exist")
	}

	ret := _m.Called(user)

	var r0 database.User
	if rf, ok := ret.Get(0).(func(requests.UserLoginRequest) database.User); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Get(0).(database.User)
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(requests.UserLoginRequest) string); ok {
		r1 = rf(user)
	} else {
		r1 = ret.Get(1).(string)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(requests.UserLoginRequest) error); ok {
		r2 = rf(user)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2

	//if user.Name == "Dmytro"{
	//	return database.User{}, "", errors.New("User doesn`t exist")
	//} else {
	//	var testUser = database.User{
	//		Id:   1,
	//		Name: "testuser",
	//		Image: database.Image{
	//			Name: "test.jpg",
	//		},
	//	}
	//	return testUser, "token123", nil
	//}
}

// Register provides a mock function with given fields: user
func (_m *AuthService) Register(user requests.UserRegistrationRequest) (database.User, string, error) {
	if user.Name == "userexist" {
		return database.User{}, "", errors.New("User already exists")
	}

	ret := _m.Called(user)

	var r0 database.User
	if rf, ok := ret.Get(0).(func(requests.UserRegistrationRequest) database.User); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Get(0).(database.User)
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(requests.UserRegistrationRequest) string); ok {
		r1 = rf(user)
	} else {
		r1 = ret.Get(1).(string)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(requests.UserRegistrationRequest) error); ok {
		r2 = rf(user)
	} else {
		r2 = ret.Error(2)
	}
	return r0, r1, r2

	//if user.Name == "Dmytro"{
	//	return database.User{}, "", errors.New("User already exists")
	//} else {
	//	var testUser = database.User{
	//		Id:   1,
	//		Name: "testuser",
	//	}
	//	return testUser, "token123", nil
	//}

}

type mockConstructorTestingTNewAuthService interface {
	mock.TestingT
	Cleanup(func())
}

// NewAuthService creates a new instance of AuthService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAuthService(t mockConstructorTestingTNewAuthService) *AuthService {
	mock := &AuthService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}