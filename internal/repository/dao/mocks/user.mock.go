// Code generated by MockGen. DO NOT EDIT.
// Source: ./webook/internal/repository/dao/user.go
//
// Generated by this command:
//
//	mockgen -source=./webook/internal/repository/dao/user.go -package=daomocks -destination=./webook/internal/repository/dao/mocks/user.mock.go
//

// Package daomocks is a generated GoMock package.
package daomocks

import (
	domain "basic-go/webook/internal/domain"
	dao "basic-go/webook/internal/repository/dao"
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockUserDAO is a mock of UserDAO interface.
type MockUserDAO struct {
	ctrl     *gomock.Controller
	recorder *MockUserDAOMockRecorder
}

// MockUserDAOMockRecorder is the mock recorder for MockUserDAO.
type MockUserDAOMockRecorder struct {
	mock *MockUserDAO
}

// NewMockUserDAO creates a new mock instance.
func NewMockUserDAO(ctrl *gomock.Controller) *MockUserDAO {
	mock := &MockUserDAO{ctrl: ctrl}
	mock.recorder = &MockUserDAOMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserDAO) EXPECT() *MockUserDAOMockRecorder {
	return m.recorder
}

// EditProfile mocks base method.
func (m *MockUserDAO) EditProfile(c context.Context, user domain.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EditProfile", c, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// EditProfile indicates an expected call of EditProfile.
func (mr *MockUserDAOMockRecorder) EditProfile(c, user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditProfile", reflect.TypeOf((*MockUserDAO)(nil).EditProfile), c, user)
}

// FindByPhone mocks base method.
func (m *MockUserDAO) FindByPhone(c context.Context, phone string) (dao.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByPhone", c, phone)
	ret0, _ := ret[0].(dao.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByPhone indicates an expected call of FindByPhone.
func (mr *MockUserDAOMockRecorder) FindByPhone(c, phone any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByPhone", reflect.TypeOf((*MockUserDAO)(nil).FindByPhone), c, phone)
}

// FindbyEmail mocks base method.
func (m *MockUserDAO) FindbyEmail(ctx context.Context, email string) (dao.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindbyEmail", ctx, email)
	ret0, _ := ret[0].(dao.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindbyEmail indicates an expected call of FindbyEmail.
func (mr *MockUserDAOMockRecorder) FindbyEmail(ctx, email any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindbyEmail", reflect.TypeOf((*MockUserDAO)(nil).FindbyEmail), ctx, email)
}

// Insert mocks base method.
func (m *MockUserDAO) Insert(ctx context.Context, u dao.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", ctx, u)
	ret0, _ := ret[0].(error)
	return ret0
}

// Insert indicates an expected call of Insert.
func (mr *MockUserDAOMockRecorder) Insert(ctx, u any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockUserDAO)(nil).Insert), ctx, u)
}

// Profile mocks base method.
func (m *MockUserDAO) Profile(c context.Context, user domain.User) (dao.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Profile", c, user)
	ret0, _ := ret[0].(dao.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Profile indicates an expected call of Profile.
func (mr *MockUserDAOMockRecorder) Profile(c, user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Profile", reflect.TypeOf((*MockUserDAO)(nil).Profile), c, user)
}
