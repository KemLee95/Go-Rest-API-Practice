package user_repository

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserModel struct {
	ID       primitive.ObjectID `bson:"_id, omitempty"`
	Email    string             `bson:"email"`
	Name     string             `bson:"name"`
	Password string             `bson:"password"`
}
