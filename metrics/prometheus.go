package metrics

import (
	"math/big"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Collector struct {
	balanceWei   *prometheus.GaugeVec
	balanceEther *prometheus.GaugeVec
	lastCheck    *prometheus.GaugeVec
	checkErrors  *prometheus.CounterVec
	checkDuration *prometheus.HistogramVec
}

func NewCollector() *Collector {
	return &Collector{
		balanceWei: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "eth_wallet_balance_wei",
				Help: "Current wallet balance in Wei",
			},
			[]string{"name", "address"},
		),
		balanceEther: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "eth_wallet_balance_ether",
				Help: "Current wallet balance in Ether",
			},
			[]string{"name", "address"},
		),
		lastCheck: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "eth_wallet_last_check_timestamp",
				Help: "Unix timestamp of the last successful balance check",
			},
			[]string{"name", "address"},
		),
		checkErrors: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "eth_wallet_check_errors_total",
				Help: "Total number of balance check errors",
			},
			[]string{"name", "address"},
		),
		checkDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "eth_wallet_check_duration_seconds",
				Help:    "Duration of balance check in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"name", "address"},
		),
	}
}

func (c *Collector) UpdateBalance(name, address string, balanceWei *big.Int, balanceEther *big.Float, checkDuration float64) {
	labels := prometheus.Labels{
		"name":    name,
		"address": address,
	}

	balanceWeiFloat, _ := new(big.Float).SetInt(balanceWei).Float64()
	balanceEtherFloat, _ := balanceEther.Float64()

	c.balanceWei.With(labels).Set(balanceWeiFloat)
	c.balanceEther.With(labels).Set(balanceEtherFloat)
	c.lastCheck.With(labels).SetToCurrentTime()
	c.checkDuration.With(labels).Observe(checkDuration)
}

func (c *Collector) RecordError(name, address string) {
	labels := prometheus.Labels{
		"name":    name,
		"address": address,
	}
	c.checkErrors.With(labels).Inc()
}

func (c *Collector) RecordCheckDuration(name, address string, duration float64) {
	labels := prometheus.Labels{
		"name":    name,
		"address": address,
	}
	c.checkDuration.With(labels).Observe(duration)
}