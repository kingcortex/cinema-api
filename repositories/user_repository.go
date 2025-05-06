package repositories

import (
	"context"

	"github.com/kingcortex/cinema-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	MongoCollection *mongo.Collection
}

func (r *UserRepository) InsertUser(user *models.User) (interface{}, error) {
	result, err := r.MongoCollection.InsertOne(context.Background(), *user)

	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}

func (r *UserRepository) FindUserByEmail(email string) (*models.User, error) {
	result := r.MongoCollection.FindOne(context.Background(), bson.D{{Key: "email", Value: email}})

	var user models.User
	err := result.Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
