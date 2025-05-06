package services

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/kingcortex/cinema-api/dtos"
	"github.com/kingcortex/cinema-api/helpers"
	"github.com/kingcortex/cinema-api/models"
	"github.com/kingcortex/cinema-api/repositories"
	"github.com/kingcortex/cinema-api/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MovieService struct {
	MongoCollection mongo.Collection
}

var validate = validator.New()

func (svc *MovieService) GetAll(c *gin.Context) {
	log.Println("MovieService.GetAll called")

	repo := repositories.MovieRepository{MongoCollection: svc.MongoCollection}

	movies, err := repo.FindAllMovie()

	if err != nil {
		log.Println("Error fetching movies:", err)
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, movies)
	log.Println("Successfully fetched movies")
}

func (svc *MovieService) CreateOneMovie(c *gin.Context) {
	log.Println("MovieService.CreateOneMovie called")

	body := c.Request.Body
	defer body.Close()
	var dto models.Movie

	c.ShouldBindJSON(&dto)
	err := validate.Struct(&dto)
	if err != nil {
		log.Println("Error decoding movie request body:", err)
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	repo := repositories.MovieRepository{MongoCollection: svc.MongoCollection}

	newMovie := models.Movie{
		ID:          primitive.NewObjectID(),
		Title:       dto.Title,
		Description: dto.Description,
		Director:    dto.Director,
		ReleaseYear: dto.ReleaseYear,
		Genre:       dto.Genre,
	}
	result, err := repo.InsertOneMovie(&newMovie)

	if err != nil {
		log.Println("Error inserting new movie:", err)
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, result)
	log.Println("Successfully inserted new movie")

}

func (svc *MovieService) UpdateMovieByID(c *gin.Context) {
	id := c.Param("id")

	// Convertir l'ID en ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid ID format:", err)
		utils.ErrorResponse(c, 400, "Invalid ID format")
		return
	}

	// Lire le corps de la requête
	var dto dtos.UpdateMovieDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		log.Println("Error decoding movie request body:", err)
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	// Valider les données
	if err := validate.Struct(&dto); err != nil {
		log.Println("Validation error:", err)
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	// Mettre à jour le film
	repo := repositories.MovieRepository{MongoCollection: svc.MongoCollection}
	updateData := helpers.BuildUpdateBson(dto)

	result, err := repo.UpdateMovieByID(objectID, &updateData)
	if err != nil {
		log.Println("Error updating movie:", err)
		utils.ErrorResponse(c, 500, err.Error())
		return
	}

	utils.SuccessResponse(c, fmt.Sprintln("Movie updated successfully", result))
	log.Println("Successfully updated movie")
}

func (svc *MovieService) GetOneMovie(c *gin.Context) {
	id := c.Param("id")

	// Convertir l'ID en ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid ID format:", err)
		utils.ErrorResponse(c, 400, "Invalid ID format")
		return
	}

	repo := repositories.MovieRepository{MongoCollection: svc.MongoCollection}

	result, err := repo.FindOneMovieByID(objectID)

	if err != nil {
		log.Println("Error Getting movie:", err)
		utils.ErrorResponse(c, 500, err.Error())
		return
	}

	utils.SuccessResponse(c, result)
}
