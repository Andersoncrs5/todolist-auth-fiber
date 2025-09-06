package userDto

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserDTO struct {
	ID           primitive.ObjectID `json:"id,omitempty"`
	Username     string             `json:"username" `
	Email        string             `json:"email" `
	CreatedAt    *time.Time         `json:"created_at" `
}