package main

import (
	"currency_mail/app/router"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	port := os.Getenv("APP_PORT")
	if port == "" {
		log.Fatalf("APP_PORT environment variable not set")
	}
	port_prefixed := ":" + port

	r := router.SetupRouter()
	r.Run(port_prefixed)
}
