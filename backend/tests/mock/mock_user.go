package mock

import (
	"backend/internal/user"

	"github.com/stretchr/testify/mock"
)

type UserServiceMock struct {
	mock.Mock
}

func (m *UserServiceMock) GetProfile(id uint) (*user.User, error) {
	args := m.Mock.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *UserServiceMock) UpdateProfile(id uint, input user.UpdateUserInput) error {
	args := m.Mock.Called(id, input)
	return args.Error(0)
}
