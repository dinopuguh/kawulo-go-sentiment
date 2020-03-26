package main

import (
	"log"

	"github.com/dinopuguh/kawulo-go-sentiment/database"
	"github.com/dinopuguh/kawulo-go-sentiment/services"
)

func main() {
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err.Error())
	}

	locs := services.FindIndonesianLocations(db)

	for _, loc := range locs {
		services.InsertSentiments(db, loc)
	}

	// err = db.Client().Disconnect(database.Ctx)
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
}
