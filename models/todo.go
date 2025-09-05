package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Todo struct {
	ID     primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"user_id" bson:"user_id"`
	Title  string             `json:"title" bson:"title"`
	Discription  string       `json:"discription" bson:"discription"`
	Done   bool               `json:"done" bson:"done"`
	CreatedAt    *time.Time   `json:"created_at" bson:"created_at"`
	UpdatedAt    *time.Time   `json:"updated_at" bson:"updated_at"`
}