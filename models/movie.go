package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Movie struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title" validate:"required,min=1"`
	Description string             `bson:"description" json:"description" validate:"required,min=1"`
	Director    string             `bson:"director" json:"director" validate:"required,min=1"`
	ReleaseYear int                `bson:"release_year" json:"release_year" validate:"required,gte=1888,lte=2100"`
	Genre       string             `bson:"genre" json:"genre" validate:"required,min=1"`
}

