package user_usecase

import (
	"context"
	"github.com/HunnTeRUS/fullcycle-auction-go/internal/entity/user_entity"
	"github.com/HunnTeRUS/fullcycle-auction-go/internal/internal_error"
	"github.com/jinzhu/copier"
)

type UserUseCase struct {
	UserRepository user_entity.UserRepositoryInterface
}

type UserOutputDTO struct {
	Id   string `json:"id,omitempty" copier:"Id"`
	Name string `json:"name,omitempty" copier:"Name"`
}

type UserUseCaseInterface interface {
	FindUserByID(ctx context.Context, userID string) (*UserOutputDTO, *internal_error.InternalError)
}

func NewUserService(userRepository user_entity.UserRepositoryInterface) UserUseCaseInterface {
	return &UserUseCase{UserRepository: userRepository}
}

func (us *UserUseCase) FindUserByID(ctx context.Context, userID string) (*UserOutputDTO, *internal_error.InternalError) {
	userEntity, err := us.UserRepository.FindUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var userOutputDTO UserOutputDTO
	if err := copier.Copy(&userOutputDTO, userEntity); err != nil {
		return nil, internal_error.NewInternalServerError(err.Error())
	}

	return &userOutputDTO, nil
}
