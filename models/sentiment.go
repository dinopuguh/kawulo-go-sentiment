package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Sentiment struct {
	ID             primitive.ObjectID `bson:"_id"`
	LocationId     string             `bson:"location_id"`
	Location       Location           `bson:"location"`
	RestaurantId   string             `bson:"restaurant_id"`
	Restaurant     Restaurant         `bson:"restaurant"`
	ReviewId       string             `bson:"review_id"`
	Review         Review             `bson:"review"`
	PublishedDate  string             `bson:"published_date"`
	Month          int32              `bson:"month"`
	Year           int32              `bson:"year"`
	TranslatedText string             `bson:"translated_text"`
	Service        float64            `bson:"service"`
	Value          float64            `bson:"value"`
	Food           float64            `bson:"food"`
	Vader          float64            `bson:"vader"`
	Wordnet        float64            `bson:"wordnet"`
	CreatedAt      time.Time          `bson:"created_at"`
}
