package main

import (
	"liteapi/external"
	"liteapi/utils"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func init() {
	// Load environment variables
	utils.LoadEnv()
}
func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET("/hotels", getRates)
	if err := r.Run(":" + os.Getenv("APP_PORT")); err != nil {
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
