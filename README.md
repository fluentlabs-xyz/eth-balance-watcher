# ETH Balance Watcher

Simple Golang service that monitors Ethereum wallet balances and exports metrics to Prometheus.

## Features

- Simple text-based wallet configuration (`name:address` format)
- Configurable balance check intervals
- Prometheus metrics with wallet name and address labels
- Docker containerization
- Health check endpoint

## Quick Start

### 1. Configure Wallets

Create `wallets.txt`:
```
alice:0xd8da6bf26964af9d7eed9e03e53415d37aa96045
bob:0x1f9840a85d5af5bf1d1762f925bdaddc4201f984
```

### 2. Run with Docker Compose

```bash
# Start the service
docker-compose up -d

# View logs
docker-compose logs -f eth-balance-watcher
```

Endpoints:
- **Metrics**: http://localhost:9090/metrics
- **Health**: http://localhost:9090/health
- **Prometheus**: http://localhost:9091

### 3. Run Locally

```bash
go mod download
go run .
```

## Configuration

Environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `ETH_RPC_URL` | Required | Ethereum RPC endpoint |
| `CHECK_INTERVAL` | `60s` | Balance check interval |
| `METRICS_PORT` | `9090` | Metrics server port |
| `WALLETS_FILE` | `wallets.txt` | Wallets config file |

## Prometheus Metrics

All metrics include labels: `{name="wallet_name", address="0x..."}`

- `eth_wallet_balance_wei{name, address}` - Balance in Wei
- `eth_wallet_balance_ether{name, address}` - Balance in Ether  
- `eth_wallet_last_check_timestamp{name, address}` - Last check time
- `eth_wallet_check_errors_total{name, address}` - Error count
- `eth_wallet_check_duration_seconds{name, address}` - Check duration

## Docker

```bash
# Build
docker build -t eth-balance-watcher .

# Run
docker run -p 9090:9090 \
  -e ETH_RPC_URL="your-rpc-url" \
  -v $(pwd)/wallets.txt:/app/wallets.txt:ro \
  eth-balance-watcher
```

## Development

```bash
# Install dependencies
make deps

# Build
make build

# Run tests
make test

# Format code
make fmt
```