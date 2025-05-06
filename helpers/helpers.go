// helpers/update_helper.go
package helpers

import (
	"github.com/kingcortex/cinema-api/dtos"
	"go.mongodb.org/mongo-driver/bson"
)

func BuildUpdateBson(movieDto dtos.UpdateMovieDto) bson.M {
	updateFields := bson.M{}

	if movieDto.Title != nil {
		updateFields["title"] = *movieDto.Title
	}
	if movieDto.Description != nil {
		updateFields["description"] = *movieDto.Description
	}
	if movieDto.Director != nil {
		updateFields["director"] = *movieDto.Director
	}
	if movieDto.ReleaseYear != nil {
		updateFields["release_year"] = *movieDto.ReleaseYear
	}
	if movieDto.Genre != nil {
		updateFields["genre"] = *movieDto.Genre
	}

	return bson.M{
		"$set": updateFields,
	}
}