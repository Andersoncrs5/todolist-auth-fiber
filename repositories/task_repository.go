package repository

import (
	"context"
	"errors"
	"fmt"
	"time"
	taskdto "todolist-auth-fiber/dtos/taskDto"
	"todolist-auth-fiber/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TaskRepository interface {
	GetById(ctx context.Context, id primitive.ObjectID) (*models.Todo, int, error)
	Create(ctx context.Context, userID primitive.ObjectID, task models.Todo) (*models.Todo, int, error)
	Delete(ctx context.Context, id primitive.ObjectID) (int, error)
	ChangeStatus(ctx context.Context, id primitive.ObjectID, task models.Todo) (*models.Todo, int, error)
	Update(ctx context.Context, id primitive.ObjectID, dto taskdto.UpdateTaskDTO) (*models.Todo, int, error)
	GetAll(ctx context.Context,userID primitive.ObjectID,title string,done *bool,createdAtBefore, createdAtAfter time.Time,page, pageSize int) ([]models.Todo, int64, error)
	DeleteAllByUserId(ctx context.Context, userId primitive.ObjectID) (int64, error)
}

type taskRepository struct {
	collection *mongo.Collection
}

func NewTaskRepository(db *mongo.Database) UserRepository {
	return &userRepository{
		collection: db.Collection("tasks"),
	}
}

func (r *taskRepository) GetById(ctx context.Context, id primitive.ObjectID) (*models.Todo, int, error) {
	var task models.Todo
	filter := bson.M{"_id": id}

	err := r.collection.FindOne(ctx, filter).Decode(&task)
	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		return nil, 404, fmt.Errorf("Task not found")
	}

	if err != nil {
		return nil, 500, fmt.Errorf("Error the get tasks by id! Error: %w", err)
	}

	return &task, 200, nil
}

func (r *taskRepository) Create(ctx context.Context, userID primitive.ObjectID, task models.Todo) (*models.Todo, int, error) {
	task.ID = primitive.NewObjectID()
	now := time.Now()

	task.Done = false
	task.UserID = userID
	task.CreatedAt = &now
	task.UpdatedAt = &now

	if _, err := r.collection.InsertOne(ctx, task); err != nil {
		return nil, 500, fmt.Errorf("Error the save task in database %w", err)
	}

	return &task, 201, nil
}

func (r *taskRepository) Delete(ctx context.Context, id primitive.ObjectID) (int, error) {
	filter := bson.M{"_id": id}
	result, err := r.collection.DeleteOne(ctx, filter)

	if err != nil {
		return 500, fmt.Errorf("Error the delete task\nError: %w", err)
	}

	if result.DeletedCount == 0 {
		return 500, errors.New("Task not deleted")
	}

	return 200, nil
}

func (r *taskRepository) ChangeStatus(ctx context.Context, id primitive.ObjectID, task models.Todo) (*models.Todo, int, error) {
	base := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "done", Value: task.Done},
		}},
	}

	opts := options.FindOneAndUpdate().SetProjection(options.After)
	var taskUpdated models.Todo

	err := r.collection.FindOneAndUpdate(ctx, bson.M{"_id": id}, base, opts).Decode(&taskUpdated)

	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		return nil, 404, fmt.Errorf("Task not found")
	}

	if err != nil {
		return nil, 500, fmt.Errorf("Error the change status tasks by id!\nError: %w", err)
	}

	return &taskUpdated, 200, nil
}

func (r *taskRepository) Update(ctx context.Context, id primitive.ObjectID, dto taskdto.UpdateTaskDTO) (*models.Todo, int, error) {
	base := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "done", Value: dto.Done},
			{Key: "discription", Value: dto.Discription},
			{Key: "title", Value: dto.Title},
		}},
	}

	opts := options.FindOneAndUpdate().SetProjection(options.After)
	var taskUpdated models.Todo

	err := r.collection.FindOneAndUpdate(ctx, bson.M{"_id": id}, base, opts).Decode(&taskUpdated)

	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		return nil, 404, fmt.Errorf("Task not found")
	}

	if err != nil {
		return nil, 500, fmt.Errorf("Error the to update tasks by id!\nError: %w", err)
	}

	return &taskUpdated, 200, nil
}

func (r *taskRepository) GetAll(
	ctx context.Context,
	userID primitive.ObjectID,
	title string,
	done *bool,
	createdAtBefore, createdAtAfter time.Time,
	page, pageSize int,
) ([]models.Todo, int64, error) {

	filter := bson.M{"user_id": userID}

	if title != "" {
		filter["title"] = bson.M{"$regex": title, "$options": "i"}
	}

	if done != nil {
		filter["done"] = *done
	}

	if !createdAtBefore.IsZero() {
		filter["created_at"] = bson.M{"$lte": createdAtBefore}
	}

	if !createdAtAfter.IsZero() {
		if _, ok := filter["created_at"]; ok {
			filter["created_at"].(bson.M)["$gte"] = createdAtAfter
		} else {
			filter["created_at"] = bson.M{"$gte": createdAtAfter}
		}
	}

	skip := int64((page - 1) * pageSize)
	limit := int64(pageSize)

	findOptions := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}

	defer cursor.Close(ctx)

	var tasks []models.Todo
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, 0, err
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

func (r *taskRepository) DeleteAllByUserId(ctx context.Context, userId primitive.ObjectID) (int64, error) {
	filter := bson.M{"user_id": userId}
	result, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}