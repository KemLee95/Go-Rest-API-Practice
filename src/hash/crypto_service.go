package hash

import (
	"crypto/md5"
	"io"
)

type IHashService interface {
	Bcrypt(value string) string
	BcryptCompare(value, hashedValue string) bool
	Sha256(value string) string
	Sha256Compare(value, hashedValue string) bool
}

type hashService struct{}

var (
	instance IHashService
)

func (ser *hashService) Bcrypt(value string) string {
	hash := md5.New()
	io.WriteString(hash, value)
	return string(hash.Sum(nil))
}

func (ser *hashService) BcryptCompare(value, hashedValue string) bool {
	return false
}

func (ser *hashService) Sha256(value string) string {
	return "123"
}

func (ser *hashService) Sha256Compare(value, hashedValue string) bool {
	return false
}

func GetHashService() IHashService {
	if instance == nil {
		instance = &hashService{}
	}
	return instance
}
