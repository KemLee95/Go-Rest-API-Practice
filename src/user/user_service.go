package user

import (
	"context"
	"encoding/json"
	"errors"

	encryptionService "github.com/kemlee/go-rest-api-practise/core/encryption"
	hashService "github.com/kemlee/go-rest-api-practise/core/hash"

	userRepository "github.com/kemlee/go-rest-api-practise/user/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateUserRequest struct {
	Email    string
	Name     string
	Password string
}

type UserAccessTokenPayload struct {
	ID    primitive.ObjectID `json:"id"`
	Email string             `json:"email"`
	Name  string             `json:"name"`
}

type IUserService interface {
	GetUserById(ctx context.Context, id string) (*userRepository.UserModel, error)
	GetUserList(ctx context.Context) ([]*userRepository.UserModel, error)
	CreateUser(ctx context.Context, data *CreateUserRequest) error
	SignIn(ctx context.Context, userName string, password string) (string, error)
	CheckEmailExist(ctx context.Context, email string) (bool, error)
	CheckIdExist(ctx context.Context, id string) (bool, error)
}

type userService struct {
	userRepo      userRepository.IUserRepository
	hashSer       hashService.IHashService
	encryptionSer encryptionService.IEncryptionService
}

var (
	userServiceInstance IUserService
)

func GetUserService(
	userRepo userRepository.IUserRepository,
	hashSer hashService.IHashService,
	encryptionSer encryptionService.IEncryptionService,
) IUserService {
	if userServiceInstance == nil {
		userServiceInstance = &userService{
			userRepo:      userRepo,
			hashSer:       hashSer,
			encryptionSer: encryptionSer,
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

func (ser *userService) SignIn(ctx context.Context, email, password string) (string, error) {
	user, err := ser.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("The email or password incorrect")
	}
	checkPassword := ser.hashSer.BcryptCompare(password, user.Password)
	if !checkPassword {
		return "", errors.New("The email or password incorrect")
	}
	userAccessTokenPayload := &UserAccessTokenPayload{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}
	userAccessTokenPayloadJson, _ := json.Marshal(userAccessTokenPayload)
	token := ser.encryptionSer.JWTEncrypt(userAccessTokenPayloadJson, encryptionService.JWTEncryptOptions{})

	return token, nil
}

func (ser *userService) CheckEmailExist(ctx context.Context, email string) (bool, error) {
	return ser.userRepo.CheckEmailExist(ctx, email)
}

func (ser *userService) CheckIdExist(ctx context.Context, id string) (bool, error) {
	return ser.userRepo.CheckIdExist(ctx, id)
}

func (ser *userService) GetUserById(ctx context.Context, id string) (*userRepository.UserModel, error) {
	return ser.userRepo.GetUserById(ctx, id)
}

func (ser *userService) GetUserList(ctx context.Context) ([]*userRepository.UserModel, error) {
	return ser.userRepo.GetUserList(ctx)
}
