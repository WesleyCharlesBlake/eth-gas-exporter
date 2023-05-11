### Build/test the Go code

To build the Go code, run the following command:

```
go build
```

To test the Go code, run the following command:

```
go test ./...
```

### Run the Docker image

To run the Docker image, run the following command:

```
docker build -t gwei_exporter .
docker run -p 9090:9090 gwei_exporter
```

The exporter will listen on port `9090` and the path `/metrics` by default.

```
curl http://localhost:9090/metrics
# HELP eth_gas_price The current gas price of ETH in gwei
# TYPE eth_gas_price gauge
eth_gas_price{type="average"} 143
eth_gas_price{type="fast"} 530
eth_gas_price{type="fastest"} 530
eth_gas_price{type="safeLow"} 143
```

### Include scrape configs

To include scrape configs in Prometheus, add the following to your Prometheus configuration file:

```
scrape_configs:
  - job_name: "gwei"
    scrape_interval: 10s
    static_configs:
      - targets: ["localhost:9090"]
```

This will tell Prometheus to scrape the exporter every 10 seconds.