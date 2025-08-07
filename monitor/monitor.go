package monitor

import (
	"context"
	"sync"
	"time"

	"github.com/fluentlabs-xyz/eth-balance-watcher/config"
	"github.com/fluentlabs-xyz/eth-balance-watcher/ethereum"
	"github.com/fluentlabs-xyz/eth-balance-watcher/metrics"
	"github.com/sirupsen/logrus"
)

type Monitor struct {
	ethClient *ethereum.Client
	metrics   *metrics.Collector
	wallets   []config.Wallet
	interval  time.Duration
	log       *logrus.Logger
}

func New(
	ethClient *ethereum.Client,
	metricsCollector *metrics.Collector,
	wallets []config.Wallet,
	interval time.Duration,
	log *logrus.Logger,
) *Monitor {
	return &Monitor{
		ethClient: ethClient,
		metrics:   metricsCollector,
		wallets:   wallets,
		interval:  interval,
		log:       log,
	}
}

func (m *Monitor) Start(ctx context.Context) {
	m.log.Info("Starting balance monitor")
	
	m.checkAllBalances()

	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			m.log.Info("Balance monitor stopped")
			return
		case <-ticker.C:
			m.checkAllBalances()
		}
	}
}

func (m *Monitor) checkAllBalances() {
	m.log.WithField("wallets_count", len(m.wallets)).Debug("Checking all wallet balances")
	
	var wg sync.WaitGroup
	
	for _, wallet := range m.wallets {
		wg.Add(1)
		go func(w config.Wallet) {
			defer wg.Done()
			m.checkWalletBalance(w)
		}(wallet)
	}
	
	wg.Wait()
}

func (m *Monitor) checkWalletBalance(wallet config.Wallet) {
	start := time.Now()
	
	walletLog := m.log.WithFields(logrus.Fields{
		"wallet_name": wallet.Name,
		"address":     wallet.Address,
	})

	balanceWei, err := m.ethClient.GetBalance(wallet.Address)
	if err != nil {
		walletLog.WithError(err).Error("Failed to get wallet balance")
		m.metrics.RecordError(wallet.Name, wallet.Address)
		return
	}

	balanceEther := ethereum.WeiToEther(balanceWei)
	duration := time.Since(start).Seconds()

	m.metrics.UpdateBalance(wallet.Name, wallet.Address, balanceWei, balanceEther, duration)

	balanceEtherFloat, _ := balanceEther.Float64()
	
	walletLog.WithFields(logrus.Fields{
		"balance_wei":   balanceWei.String(),
		"balance_ether": balanceEtherFloat,
		"duration_ms":   duration * 1000,
	}).Debug("Balance updated")
}

func (m *Monitor) CheckOnce() {
	m.checkAllBalances()
}