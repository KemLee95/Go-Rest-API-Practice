package repository

import "go-practise/model"

type UserRepo interface {
	Insert(u model.User) error
	FindUserByEmail(email string) (model.User, error)
	CheckUserLogin(email, password string) (model.User, error)
}
