package repository

import (
	"context"
	"errors"
	"fmt"
	"time"
	"todolist-auth-fiber/dtos/userDto"
	"todolist-auth-fiber/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository interface {
	GetEmail(ctx context.Context, email string) (*models.User, error)
	GetId(ctx context.Context, id primitive.ObjectID) (*models.User, error)
	Save(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
	Update(ctx context.Context, id primitive.ObjectID, update userDto.UpdateUserDTO) (*models.User, uint, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByUserName(ctx context.Context, username string) (bool, error)
	SetRefreshToken(ctx context.Context, id primitive.ObjectID, refreshToken string) (*models.User, int, error)
}

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) UserRepository { 
	return &userRepository{
		collection: db.Collection("users"),
	}
}

func (u *userRepository) GetEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	filter := bson.M{"email": email}

	err := u.collection.FindOne(ctx, filter).Decode(&user);
	if err != nil  {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil ,nil;
		}
		return nil, fmt.Errorf("Fail the to search user by email");
	}

	return &user, nil
}

func (u *userRepository) GetId(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	var user models.User
	filter := bson.M{"_id": id}

	err := u.collection.FindOne(ctx, filter).Decode(&user);
	if err != nil  {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil ,nil;
		}

		return nil, fmt.Errorf("Fail the to search user by id");
	}

	return &user, nil
}

func (u *userRepository) Save(ctx context.Context, user *models.User) (*models.User, error) {
	user.ID = primitive.NewObjectID()
	now := time.Now()

    user.CreatedAt = &now
    user.UpdatedAt = &now

	_, err := u.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("Error the save user in database")
	}

	return user, nil
}

func (u *userRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{ "_id": id }
	result, err := u.collection.DeleteOne(ctx, filter)

	if err != nil {
		return fmt.Errorf("Error the delete user")
	}

	if result.DeletedCount == 0 {
		return errors.New("User not deleted")
	}

	return nil
}

func (u *userRepository) Update(ctx context.Context, id primitive.ObjectID, update userDto.UpdateUserDTO) (*models.User, uint, error) {
	base := bson.D{
		{Key : "$set", Value: bson.D{
			{Key: "username", Value: update.Username},
			{Key: "password", Value: update.Password},
			{Key: "updated_at", Value: time.Now()},
		}},
	}
	
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedUser models.User
	err := u.collection.FindOneAndUpdate(ctx, bson.M{"_id": id}, base, opts).Decode(&updatedUser)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, 404, fmt.Errorf("User not found")
		}
		return nil, 500, fmt.Errorf("fail to update user: %w", err)
	}

	return &updatedUser, 200, nil
}

func (u *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
    filter := bson.M{"email": email}
    opts := options.FindOne().SetProjection(bson.M{"_id": 1})

    var result struct {
        ID primitive.ObjectID `bson:"_id"`
    }

    err := u.collection.FindOne(ctx, filter, opts).Decode(&result)
    if err != nil {
        if errors.Is(err, mongo.ErrNoDocuments) {
            return false, nil
        }
        return false, fmt.Errorf("fail to check if user exists by email: %w", err)
    }

    return true, nil
}

func (u *userRepository) ExistsByUserName(ctx context.Context, username string) (bool, error) { 
	filter := bson.M{"username": username}
	opts := options.FindOne().SetProjection(bson.M{"_id": 1})

	var result struct {
        ID primitive.ObjectID `bson:"_id"`
    }

	err := u.collection.FindOne(ctx, filter, opts).Decode(&result)
	if err != nil {
        if errors.Is(err, mongo.ErrNoDocuments) {
            return false, nil
        }
        return false, fmt.Errorf("fail to check if user exists by username: %w", err)
    }

	return true, nil
}

func (u *userRepository) SetRefreshToken(ctx context.Context, id primitive.ObjectID, refreshToken string) (*models.User, int, error) {
	base := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "refresh_token", Value: refreshToken},
			{Key: "updated_at", Value: time.Now()},
		}},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var userUpdated models.User
	err := u.collection.FindOneAndUpdate(ctx, bson.M{"_id": id}, base, opts).Decode(&userUpdated)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, 404, fmt.Errorf("User not found")
		}
		return nil, 500,fmt.Errorf("fail to update user: %w", err)
	}

	return &userUpdated, 200, nil
}