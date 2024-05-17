package utils

import (
	"currency_mail/app/models"

	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"log"
	"regexp"
	"time"
	"github.com/jackc/pgx/v4/pgxpool"
)

func IsValidEmail(email string) bool {
	const emailRegex = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

func Rate() (models.RateResponse, error) {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		return models.RateResponse{}, fmt.Errorf("API_KEY environment variable not set")
	}
	url := os.Getenv("URL")
	if url == "" {
		return models.RateResponse{}, fmt.Errorf("URL environment variable not set")
	}

	url = fmt.Sprintf(url, apiKey)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return models.RateResponse{}, fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return models.RateResponse{}, fmt.Errorf("failed to get data from API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.RateResponse{}, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return models.RateResponse{}, fmt.Errorf("failed to decode JSON response: %v", err)
	}

	rates, ok := result["conversion_rates"].(map[string]interface{})
	if !ok {
		return models.RateResponse{}, fmt.Errorf("conversion_rates not found or invalid")
	}

	uaRate, ok := rates["UAH"].(float64)
	if !ok {
		return models.RateResponse{}, fmt.Errorf("UAH rate not found or invalid")
	}

	return models.RateResponse{ExchangeRate: uaRate}, nil
}

func GetDB() *pgxpool.Pool {
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	if dbName == "" || dbUser == "" || dbPassword == "" || dbHost == "" || dbPort == "" {
		log.Fatalf("Database environment variables not set")
	}

	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)

	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		log.Fatalf("Unable to parse database URL: %v\n", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbpool, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	log.Println("Database connection established")

	return dbpool
}