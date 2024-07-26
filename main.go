package main

import (
	"net/http"

	"hotel-api/external"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/hotels", getRates)
	if err := r.Run(":8080"); err != nil {
		panic(err.Error())
	}
}

func getRates(c *gin.Context) {
	var hotelbedsRequest external.HotelbedsRequest

	if err := c.ShouldBindQuery(&hotelbedsRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hotelbedsExternalRequest, err := hotelbedsRequest.ToExternalRequest()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response, err := external.FetchHotelbedsRates(hotelbedsRequest.Currency, hotelbedsExternalRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
