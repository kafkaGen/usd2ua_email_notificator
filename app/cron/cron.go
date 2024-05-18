package main

import (
	"context"
	"currency_mail/app/db"
	"currency_mail/app/utils"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/robfig/cron"
)

func SendDailyEmails(path_to_template string) {
	conn := db.GetDB()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := conn.Query(ctx, "SELECT email FROM Subscribers")
	if err != nil {
		log.Printf("Error querying subscribers: %v", err)
		return
	}
	defer rows.Close()

	var emails []string
	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			log.Printf("Error scanning email: %v", err)
			continue
		}
		emails = append(emails, email)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error reading rows: %v", err)
		return
	}

	exchangeRate, err := utils.Rate()
	if err != nil {
		log.Fatalf("Error getting exchange rate: %v", err)
	}

	template, err := utils.LoadEmailTemplate(path_to_template)
	if err != nil {
		log.Fatalf("Error loading email template: %v", err)
	}

	subject, body := utils.FormatEmail(template, exchangeRate.ExchangeRate)

	for _, email := range emails {
		if err := utils.SendEmail(email, subject, body); err != nil {
			log.Printf("Error sending email to %s: %v", email, err)
		}
	}

	log.Print("Daily emails sent")
}

func main() {
	godotenv.Load()

	c := cron.New()

	// TODO: pass cron time as env variable also
	c.AddFunc("@daily", func() {
		log.Println("Cron job triggered - calling SendDailyEmails")
		// TODO: change to Strategy pattern, not send path to template directly
		SendDailyEmails("messages/usd2ua.json")
		log.Println("Finished calling SendDailyEmails")
	})

	c.Start()
	log.Println("Cron job started.")

	// Keep the program running to allow the cron jobs to execute.
	select {}
}
