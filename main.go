package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
	"github.com/dinopuguh/kawulo-go-sentiment/database"
	"github.com/dinopuguh/kawulo-go-sentiment/models"
	"github.com/dinopuguh/kawulo-go-sentiment/services"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	brokers        = fmt.Sprintf("%v:%v", os.Getenv("KAFKA_SERVER"), os.Getenv("KAFKA_PORT"))
	group          = os.Getenv("KAWULO_SENTIMENT_CONSUMER")
	topics         = os.Getenv("KAWULO_TRANSLATE_TOPICS")
	version        = os.Getenv("KAFKA_VERSION")
	yandexAPIKeyId = 1
)

type Consumer struct {
	ready chan bool
}

func main() {
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	logrus.SetFormatter(customFormatter)

	logrus.Infof("Brokers: %v, Group: %v, Topics: %v, Version: %v", brokers, group, topics, version)

	version, err := sarama.ParseKafkaVersion(version)
	if err != nil {
		logrus.Panicf("Error parsing Kafka version: %v", err)
	}

	kafkaConfig := getKafkaConfig("", "")
	kafkaConfig.Version = version

	consumer := Consumer{
		ready: make(chan bool),
	}

	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(strings.Split(brokers, ","), group, kafkaConfig)
	if err != nil {
		logrus.Panicf("Error creating consumer group client: %v", err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(5)
	go func() {
		defer wg.Done()
		for {
			if err := client.Consume(ctx, strings.Split(topics, ","), &consumer); err != nil {
				logrus.Panicf("Error from consumer: %v", err)
			}

			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready
	logrus.Infoln("Sarama consumer up and running!...")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		logrus.Infoln("terminating: context cancelled")
	case <-sigterm:
		logrus.Infoln("terminating: via signal")
	}
	cancel()
	wg.Wait()
	if err = client.Close(); err != nil {
		logrus.Panicf("Error closing client: %v", err)
	}
}

func getKafkaConfig(username, password string) *sarama.Config {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	kafkaConfig.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin

	if username != "" {
		kafkaConfig.Net.SASL.Enable = true
		kafkaConfig.Net.SASL.User = username
		kafkaConfig.Net.SASL.Password = password
	}
	return kafkaConfig
}

func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	close(consumer.ready)
	return nil
}

func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	db, err := database.Connect()
	if err != nil {
		logrus.Fatal(err)
	}

	for message := range claim.Messages() {
		translatedMsg := &models.TranslatedMessage{}

		err := json.Unmarshal([]byte(message.Value), translatedMsg)
		if err != nil {
			logrus.Errorf("Unable to unmarshal kafka message: %v", err)
			return err
		}

		saveSentiment(db, session, message, translatedMsg, translatedMsg.Translated)
	}

	return nil
}

func saveSentiment(db *mongo.Database, session sarama.ConsumerGroupSession, message *sarama.ConsumerMessage, translatedMsg *models.TranslatedMessage, translatedText string) {
	vaderScore := services.VaderAnalyze(translatedText)
	wordnetScore := services.WordnetAnalyze(translatedText)

	service, _ := strconv.ParseFloat(translatedMsg.Review.Rating, 64)
	value, _ := strconv.ParseFloat(translatedMsg.Review.Rating, 64)
	food, _ := strconv.ParseFloat(translatedMsg.Review.Rating, 64)

	publishedDate, err := time.Parse("2006-01-02T15:04:05-04:00", translatedMsg.Review.PublishedDate)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	result := models.Sentiment{
		ID:             primitive.NewObjectID(),
		PublishedDate:  translatedMsg.Review.PublishedDate,
		LocationId:     translatedMsg.Restaurant.LocationID,
		Location:       translatedMsg.Location,
		RestaurantId:   translatedMsg.Restaurant.LocationId,
		Restaurant:     translatedMsg.Restaurant,
		ReviewId:       translatedMsg.Review.Id,
		Review:         translatedMsg.Review,
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

	logrus.Infof("(%v/%v) - %v / %v / %v / %v / %v - %v - %v [%v / %v]", result.Month, result.Year, result.Service, result.Value, result.Food, result.Vader, result.Wordnet, result.RestaurantId, result.Review.Lang, message.Partition, message.Offset)

	err = services.InsertSentiment(db, result)
	if err == nil {
		session.MarkMessage(message, "")
	}
}
