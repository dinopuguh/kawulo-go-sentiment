package services

import (
	"log"

	"github.com/dinopuguh/kawulo-go-sentiment/database"
	"github.com/dinopuguh/kawulo-go-sentiment/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CheckSentimentExist(db *mongo.Database, rev_id string) bool {
	ctx := database.Ctx

	var result models.Sentiment
	err := db.Collection("sentiment").FindOne(ctx, bson.M{"review_id": rev_id}).Decode(&result)
	if err != nil {
		return false
	}

	return true
}

func InsertSentiments(db *mongo.Database, loc models.Location) {
	rests := FindRestaurantByLocation(db, loc.LocationId)

	for _, rest := range rests {
		log.Print(rest.Name)
	}
}
