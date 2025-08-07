package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	EthereumRPC   string        `yaml:"ethereum_rpc"`
	CheckInterval time.Duration `yaml:"check_interval"`
	MetricsPort   int           `yaml:"metrics_port"`
	WalletsFile   string        `yaml:"wallets_file"`
	Wallets       []Wallet
}

type Wallet struct {
	Name    string
	Address string
}

func Load() (*Config, error) {
	cfg := &Config{
		EthereumRPC:   getEnv("ETH_RPC_URL", "https://eth-mainnet.g.alchemy.com/v2/your-api-key"),
		CheckInterval: getDurationEnv("CHECK_INTERVAL", 60*time.Second),
		MetricsPort:   getIntEnv("METRICS_PORT", 9090),
		WalletsFile:   getEnv("WALLETS_FILE", "wallets.txt"),
	}

	configFile := getEnv("CONFIG_FILE", "config.yaml")
	if _, err := os.Stat(configFile); err == nil {
		data, err := os.ReadFile(configFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}
	}

	wallets, err := loadWallets(cfg.WalletsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load wallets: %w", err)
	}
	cfg.Wallets = wallets

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

func loadWallets(filename string) ([]Wallet, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open wallets file: %w", err)
	}
	defer file.Close()

	var wallets []Wallet
	scanner := bufio.NewScanner(file)
	lineNum := 0
	
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid format on line %d: expected 'name:address', got '%s'", lineNum, line)
		}

		name := strings.TrimSpace(parts[0])
		address := strings.TrimSpace(parts[1])

		if name == "" {
			return nil, fmt.Errorf("empty wallet name on line %d", lineNum)
		}

		if !isValidEthereumAddress(address) {
			return nil, fmt.Errorf("invalid Ethereum address '%s' on line %d", address, lineNum)
		}

		wallets = append(wallets, Wallet{
			Name:    name,
			Address: address,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading wallets file: %w", err)
	}

	if len(wallets) == 0 {
		return nil, fmt.Errorf("no wallets found in file")
	}

	return wallets, nil
}

func (c *Config) Validate() error {
	if c.EthereumRPC == "" {
		return fmt.Errorf("ethereum_rpc is required")
	}

	if c.CheckInterval < 10*time.Second {
		return fmt.Errorf("check_interval must be at least 10 seconds")
	}

	if c.MetricsPort < 1 || c.MetricsPort > 65535 {
		return fmt.Errorf("invalid metrics_port: %d", c.MetricsPort)
	}

	if len(c.Wallets) == 0 {
		return fmt.Errorf("no wallets configured")
	}

	return nil
}

func isValidEthereumAddress(address string) bool {
	if !strings.HasPrefix(address, "0x") {
		return false
	}
	
	if len(address) != 42 {
		return false
	}

	for i := 2; i < len(address); i++ {
		ch := address[i]
		if !((ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')) {
			return false
		}
	}

	return true
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var intValue int
		if _, err := fmt.Sscanf(value, "%d", &intValue); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}