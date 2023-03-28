package user

import (
	"context"

	hashService "github.com/kemlee/go-rest-api-practise/hash"
	userRepository "github.com/kemlee/go-rest-api-practise/user/repository"
)

type CreateUserRequest struct {
	Email    string
	Name     string
	Password string
}

type IUserService interface {
	CreateUser(ctx context.Context, data *CreateUserRequest) error
	SignIn(ctx context.Context, userName string, password string)
	CheckEmailExist(ctx context.Context, email string) (bool, error)
}

type userService struct {
	userRepo userRepository.IUserRepository
	hashSer  hashService.IHashService
}

var (
	userServiceInstance IUserService
)

func GetUserService(
	userRepo userRepository.IUserRepository,
	hashSer hashService.IHashService,
) IUserService {
	if userServiceInstance == nil {
		userServiceInstance = &userService{
			userRepo: userRepo,
			hashSer:  hashSer,
		}
	}
	return userServiceInstance
}

func (ser *userService) CreateUser(ctx context.Context, data *CreateUserRequest) error {
	userModel := &userRepository.UserModel{
		Email:    data.Email,
		Name:     data.Name,
		Password: ser.hashSer.Bcrypt(data.Password),
	}

	return ser.userRepo.CreateUser(ctx, userModel)
}

func (ser *userService) SignIn(ctx context.Context, email, password string) {

}

func (ser *userService) CheckEmailExist(ctx context.Context, email string) (bool, error) {
	return ser.userRepo.CheckEmailExist(ctx, email)
}
