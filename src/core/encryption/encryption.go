package encryption

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type encryptionService struct{}

type IEncryptionService interface {
	JWTEncrypt(payload []byte, options JWTEncryptOptions) (string, error)
}

type AppClaims struct {
	payload []byte
	jwt.RegisteredClaims
}
type JWTEncryptOptions struct {
	SecretKey string
	ExpireIn  int
	Issuer    string
	Subject   string
}

func New() *encryptionService {
	return &encryptionService{}
}

func (ser *encryptionService) JWTEncrypt(payload []byte, options JWTEncryptOptions) (string, error) {
	claims := &AppClaims{
		payload,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(options.ExpireIn))),
			Issuer:    options.Issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   options.Subject,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	key := []byte(options.SecretKey)
	signedToken, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
