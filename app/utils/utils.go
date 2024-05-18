package utils

import (
	"currency_mail/app/models"
	"io/ioutil"
	"strconv"
	"strings"

	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/go-gomail/gomail"
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

func SendEmail(email, subject string, body string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPortStr := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")

	if smtpHost == "" || smtpPortStr == "" || smtpUser == "" || smtpPass == "" {
		return fmt.Errorf("SMTP configuration not set")
	}

	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		return fmt.Errorf("invalid SMTP port: %v", err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", smtpUser)
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)

	return d.DialAndSend(m)
}

func LoadEmailTemplate(filePath string) (models.EmailTemplate, error) {
	var template models.EmailTemplate

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return template, fmt.Errorf("failed to read template file: %w", err)
	}

	err = json.Unmarshal(data, &template)
	if err != nil {
		return template, fmt.Errorf("failed to unmarshal template: %w", err)
	}

	return template, nil
}

func FormatEmail(template models.EmailTemplate, exchangeRate float64) (string, string) {
	exchangeRateStr := fmt.Sprintf("%.2f", exchangeRate)
	body := strings.Replace(template.Body, "[Current Exchange Rate]", exchangeRateStr, -1)
	return template.Subject, body
}
