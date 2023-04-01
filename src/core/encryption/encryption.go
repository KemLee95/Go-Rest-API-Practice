package encryption

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type encryptionService struct{}

type IEncryptionService interface {
	JWTEncrypt(payload []byte, options JWTEncryptOptions) string
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

func (ser *encryptionService) JWTEncrypt(payload []byte, options JWTEncryptOptions) string {
	claims := &AppClaims{
		payload,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(options.ExpireIn))),
			Issuer:    options.Issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   options.Subject,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	signedToken, err := token.SignedString(options.SecretKey)
	if err != nil {
		return ""
	}
	return signedToken
}
