package ethereum

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Client struct {
	client  *ethclient.Client
	timeout time.Duration
}

func NewClient(rpcURL string) (*Client, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum node: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if _, err := client.ChainID(ctx); err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	return &Client{
		client:  client,
		timeout: 10 * time.Second,
	}, nil
}

func (c *Client) GetBalance(address string) (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	addr := common.HexToAddress(address)
	balance, err := c.client.BalanceAt(ctx, addr, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance for %s: %w", address, err)
	}

	return balance, nil
}

func (c *Client) GetBalanceInEther(address string) (*big.Float, error) {
	balanceWei, err := c.GetBalance(address)
	if err != nil {
		return nil, err
	}

	ether := new(big.Float).SetInt(balanceWei)
	ether.Quo(ether, big.NewFloat(1e18))

	return ether, nil
}

func (c *Client) Close() {
	c.client.Close()
}

func WeiToEther(wei *big.Int) *big.Float {
	ether := new(big.Float).SetInt(wei)
	return ether.Quo(ether, big.NewFloat(1e18))
}