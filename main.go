package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// GetETHGasPrices gets the current ETH gas prices in gwei.
func GetETHGasPrices() (map[string]float64, error) {
	// Get the prices from the Ethereum blockchain.
	resp, err := http.Get("https://ethgasstation.info/api/ethgasAPI.json")
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var gasPrices map[string]interface{}
	err = json.Unmarshal(body, &gasPrices)
	if err != nil {
		return nil, err
	}

	// Get the gas prices from the JSON object.
	fastPrice, ok := gasPrices["fast"].(float64)
	if !ok {
		return nil, fmt.Errorf("Cannot convert fastPrice to float64")
	}
	fastestPrice, ok := gasPrices["fastest"].(float64)
	if !ok {
		return nil, fmt.Errorf("Cannot convert fastestPrice to float64")
	}
	safeLowPrice, ok := gasPrices["safeLow"].(float64)
	if !ok {
		return nil, fmt.Errorf("Cannot convert safeLowPrice to float64")
	}
	averagePrice, ok := gasPrices["average"].(float64)
	if !ok {
		return nil, fmt.Errorf("Cannot convert averagePrice to float64")
	}

	// Return the gas prices.
	return map[string]float64{
		"fast":    fastPrice,
		"fastest": fastestPrice,
		"safeLow": safeLowPrice,
		"average": averagePrice,
	}, nil
}

// main is the entry point for the program.
func main() {
	// Create a new Prometheus registry.
	reg := prometheus.NewRegistry()

	// Create a new Gauge metric for the ETH gas price.
	ethGasPriceGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "eth_gas_price",
		Help: "The current gas price of ETH in gwei",
	}, []string{"type"})

	// Register the ETH gas price gauge with the registry.
	reg.MustRegister(ethGasPriceGauge)

	// Get the current ETH gas prices.
	prices, err := GetETHGasPrices()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Set the value of the ETH gas price gauge.
	ethGasPriceGauge.WithLabelValues("fast").Set(prices["fast"])
	ethGasPriceGauge.WithLabelValues("fastest").Set(prices["fastest"])
	ethGasPriceGauge.WithLabelValues("safeLow").Set(prices["safeLow"])
	ethGasPriceGauge.WithLabelValues("average").Set(prices["average"])

	// Start a Prometheus HTTP server on port 9090.
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	http.ListenAndServe(":9090", nil)
}
