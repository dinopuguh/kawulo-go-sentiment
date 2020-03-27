package services

import (
	"log"

	"github.com/dinopuguh/kawulo-go-sentiment/database"
	"github.com/dinopuguh/kawulo-go-sentiment/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindAllReviews(db *mongo.Database) []models.Review {
	ctx := database.Ctx

	crs, err := db.Collection("review").Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err.Error())
	}
	defer crs.Close(ctx)

	result := make([]models.Review, 0)
	for crs.Next(ctx) {
		var row models.Review
		err := crs.Decode(&row)
		if err != nil {
			log.Fatal(err.Error())
		}

		result = append(result, row)
	}

	return result
}

func FindReviewByRestaurant(db *mongo.Database, restObjId primitive.ObjectID) []models.Review {
	ctx := database.Ctx

	csr, err := db.Collection("review").Find(ctx, bson.M{"restaurant_ObjectId": restObjId})
	if err != nil {
		log.Fatal(err.Error())
	}
	defer csr.Close(ctx)

	result := make([]models.Review, 0)
	for csr.Next(ctx) {
		var row models.Review
		err := csr.Decode(&row)
		if err != nil {
			log.Fatal(err.Error())
		}

		result = append(result, row)
	}

	return result
}
