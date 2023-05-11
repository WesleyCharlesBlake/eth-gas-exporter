package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestGetETHGasPrices(t *testing.T) {
	// Create a new HTTP client.
	client := &http.Client{}

	// Make a request to the Ethereum blockchain.
	resp, err := client.Get("https://ethgasstation.info/api/ethgasAPI.json")
	if err != nil {
		t.Fatal(err)
	}

	// Read the response body.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	// Unmarshal the JSON response into a map of gas prices.
	gasPrices := make(map[string]interface{})
	err = json.Unmarshal(body, &gasPrices)
	if err != nil {
		t.Fatal(err)
	}

	// Get the gas prices from the map.
	fastPrice, ok := gasPrices["fast"].(float64)
	if !ok {
		t.Errorf("Cannot convert fastPrice to float64")
	}
	fastestPrice, ok := gasPrices["fastest"].(float64)
	if !ok {
		t.Errorf("Cannot convert fastestPrice to float64")
	}
	safeLowPrice, ok := gasPrices["safeLow"].(float64)
	if !ok {
		t.Errorf("Cannot convert safeLowPrice to float64")
	}
	averagePrice, ok := gasPrices["average"].(float64)
	if !ok {
		t.Errorf("Cannot convert averagePrice to float64")
	}

	// Check that the gas prices are valid.
	if fastPrice <= 0 {
		t.Errorf("fastPrice must be greater than 0")
	}
	if fastestPrice <= 0 {
		t.Errorf("fastestPrice must be greater than 0")
	}
	if safeLowPrice <= 0 {
		t.Errorf("safeLowPrice must be greater than 0")
	}
	if averagePrice <= 0 {
		t.Errorf("averagePrice must be greater than 0")
	}

	// Get the gas prices from the GetETHGasPrices function.
	actualGasPrices, err := GetETHGasPrices()
	if err != nil {
		t.Fatal(err)
	}

	// Check that the gas prices from the two functions are equal.
	for key, value := range actualGasPrices {
		if value != gasPrices[key] {
			t.Errorf("Expected gas price for %s to be %f, got %f", key, gasPrices[key], value)
		}
	}
}
