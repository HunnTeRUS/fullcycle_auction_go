package user_entity

import (
	"context"
	"github.com/HunnTeRUS/fullcycle-auction-go/internal/internal_error"
)

type User struct {
	Id   string `copier:"Id"`
	Name string `copier:"Name"`
}

func (u *User) Validate() *internal_error.InternalError {
	if len(u.Name) <= 1 {
		return internal_error.NewBadRequestError("name is not a valid field value")
	}

	return nil
}

type UserRepositoryInterface interface {
	FindUserByID(ctx context.Context, userId string) (*User, *internal_error.InternalError)
}
