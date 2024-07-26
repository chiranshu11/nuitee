package main

import (
	"encoding/json"
	"fmt"
	"liteapi/constants"
	"liteapi/utilities"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	apiKey    = "db11033c50b5ed53ab7b815cb1b2eaee"
	apiSecret = "704773c0e3"
)

type HotelbedsRequest struct {
	CheckIn          string      `form:"checkin" binding:"required"`
	CheckOut         string      `form:"checkout" binding:"required"`
	Currency         string      `form:"currency" binding:"required"`
	GuestNationality string      `form:"guestNationality" binding:"required"`
	HotelIds         string      `form:"hotelIds" binding:"required"`
	Occupancies      []Occupancy `json:"occupancies"`
}

type Occupancy struct {
	Rooms  int `json:"rooms"`
	Adults int `json:"adults"`
}

type Geolocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Radius    int     `json:"radius"`
}

type HotelbedsResponse struct {
	HotelResults struct {
		HotelResults []HotelResult `json:"hotelResults"`
	} `json:"hotelResults"`
}

type HotelResult struct {
	HotelInfo struct {
		HotelCode  string  `json:"hotelCode"`
		HotelName  string  `json:"hotelName"`
		Address    string  `json:"hotelAddress"`
		Picture    string  `json:"hotelPictureUrl"`
		Longitude  float64 `json:"longitude,string"`
		Latitude   float64 `json:"latitude,string"`
		StarRating int     `json:"starRating"`
	} `json:"hotelInfo"`
	MinPrice    float64 `json:"minPrice"`
	RateDetails struct {
		RateDetails []RateDetail `json:"rateDetails"`
	} `json:"rateDetails"`
}

type RateDetail struct {
	RateDetailCode string  `json:"rateDetailCode"`
	TotalPrice     float64 `json:"totalPrice"`
	Tax            float64 `json:"tax"`
	HotelFees      float64 `json:"hotelFees"`
}

type LiteAPIResponse struct {
	Data     []LiteAPIData `json:"data"`
	Supplier SupplierInfo  `json:"supplier"`
}

type LiteAPIData struct {
	HotelId  string  `json:"hotelId"`
	Currency string  `json:"currency"`
	Price    float64 `json:"price"`
}

type SupplierInfo struct {
	Request  string `json:"request"`
	Response string `json:"response"`
}

func main() {
	r := gin.Default()
	r.POST("/getRates", getRates)
	r.Run(":9001")
}

func getRates(c *gin.Context) {
	var hotelbedsRequest HotelbedsRequest

	if err := c.ShouldBindQuery(&hotelbedsRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hotelbedsResponse, err := fetchHotelbedsRates(hotelbedsRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var liteAPIData []LiteAPIData
	for _, hotel := range hotelbedsResponse.HotelResults.HotelResults {
		for _, rateDetail := range hotel.RateDetails.RateDetails {
			liteAPIData = append(liteAPIData, LiteAPIData{
				HotelId:  hotel.HotelInfo.HotelCode,
				Currency: hotelbedsRequest.Currency,
				Price:    rateDetail.TotalPrice,
			})
		}
	}

	req, _ := json.Marshal(hotelbedsRequest)
	reqString := string(req)

	resp, _ := json.Marshal(hotelbedsResponse)
	respString := string(resp)

	response := LiteAPIResponse{
		Data: liteAPIData,
		Supplier: SupplierInfo{
			Request:  reqString,
			Response: respString,
		},
	}

	c.JSON(http.StatusOK, response)
}

func fetchHotelbedsRates(request HotelbedsRequest) (HotelbedsResponse, error) {
	reqUrl := fmt.Sprintf("%s/hotels?checkin=%s&checkout=%s&currency=%s&guestNationality=%s&hotelIds=%s&occupancies=%s",
		constants.APIBASEURL, request.CheckIn, request.CheckOut, request.Currency, request.GuestNationality, request.HotelIds, request.Occupancies)

	client := &http.Client{}
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return HotelbedsResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", apiKey)
	req.Header.Set("X-Signature", utilities.GenerateSignature(apiKey+apiSecret))
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return HotelbedsResponse{}, err
	}
	defer resp.Body.Close()

	var hotelbedsResponse HotelbedsResponse
	fmt.Println(" body \n", resp.Body)
	err = json.NewDecoder(resp.Body).Decode(&hotelbedsResponse)
	if err != nil {
		return HotelbedsResponse{}, err
	}

	return hotelbedsResponse, nil
}
