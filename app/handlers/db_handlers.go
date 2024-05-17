package handlers

import (
	"context"
	"currency_mail/app/models"
	"currency_mail/app/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SubscribeHandler(c *gin.Context) {
	var req models.SubscribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIError{Message: err.Error()})
		return
	}

	if !utils.IsValidEmail(req.Email) {
		c.JSON(http.StatusBadRequest, models.APIError{Message: "Invalid email address"})
		return
	}

	conn := utils.GetDB()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var emailCount int
	err := conn.QueryRow(ctx, "SELECT COUNT(*) FROM Subscribers WHERE email=$1", req.Email).Scan(&emailCount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{Message: "Failed to check existing subscriptions"})
		return
	}

	if emailCount > 0 {
		c.JSON(http.StatusConflict, models.APIError{Message: "Email is already subscribed"})
		return
	}

	_, err = conn.Exec(ctx, "INSERT INTO Subscribers (email) VALUES ($1)", req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{Message: "Failed to subscribe"})
		return
	}

	conn.Close()

	c.JSON(http.StatusOK, gin.H{
		"message": "Subscription successful",
		"email":   req.Email,
	})
}

func UnsubscribeHandler(c *gin.Context) {
	var req models.SubscribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIError{Message: err.Error()})
		return
	}

	if !utils.IsValidEmail(req.Email) {
		c.JSON(http.StatusBadRequest, models.APIError{Message: "Invalid email address"})
		return
	}

	conn := utils.GetDB()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var emailCount int
	err := conn.QueryRow(ctx, "SELECT COUNT(*) FROM Subscribers WHERE email=$1", req.Email).Scan(&emailCount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{Message: "Failed to check existing subscriptions"})
		return
	}

	if emailCount == 0 {
		c.JSON(http.StatusNotFound, models.APIError{Message: "Email not found"})
		return
	}

	_, err = conn.Exec(ctx, "DELETE FROM Subscribers WHERE email=$1", req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{Message: "Failed to unsubscribe"})
		return
	}

	conn.Close()

	c.JSON(http.StatusOK, gin.H{
		"message": "Unsubscription successful",
		"email":   req.Email,
	})
}
