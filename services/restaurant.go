package services

import (
	"log"

	"github.com/dinopuguh/kawulo-go-sentiment/database"
	"github.com/dinopuguh/kawulo-go-sentiment/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindAllRestaurants(db *mongo.Database) []models.Restaurant {
	ctx := database.Ctx

	csr, err := db.Collection("restaurant").Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err.Error())
	}
	defer csr.Close(ctx)

	result := make([]models.Restaurant, 0)
	for csr.Next(ctx) {
		var row models.Restaurant
		err := csr.Decode(&row)
		if err != nil {
			log.Fatal(err.Error())
		}

		result = append(result, row)
	}

	return result
}

func FindRestaurantByLocation(db *mongo.Database, locId string) []models.Restaurant {
	ctx := database.Ctx

	csr, err := db.Collection("restaurant").Find(ctx, bson.M{"locationID": locId})
	if err != nil {
		log.Fatal(err.Error())
	}
	defer csr.Close(ctx)

	result := make([]models.Restaurant, 0)
	for csr.Next(ctx) {
		var row models.Restaurant
		err := csr.Decode(&row)
		if err != nil {
			log.Fatal(err.Error())
		}

		result = append(result, row)
	}

	return result
}
