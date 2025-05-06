package repositories

import (
	"context"
	"log"

	"github.com/kingcortex/cinema-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MovieRepository struct {
	MongoCollection mongo.Collection
}

func (r *MovieRepository) InsertOneMovie(movie *models.Movie) (interface{}, error) {
	result, err := r.MongoCollection.InsertOne(context.Background(), movie)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}

func (r *MovieRepository) FindOneMovieByID(movieID primitive.ObjectID) (*models.Movie, error) {
	var movie models.Movie
	err := r.MongoCollection.FindOne(context.Background(), bson.M{"_id": movieID}).Decode(&movie)

	if err != nil {
		return nil, err
	}

	return &movie, nil

}

func (r *MovieRepository) FindAllMovie() (*[]models.Movie, error) {
	var movies []models.Movie
	cursor, err := r.MongoCollection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}

	err = cursor.All(context.Background(), &movies)

	if err != nil {
		return nil, err
	}

	return &movies, nil

}

func (r *MovieRepository) UpdateMovieByID(objectID primitive.ObjectID, updateMovieBson *bson.M) (int64, error) {
	log.Println(updateMovieBson)
	result, err := r.MongoCollection.UpdateByID(context.Background(), objectID, *updateMovieBson)

	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}