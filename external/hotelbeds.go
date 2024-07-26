package external

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"liteapi/constants"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	apiKey    = os.Getenv("API_KEY")
	apiSecret = os.Getenv("API_SECRET")
)

type HotelbedsRequest struct {
	CheckIn          string `json:"checkIn" form:"checkin" binding:"required"`
	CheckOut         string `json:"checkOut" form:"checkout" binding:"required"`
	Currency         string `json:"currency" form:"currency" binding:"required"`
	GuestNationality string `json:"guestNationality" form:"guestNationality" binding:"required"`
	HotelIds         string `json:"hotelIds" form:"hotelIds" binding:"required"`
	Occupancies      string `json:"occupancies" form:"occupancies" binding:"required"`
}

type HotelbedsExternalRequest struct {
	Stays struct {
		CheckIn  string `json:"checkIn"`
		CheckOut string `json:"checkOut"`
	} `json:"stay"`
	HotelIds struct {
		Hotel []string `json:"hotel"`
	} `json:"hotels"`
	Occupancies []struct {
		Adults   int `json:"adults"`
		Children int `json:"children"`
		Rooms    int `json:"rooms"`
	} `json:"occupancies"`
}

type HotelbedsResponse struct {
	Hotels struct {
		Hotels []HotelResult `json:"hotels"`
	} `json:"hotels"`
}

type HotelResult struct {
	Code     int    `json:"code"`
	Name     string `json:"name"`
	MinRate  string `json:"minRate"`
	Currency string `json:"currency"`
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

func (h *HotelbedsRequest) ToExternalRequest() (HotelbedsExternalRequest, error) {
	hotelIds := strings.Split(h.HotelIds, ",")

	var occupancies []struct {
		Adults   int `json:"adults"`
		Children int `json:"children"`
		Rooms    int `json:"rooms"`
	}
	if err := json.Unmarshal([]byte(h.Occupancies), &occupancies); err != nil {
		return HotelbedsExternalRequest{}, err
	}

	externalRequest := HotelbedsExternalRequest{
		Stays: struct {
			CheckIn  string `json:"checkIn"`
			CheckOut string `json:"checkOut"`
		}{
			CheckIn:  h.CheckIn,
			CheckOut: h.CheckOut,
		},
		HotelIds: struct {
			Hotel []string `json:"hotel"`
		}{
			Hotel: hotelIds,
		},
		Occupancies: occupancies,
	}

	return externalRequest, nil
}

func FetchHotelbedsRates(targetCurrency string, request HotelbedsExternalRequest) (LiteAPIResponse, error) {
	reqUrl := fmt.Sprintf("%s/hotels", constants.NuiteeApiBaseUrl)

	requestBody, err := json.Marshal(request)
	if err != nil {
		return LiteAPIResponse{}, err
	}

	client := &http.Client{}
	client.Timeout = time.Second * 5
	req, err := http.NewRequest("POST", reqUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		return LiteAPIResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", apiKey)
	req.Header.Set("X-Signature", generateSignature())
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return LiteAPIResponse{}, err
	}

	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return LiteAPIResponse{}, err
	}

	var hotelbedsResponse HotelbedsResponse
	if err := json.NewDecoder(bytes.NewReader(responseBody)).Decode(&hotelbedsResponse); err != nil {
		return LiteAPIResponse{}, err
	}

	exchangeRates, err := FetchExchangeRates(targetCurrency)
	if err != nil {
		return LiteAPIResponse{}, err
	}

	var liteAPIData []LiteAPIData
	for _, hotel := range hotelbedsResponse.Hotels.Hotels {
		price, err := strconv.ParseFloat(hotel.MinRate, 64)
		if err != nil {
			return LiteAPIResponse{}, err
		}

		convertedPrice := price
		if hotel.Currency != targetCurrency {
			conversionRate, exists := exchangeRates[hotel.Currency]
			if exists {
				convertedPrice = price * conversionRate
			}
		}

		liteAPIData = append(liteAPIData, LiteAPIData{
			HotelId:  fmt.Sprintf("%d", hotel.Code),
			Currency: targetCurrency,
			Price:    convertedPrice,
		})
	}

	response := LiteAPIResponse{
		Data: liteAPIData,
		Supplier: SupplierInfo{
			Request:  string(requestBody),
			Response: string(responseBody),
		},
	}

	return response, nil
}

func generateSignature() string {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	signature := sha256.Sum256([]byte(apiKey + apiSecret + timestamp))
	return hex.EncodeToString(signature[:])
}
