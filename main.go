package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/dinopuguh/kawulo-go-sentiment/database"
	"github.com/dinopuguh/kawulo-go-sentiment/models"
	"github.com/dinopuguh/kawulo-go-sentiment/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	yandexAPIKeyId := 1

	db, err := database.Connect()
	if err != nil {
		log.Fatal(err.Error())
	}

	locs := services.FindIndonesianLocations(db)

	for _, loc := range locs {
		log.Println("------------------------", loc.Name)
		rests := services.FindRestaurantByLocation(db, loc.LocationId)

		for _, rest := range rests {
			log.Println("---------------------", rest.Name)
			revs := services.FindReviewByRestaurant(db, rest.ID)

			for _, rev := range revs {
				log.Println("------------------", rev.Id)
				sentExist := services.CheckSentimentExist(db, rev.Id)

				if sentExist {
					log.Println("Sentiment with review id", rev.Id, "is already exist")
					continue
				}

				text := rev.Text
				lang := rev.Lang

				translatedText := text
				if lang != "en" {
					translatedText, err = services.TranslateReview(text, lang, os.Getenv("YANDEX_API_KEY"+strconv.Itoa(yandexAPIKeyId)))
					if err != nil && yandexAPIKeyId < 4 {
						yandexAPIKeyId++
						log.Println("Change Yandex API Key", yandexAPIKeyId)
						translatedText, _ = services.TranslateReview(text, lang, os.Getenv("YANDEX_API_KEY"+strconv.Itoa(yandexAPIKeyId)))
					}
				}

				vaderScore := services.VaderAnalyze(translatedText)
				wordnetScore := services.WordnetAnalyze(translatedText)

				service, _ := strconv.ParseFloat(rev.Rating, 64)
				value, _ := strconv.ParseFloat(rev.Rating, 64)
				food, _ := strconv.ParseFloat(rev.Rating, 64)

				publishedDate, err := time.Parse("2006-01-02T15:04:05-04:00", rev.PublishedDate)
				if err != nil {
					log.Fatal(err.Error())
				}

				result := models.Sentiment{
					ID:             primitive.NewObjectID(),
					PublishedDate:  rev.PublishedDate,
					LocationId:     rest.LocationID,
					Location:       loc,
					RestaurantId:   rest.LocationId,
					Restaurant:     rest,
					ReviewId:       rev.Id,
					Review:         rev,
					Month:          int32(publishedDate.Month()),
					Year:           int32(publishedDate.Year()),
					TranslatedText: translatedText,
					Service:        service,
					Value:          value,
					Food:           food,
					Vader:          vaderScore,
					Wordnet:        wordnetScore,
					CreatedAt:      time.Now(),
				}

				services.InsertSentiment(db, result)
			}
		}
	}
}
