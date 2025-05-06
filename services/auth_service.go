package services

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kingcortex/cinema-api/config"
	"github.com/kingcortex/cinema-api/models"
	"github.com/kingcortex/cinema-api/repositories"
	"github.com/kingcortex/cinema-api/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	MongoCollection *mongo.Collection
}

const PasswordCost = 10

func (svc *AuthService) Register(c *gin.Context) {
	type Payload struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}
	var payload Payload

	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.ErrorResponse(c, 400, "Invalid JSON format")
		return
	}
	err := validate.Struct(&payload)
	if err != nil {
		log.Println("decode error", err)
		utils.ErrorResponse(c, 400, err.Error())
		return
	}
	repo := repositories.UserRepository{MongoCollection: svc.MongoCollection}
	existUser, err := repo.FindUserByEmail(payload.Email)
	if err != nil && err != mongo.ErrNoDocuments {
		utils.ErrorResponse(c, 500, "Database error")
		return
	}
	if existUser != nil {
		utils.ErrorResponse(c, 400, "Email already in use")
		return
	}
	password, err := bcrypt.GenerateFromPassword([]byte(payload.Password), PasswordCost)

	if err != nil {
		log.Println("failed to hash password:", err)
		utils.ErrorResponse(c, 500, "Internal server error")
		return
	}

	user := models.User{ID: primitive.NewObjectID(), Email: payload.Email, Password: string(password)}

	result, err := repo.InsertUser(&user)
	if err != nil {
		log.Println("save error", err)
		utils.ErrorResponse(c, 500, err.Error())
		return
	}
	utils.SuccessResponse(c, map[string]any{
		"message": "User registered successfully",
		"user":    result,
	})
}

func (svc *AuthService) Login(c *gin.Context) {
	type Payload struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}
	var payload Payload

	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.ErrorResponse(c, 400, "Invalid JSON format")
		return
	}
	err := validate.Struct(&payload)
	if err != nil {
		log.Println("decode error", err)
		utils.ErrorResponse(c, 400, err.Error())
		return
	}
	repo := repositories.UserRepository{MongoCollection: svc.MongoCollection}
	user, _ := repo.FindUserByEmail(payload.Email)
	if user == nil {
		utils.ErrorResponse(c, 401, "invalid email or password")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))

	if err != nil {
		utils.ErrorResponse(c, 401, "invalid email or password")
		return
	}

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID.Hex(),
		"exp": time.Now().Add(time.Hour * 24 * 60).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(config.MyConfig.JWT_SECRET))

	if err != nil {
		log.Println(err)
		utils.ErrorResponse(c, 400, "Failed to create token")
		return
	}

	utils.SuccessResponse(c, map[string]string{"token": tokenString})

}
