package main

import (
	"currency_mail/app/router"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	r := router.SetupRouter()
	r.Run(":8080")
}
