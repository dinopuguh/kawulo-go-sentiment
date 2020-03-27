package services

import (
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/aaaton/golem"
	"github.com/aaaton/golem/dicts/en"
	"github.com/bbalet/stopwords"
	"github.com/dinopuguh/kawulo-go-sentiment/database"
	"github.com/dinopuguh/kawulo-go-sentiment/models"
	"github.com/jonreiter/govader"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/jdkato/prose.v2"
)

func CheckSentimentExist(db *mongo.Database, revId string) bool {
	ctx := database.Ctx

	var result models.Sentiment
	err := db.Collection("sentiment").FindOne(ctx, bson.M{"review_id": revId}).Decode(&result)
	if err != nil {
		return false
	}

	return true
}

func VaderAnalyze(text string) float64 {
	analyzer := govader.NewSentimentIntensityAnalyzer()
	sentiment := analyzer.PolarityScores(text)

	return sentiment.Positive - sentiment.Negative
}

func WordnetAnalyze(text string) {
	// var sentiment float64
	// var count int32

	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err.Error())
	}

	newText := reg.ReplaceAllString(text, " ")

	cleanText := stopwords.CleanString(newText, "en", false)

	doc, err := prose.NewDocument(cleanText)
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, tok := range doc.Tokens() {
		lemmatizer, err := golem.New(en.New())
		if err != nil {
			log.Fatal(err.Error())
		}

		lemmatized := lemmatizer.Lemma(tok.Text)
		log.Println(tok.Text, lemmatized, tok.Tag)
	}

}

func InsertSentiments(db *mongo.Database, loc models.Location) {
	rests := FindRestaurantByLocation(db, loc.LocationId)

	for _, rest := range rests {
		log.Println("---------------------", rest.Name)
		InsertSentiment(db, loc, rest)
	}
}

func InsertSentiment(db *mongo.Database, loc models.Location, rest models.Restaurant) {
	revs := FindReviewByRestaurant(db, rest.ID)

	for _, rev := range revs {
		log.Println("------------------", rev.Id)
		sentExist := CheckSentimentExist(db, rev.Id)

		if sentExist {
			log.Println("Sentiment with review id", rev.Id, "is already exist")
			continue
		}

		text := rev.Text
		lang := rev.Lang

		translatedText := text
		if lang != "en" {
			translatedText = TranslateReview(text, lang, "trnsl.1.1.20191231T003150Z.5e39cb9fd8acfcf0.319e28c6c047015447eaa9f45951fd16a87f9a8c")
		}

		vaderScore := VaderAnalyze(translatedText)
		WordnetAnalyze(translatedText)

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
