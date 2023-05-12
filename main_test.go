package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestFetchGasPrices(t *testing.T) {
	// Mock Etherscan API response
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`{"status":"1","message":"OK","result":{"SafeGasPrice":"44","ProposeGasPrice":"44","FastGasPrice":"46"}}`))
	}))
	defer server.Close()

	// Temporarily replace EtherscanAPIURL with our mock server's URL
	oldURL := EtherscanAPIURL
	EtherscanAPIURL = server.URL
	defer func() { EtherscanAPIURL = oldURL }() // restore original URL after test

	// Call fetchGasPrices
	fetchGasPrices()

	// Check that the metrics were correctly updated
	assert.Equal(t, 44.0, testutil.ToFloat64(gasPrice.WithLabelValues("safe")))
	assert.Equal(t, 44.0, testutil.ToFloat64(gasPrice.WithLabelValues("propose")))
	assert.Equal(t, 46.0, testutil.ToFloat64(gasPrice.WithLabelValues("fast")))
}
