package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type GasTracker struct {
	SafeGasPrice    string `json:"SafeGasPrice"`
	ProposeGasPrice string `json:"ProposeGasPrice"`
	FastGasPrice    string `json:"FastGasPrice"`
}

var EtherscanAPIURL = "https://api.etherscan.io/api?module=gastracker&action=gasoracle&apikey="

var (
	gasPrice = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "eth_gas_price",
			Help: "Current ethereum gas prices.",
		},
		[]string{"type"},
	)
)

func fetchGasPrices() {
	apiKey := os.Getenv("ETHERSCAN_API_KEY")
	resp, err := http.Get(EtherscanAPIURL + apiKey)
	if err != nil {
		log.Println("Error fetching gas prices:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading gas prices response:", err)
		return
	}

	var result struct {
		Status  string     `json:"status"`
		Message string     `json:"message"`
		Result  GasTracker `json:"result"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println("Error parsing gas prices JSON:", err)
		return
	}

	safeGasPrice, err := strconv.ParseFloat(result.Result.SafeGasPrice, 64)
	if err == nil {
		gasPrice.WithLabelValues("safe").Set(safeGasPrice)
	}

	proposeGasPrice, err := strconv.ParseFloat(result.Result.ProposeGasPrice, 64)
	if err == nil {
		gasPrice.WithLabelValues("propose").Set(proposeGasPrice)
	}

	fastGasPrice, err := strconv.ParseFloat(result.Result.FastGasPrice, 64)
	if err == nil {
		gasPrice.WithLabelValues("fast").Set(fastGasPrice)
	}
}

func recordMetrics() {
	go func() {
		for {
			fetchGasPrices()
			time.Sleep(10 * time.Second)
		}
	}()
}

func main() {
	// Create a new registry
	r := prometheus.NewRegistry()

	// Register the custom collector with the registry
	r.MustRegister(gasPrice)

	// Start the metrics recording goroutine
	recordMetrics()

	// Start the HTTP server with the custom registry
	http.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{}))
	http.ListenAndServe(":9100", nil)
}
