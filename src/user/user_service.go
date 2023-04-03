package user

import (
	"context"
	"encoding/json"
	"errors"

	config "github.com/kemlee/go-rest-api-practise/config"
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

type LoginPayload struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpireIn     int    `json:"expire_in"`
}

type IUserService interface {
	GetUserById(ctx context.Context, id string) (*userRepository.UserModel, error)
	GetUserList(ctx context.Context) ([]*userRepository.UserModel, error)
	CreateUser(ctx context.Context, data *CreateUserRequest) error
	SignIn(ctx context.Context, userName string, password string) (*LoginPayload, error)
	CheckEmailExist(ctx context.Context, email string) (bool, error)
	CheckIdExist(ctx context.Context, id string) (bool, error)
}

type userService struct {
	userRepo      userRepository.IUserRepository
	hashSer       hashService.IHashService
	encryptionSer encryptionService.IEncryptionService
	config        *config.ApiConfig
}

var (
	userServiceInstance IUserService
)

func GetUserService(
	userRepo userRepository.IUserRepository,
	hashSer hashService.IHashService,
	encryptionSer encryptionService.IEncryptionService,
	config *config.ApiConfig,
) IUserService {
	if userServiceInstance == nil {
		userServiceInstance = &userService{
			userRepo:      userRepo,
			hashSer:       hashSer,
			encryptionSer: encryptionSer,
			config:        config,
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

func (ser *userService) SignIn(ctx context.Context, email, password string) (*LoginPayload, error) {
	result := &LoginPayload{}
	user, err := ser.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return result, err
	}
	if user == nil {
		return result, errors.New("The email or password incorrect")
	}
	checkPassword := ser.hashSer.BcryptCompare(password, user.Password)
	if !checkPassword {
		return result, errors.New("The email or password incorrect")
	}
	userAccessTokenPayload := &UserAccessTokenPayload{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}
	ser.getJWTToken(result, userAccessTokenPayload, "")

	return result, nil
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

func (ser *userService) getJWTToken(result *LoginPayload, payload *UserAccessTokenPayload, refreshToken string) {
	jwtConfig := ser.config.JWTConfig
	userAccessTokenPayloadJson, _ := json.Marshal(*payload)

	token, _ := ser.encryptionSer.JWTEncrypt(userAccessTokenPayloadJson, encryptionService.JWTEncryptOptions{
		SecretKey: jwtConfig.PrivateKey,
		ExpireIn:  jwtConfig.AccessTokenExpiration,
		Issuer:    payload.Email,
		Subject:   "App",
	})
	if refreshToken == "" {
		refreshToken, _ = ser.encryptionSer.JWTEncrypt(userAccessTokenPayloadJson, encryptionService.JWTEncryptOptions{
			SecretKey: jwtConfig.PrivateKey,
			ExpireIn:  jwtConfig.RefreshTokenExpiration,
			Issuer:    payload.Email,
			Subject:   "App",
		})
	}
	result.AccessToken = token
	result.RefreshToken = refreshToken
	result.ExpireIn = jwtConfig.AccessTokenExpiration
}
