package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Todo struct {
	ID     primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"user_id" bson:"user_id"`
	Title  string             `json:"title" bson:"title"`
	Discription  string        `json:"discription" bson:"discription"`
	Done   bool               `json:"done" bson:"done"`
}