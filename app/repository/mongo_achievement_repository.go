package repository

import (
	"context"
	"time"

	"github.com/Lutfania/ekrp/app/models"
	"github.com/Lutfania/ekrp/database" // pastikan path sesuai
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoAchievementRepository struct{}

func NewMongoAchievementRepository() *MongoAchievementRepository {
	return &MongoAchievementRepository{}
}

func (r *MongoAchievementRepository) Insert(doc *models.MongoAchievement) (string, error) {
	coll := database.Collection("achievements")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	doc.CreatedAt = time.Now()
	res, err := coll.InsertOne(ctx, doc)
	if err != nil {
		return "", err
	}
	oid := res.InsertedID.(primitive.ObjectID)
	return oid.Hex(), nil
}

func (r *MongoAchievementRepository) FindByIDHex(hexID string) (*models.MongoAchievement, error) {
	oid, err := primitive.ObjectIDFromHex(hexID)
	if err != nil {
		return nil, err
	}
	coll := database.Collection("achievements")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var doc models.MongoAchievement
	if err := coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&doc); err != nil {
		return nil, err
	}
	return &doc, nil
}

func (r *MongoAchievementRepository) UpdateByHex(hexID string, update bson.M) error {
	oid, err := primitive.ObjectIDFromHex(hexID)
	if err != nil {
		return err
	}

	coll := database.Collection("achievements")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = coll.UpdateByID(ctx, oid, update) // ⬅️ PENTING
	return err
}

func (r *MongoAchievementRepository) DeleteByHex(hexID string) error {
	oid, err := primitive.ObjectIDFromHex(hexID)
	if err != nil {
		return err
	}
	coll := database.Collection("achievements")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = coll.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}
