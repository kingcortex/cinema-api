package middleware

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kingcortex/cinema-api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TokenData struct {
	UserID primitive.ObjectID
}

// ValidateToken valide le token JWT et retourne le UserID extrait
func ValidateToken(tokenString string) (*TokenData, error) {
	// Parse le token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Vérifie expiration
	exp, ok := claims["exp"].(float64)
	if !ok || int64(exp) < time.Now().Unix() {
		return nil, errors.New("token expired")
	}

	// Récupère l'ID utilisateur
	userIdStr, ok := claims["sub"].(string)
	if !ok {
		return nil, errors.New("invalid user ID in token")
	}

	objID, err := primitive.ObjectIDFromHex(userIdStr)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	return &TokenData{
		UserID: objID,
	}, nil
}

func AuthMiddleware(authService *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Authorization header missing")
			c.Abort()
			return
		}

		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		// Valide le token
		tokenData, err := ValidateToken(tokenString)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		// Vérifie si le user existe
		var user bson.M
		err = authService.FindOne(c, bson.M{"_id": tokenData.UserID}).Decode(&user)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "User not found")
			c.Abort()
			return
		}

		// Passe le userId dans le contexte
		c.Set("userId", tokenData.UserID.Hex())

		c.Next()
	}
}
