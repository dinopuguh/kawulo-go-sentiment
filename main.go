package main

import (
	"log"
	"strconv"
	"time"

	"github.com/dinopuguh/kawulo-go-sentiment/database"
	"github.com/dinopuguh/kawulo-go-sentiment/models"
	"github.com/dinopuguh/kawulo-go-sentiment/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
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
					translatedText = services.TranslateReview(text, lang, "trnsl.1.1.20191231T003150Z.5e39cb9fd8acfcf0.319e28c6c047015447eaa9f45951fd16a87f9a8c")
				}

				vaderScore := services.VaderAnalyze(translatedText)
				services.WordnetAnalyze(translatedText)

				service, _ := strconv.ParseFloat(rev.Rating, 64)
				value, _ := strconv.ParseFloat(rev.Rating, 64)
				food, _ := strconv.ParseFloat(rev.Rating, 64)

				publishedDate, err := time.Parse("2006-01-02T15:04:05-04:00", rev.PublishedDate)
				if err != nil {
					log.Fatal(err.Error())
				}

				var result models.Sentiment

				result.ID = primitive.NewObjectID()
				result.PublishedDate = rev.PublishedDate
				result.LocationId = rest.LocationID
				result.Location = loc
				result.RestaurantId = rest.LocationId
				result.Restaurant = rest
				result.ReviewId = rev.Id
				result.Review = rev
				result.Month = int32(publishedDate.Month())
				result.Year = int32(publishedDate.Year())
				result.TranslatedText = translatedText
				result.Service = service
				result.Value = value
				result.Food = food
				result.Vader = vaderScore
				result.Wordnet = 0
				result.CreatedAt = time.Now()

				log.Println(result)
			}
		}
	}
}
