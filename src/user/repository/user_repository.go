package user_repository

import (
	"context"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IUserRepository interface {
	CreateUser(ctx context.Context, user *UserModel) error
	FindUserByEmail(ctx context.Context, email string) (UserModel, error)
	CheckEmailExist(ctx context.Context, email string) (bool, error)
}

type userRepository struct {
	db *mongo.Collection
}

var (
	instance *userRepository
	once     sync.Once
)

func GetUserRepository(db *mongo.Database) *userRepository {
	once.Do(func() {
		instance = &userRepository{
			db: db.Collection("users"),
		}
	})
	return instance
}

func (repo userRepository) CreateUser(ctx context.Context, user *UserModel) error {
	bTypes, err := bson.Marshal(user)
	if err != nil {
		return err
	}
	_, err = repo.db.InsertOne(ctx, bTypes)
	if err != nil {
		return err
	}
	return nil
}

func (repo userRepository) FindUserByEmail(ctx context.Context, email string) (UserModel, error) {
	user := UserModel{}
	filter := bson.M{email: email}
	options := options.FindOne().SetProjection(
		bson.D{
			{Key: "email", Value: true},
			{Key: "name", Value: true},
		},
	)
	err := repo.db.FindOne(ctx, filter, options).Decode(user)

	if err != nil {
		return user, err
	}
	return user, nil
}

func (repo userRepository) CheckEmailExist(ctx context.Context, email string) (bool, error) {
	filters := bson.M{email: email}
	counter, err := repo.db.CountDocuments(ctx, filters)
	if err != nil {
		return true, err
	} else if counter > 0 {
		return true, nil
	}
	return false, nil
}
