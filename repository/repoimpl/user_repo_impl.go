package repository

import (
	"context"
	"go-practise/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepoImpl struct {
	DB *mongo.Database
}

func NewUserRepo(db *mongo.Database) *UserRepoImpl {
	return &UserRepoImpl{
		DB: db,
	}
}

func (mongo *UserRepoImpl) Insert(user model.User) error {
	btypes, _ := bson.Marshal(user)
	_, err := mongo.DB.Collection("users").InsertOne(context.Background(), btypes)
	if err != nil {
		return err
	}
	return nil
}

func (mongo *UserRepoImpl) FindUserByEmail(email string) (model.User, error) {
	user := model.User{}
	filter := bson.M{
		"email": email,
	}
	opts := options.FindOne().SetProjection(bson.D{
		{Key: "email", Value: true},
		{Key: "name", Value: true},
	})
	result := mongo.DB.Collection("users").FindOne(context.Background(), filter, opts)
	err := result.Decode(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (mongo *UserRepoImpl) CheckUserLogin(email, password string) (model.User, error) {
	user := model.User{}
	filter := bson.M{
		"email":    email,
		"password": password,
	}
	opts := options.FindOne().SetProjection(bson.D{
		{Key: "email", Value: true},
		{Key: "name", Value: true},
	})
	result := mongo.DB.Collection("users").FindOne(context.Background(), filter, opts)
	err := result.Decode(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}
