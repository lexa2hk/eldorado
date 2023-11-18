package services

import (
	"errors"
)

var (
	ErrNilTasksStorage = errors.New("the tasks storage could not be nil")
	ErrNilUsersStorage = errors.New("the users storage could not be nil")
)
