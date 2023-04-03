package user_repository

import (
	"context"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IUserRepository interface {
	GetUserById(ctx context.Context, id string) (*UserModel, error)
	GetUserByEmail(ctx context.Context, email string) (*UserModel, error)
	GetUserList(ctx context.Context) ([]*UserModel, error)
	CreateUser(ctx context.Context, user *UserModel) error
	FindUserByEmail(ctx context.Context, email string) (UserModel, error)
	CheckEmailExist(ctx context.Context, email string) (bool, error)
	CheckIdExist(ctx context.Context, id string) (bool, error)
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

func (repo *userRepository) CreateUser(ctx context.Context, user *UserModel) error {
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

func (repo *userRepository) FindUserByEmail(ctx context.Context, email string) (UserModel, error) {
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

func (repo *userRepository) CheckEmailExist(ctx context.Context, email string) (bool, error) {
	filters := bson.M{email: email}
	counter, err := repo.db.CountDocuments(ctx, filters)
	if err != nil {
		return true, err
	} else if counter > 0 {
		return true, nil
	}
	return false, nil
}

func (repo *userRepository) CheckIdExist(ctx context.Context, id string) (bool, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	filters := bson.M{"_id": _id}
	counter, err := repo.db.CountDocuments(ctx, filters)
	if err != nil {
		return true, err
	} else if counter > 0 {
		return true, nil
	}
	return false, nil
}

func (repo *userRepository) GetUserById(ctx context.Context, id string) (*UserModel, error) {
	user := &UserModel{}
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return user, err
	}

	filters := bson.M{
		"_id": _id,
	}
	options := options.FindOne().SetProjection(bson.D{
		{Key: "email", Value: 1},
		{Key: "name", Value: 1},
	})
	err = repo.db.FindOne(ctx, filters, options).Decode(user)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (repo *userRepository) GetUserByEmail(ctx context.Context, email string) (*UserModel, error) {
	user := &UserModel{}
	filters := bson.M{
		"email": email,
	}
	options := options.FindOne().SetProjection(bson.D{
		{Key: "email", Value: 1},
		{Key: "name", Value: 1},
		{Key: "password", Value: 1},
	})
	err := repo.db.FindOne(ctx, filters, options).Decode(user)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (repo *userRepository) GetUserList(ctx context.Context) ([]*UserModel, error) {
	options := options.Find().SetProjection(bson.D{
		{Key: "email", Value: 1},
		{Key: "name", Value: 1},
	})
	cur, err := repo.db.Find(ctx, bson.D{}, options)
	if err != nil {
		return []*UserModel{}, err
	}
	defer cur.Close(ctx)
	var users []*UserModel
	for cur.Next(ctx) {
		user := &UserModel{}
		err := cur.Decode(user)
		if err != nil {
			return []*UserModel{}, err
		}
		users = append(users, user)
	}
	return users, nil
}
