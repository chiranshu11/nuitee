package external

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func FetchExchangeRates(targetCurrency string) (map[string]float64, error) {
	resp, err := http.Get(fmt.Sprintf("%s%s", ratesAPI, targetCurrency))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get exchange rates: %s", resp.Status)
	}

	var ratesResponse struct {
		Rates map[string]float64 `json:"rates"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&ratesResponse); err != nil {
		return nil, err
	}

	return ratesResponse.Rates, nil
}
