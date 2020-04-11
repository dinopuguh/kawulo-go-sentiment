package services

import (
	"log"
	"regexp"
	"strings"

	"github.com/aaaton/golem/v4"
	"github.com/aaaton/golem/v4/dicts/en"
	"github.com/bbalet/stopwords"
	"github.com/dinopuguh/gosentiwordnet"
	"github.com/dinopuguh/kawulo-go-sentiment/database"
	"github.com/dinopuguh/kawulo-go-sentiment/models"
	"github.com/jonreiter/govader"
	"go.mongodb.org/mongo-driver/bson"
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

func InsertSentiment(db *mongo.Database, sentiment models.Sentiment) {
	ctx := database.Ctx

	crs, err := db.Collection("sentiment").InsertOne(ctx, sentiment)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Insert sentiment success -", crs.InsertedID)
}

func VaderAnalyze(text string) float64 {
	analyzer := govader.NewSentimentIntensityAnalyzer()
	sentiment := analyzer.PolarityScores(text)

	return sentiment.Positive - sentiment.Negative
}

func getWordnetPosTag(posTag string) string {
	result := "n"

	wordnetPosTag := make(map[string]string)
	wordnetPosTag["J"] = "a"
	wordnetPosTag["N"] = "n"
	wordnetPosTag["V"] = "v"
	wordnetPosTag["R"] = "r"

	if val, ok := wordnetPosTag[string(posTag[0])]; ok {
		result = val
	}

	return result
}

func WordnetAnalyze(text string) float64 {
	sa := gosentiwordnet.NewGoSentiwordnet()

	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err.Error())
	}

	newText := reg.ReplaceAllString(text, " ")

	cleanText := strings.ToLower(stopwords.CleanString(newText, "en", false))

	doc, err := prose.NewDocument(cleanText)
	if err != nil {
		log.Fatal(err.Error())
	}

	var degree float64
	var sentimentTemp float64 = 0
	count := 0
	for _, tok := range doc.Tokens() {
		lemmatizer, err := golem.New(en.New())
		if err != nil {
			log.Fatal(err.Error())
		}

		lemmatized := lemmatizer.Lemma(tok.Text)
		posTag := getWordnetPosTag(tok.Tag)

		exist, sentiment := sa.GetSentimentScore(lemmatized, posTag, "1")
		if exist {
			count++

			if sentiment.Positive == 0 && sentiment.Negative == 0 {
				degree = 0
			} else if sentiment.Positive >= sentiment.Negative {
				degree = sentiment.Positive / (sentiment.Positive + sentiment.Negative)
			} else {
				degree = -1 * sentiment.Negative / (sentiment.Positive + sentiment.Negative)
			}

			sentimentTemp += degree - (sentiment.Objective * degree)
		}
	}

	if count == 0 {
		return 0
	}

	return sentimentTemp / float64(count)
}
