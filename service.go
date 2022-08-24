package genericscrud

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoService[T any] struct {
	db             *mongo.Database
	collectionName string
}

func (s *MongoService[T]) Find(ctx context.Context, id string) (*Model[T], error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid object id format: %w", err)
	}

	var m *Model[T]
	if err := s.db.Collection(s.collectionName).FindOne(ctx, bson.M{"_id": objID}).Decode(&m); err != nil {
		return nil, err
	}

	return m, nil
}

func (s *MongoService[T]) Create(ctx context.Context, create *T) (*Model[T], error) {
	now := time.Now()

	m := &Model[T]{
		CreatedAt: now,
		UpdatedAt: now,
		Data:      *create,
	}

	result, err := s.db.Collection(s.collectionName).InsertOne(ctx, m)
	if err != nil {
		return nil, err
	}

	m.ID = result.InsertedID.(primitive.ObjectID)

	return m, nil
}

func (s *MongoService[T]) Update(ctx context.Context, id string, update *T) (*Model[T], error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid object id format: %w", err)
	}

	updateBSON := &updateBSON[*T]{Data: update, UpdatedAt: time.Now()}

	_, err = s.db.Collection(s.collectionName).UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": updateBSON})
	if err != nil {
		return nil, err
	}

	return s.Find(ctx, id)
}

func (s *MongoService[T]) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid object id format: %w", err)
	}

	_, err = s.db.Collection(s.collectionName).DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}

	return nil
}

type updateBSON[T any] struct {
	Data      T         `bson:"inline"`
	UpdatedAt time.Time `bson:"updated_at"`
}
