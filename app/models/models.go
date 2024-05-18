package models

type SubscribeRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type RateResponse struct {
	ExchangeRate float64 `json:"exchange_rate"`
}

type APIError struct {
	Message string `json:"message"`
}

type EmailTemplate struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
}
