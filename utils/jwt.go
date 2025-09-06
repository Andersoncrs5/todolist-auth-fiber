package utils

import (
	"fmt"
	"os"
	"time"
	"todolist-auth-fiber/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var jwtSecret string

const (
	AccessTokenExpiration  = time.Hour * 24
	RefreshTokenExpiration = time.Hour * 24 * 7
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(fmt.Errorf("error loading .env file: %v", err))
	}

	jwtSecret = os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		panic("JWT_SECRET not set in .env file")
	}
}

func GenerateToken(user *models.User, expiration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID.Hex(),
		"email":    user.Email,
		"username": user.Username,
		"exp":      time.Now().Add(expiration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func GenerateAccessToken(user *models.User) (string, error) {
	return GenerateToken(user, AccessTokenExpiration)
}

func GenerateRefreshToken(user *models.User) (string, error) {
	return GenerateToken(user, RefreshTokenExpiration)
}

func parseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
}

func ExtractUserID(tokenString string) (primitive.ObjectID, error) {
	token, err := parseToken(tokenString)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("failed to parse token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return primitive.NilObjectID, fmt.Errorf("invalid token")
	}

	userIDHex, ok := claims["user_id"].(string)
	if !ok {
		return primitive.NilObjectID, fmt.Errorf("user_id not found in token claims")
	}

	userID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("invalid user_id in token claims: %v", err)
	}

	return userID, nil
}
