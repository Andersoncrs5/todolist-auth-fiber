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

func GenerateJWT(user models.User) (string, error)  {
	if err := godotenv.Load(); err != nil {
		return "", fmt.Errorf("error loading .env file: %v", err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", fmt.Errorf("JWT_SECRET not set in .env file")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.Hex(),
		"email": user.Email,
		"username": user.Username,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}
	return tokenString, nil
}

func GenerateRefreshToken(user models.User) (string, error) {
	if err := godotenv.Load(); err != nil {
		return "", fmt.Errorf("error loading .env file: %v", err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", fmt.Errorf("JWT_SECRET not set in .env file")
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.Hex(),
		"email": user.Email,
		"username": user.Username,
		"exp": time.Now().Add(time.Hour * (24 * 7)).Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return refreshTokenString, nil
}

func VerifyJWT(tokenString string) (primitive.ObjectID, error) {
	if err := godotenv.Load(); err != nil {
		return primitive.NilObjectID, fmt.Errorf("error loading .env file: %v", err)
	}
	
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return primitive.NilObjectID, fmt.Errorf("JWT_SECRET not set in .env file")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("failed to parse token: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDHex, ok := claims["user_id"].(string)
		if !ok {
			return primitive.NilObjectID, fmt.Errorf("user_id not found in token claims")
		}
		userID, err := primitive.ObjectIDFromHex(userIDHex)
		if err != nil {
			return primitive.NilObjectID, fmt.Errorf("invalid user_id in token claims: %v", err)
		}
		return userID, nil
	} else {
		return primitive.NilObjectID, fmt.Errorf("invalid token")
	}
}