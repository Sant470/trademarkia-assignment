package store

import "github.com/sant470/trademark/dtos"

type Store interface {
	CheckUser(username string) (bool, error)
	AddUser(user *dtos.RegisterRequest) error
}
