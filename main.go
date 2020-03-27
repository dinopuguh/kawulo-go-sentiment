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
		log.Println("------------------------", loc.Name)
		services.InsertSentiments(db, loc)
	}
}
