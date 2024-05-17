package handlers

import (
	"currency_mail/app/models"
	"currency_mail/app/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RateHandler(c *gin.Context) {
	response, err := utils.Rate()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
