// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	database "github.com/DmytroKha/nix-chat/internal/infra/database"
	mock "github.com/stretchr/testify/mock"
)

// BlacklistService is an autogenerated mock type for the BlacklistService type
type BlacklistService struct {
	mock.Mock
}

// Delete provides a mock function with given fields: id
func (_m *BlacklistService) Delete(id int64) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(int64) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Find provides a mock function with given fields: userId, roomId
func (_m *BlacklistService) Find(userId int64, roomId int64) (database.Blacklist, error) {
	ret := _m.Called(userId, roomId)

	var r0 database.Blacklist
	if rf, ok := ret.Get(0).(func(int64, int64) database.Blacklist); ok {
		r0 = rf(userId, roomId)
	} else {
		r0 = ret.Get(0).(database.Blacklist)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64, int64) error); ok {
		r1 = rf(userId, roomId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindAll provides a mock function with given fields: userId
func (_m *BlacklistService) FindAll(userId int64) ([]database.Blacklist, error) {
	ret := _m.Called(userId)

	var r0 []database.Blacklist
	if rf, ok := ret.Get(0).(func(int64) []database.Blacklist); ok {
		r0 = rf(userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]database.Blacklist)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: bl
func (_m *BlacklistService) Save(bl database.Blacklist) (database.Blacklist, error) {
	ret := _m.Called(bl)

	var r0 database.Blacklist
	if rf, ok := ret.Get(0).(func(database.Blacklist) database.Blacklist); ok {
		r0 = rf(bl)
	} else {
		r0 = ret.Get(0).(database.Blacklist)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(database.Blacklist) error); ok {
		r1 = rf(bl)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewBlacklistService interface {
	mock.TestingT
	Cleanup(func())
}

// NewBlacklistService creates a new instance of BlacklistService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewBlacklistService(t mockConstructorTestingTNewBlacklistService) *BlacklistService {
	mock := &BlacklistService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
