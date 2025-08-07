package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fluentlabs-xyz/eth-balance-watcher/config"
	"github.com/fluentlabs-xyz/eth-balance-watcher/ethereum"
	"github.com/fluentlabs-xyz/eth-balance-watcher/metrics"
	"github.com/fluentlabs-xyz/eth-balance-watcher/monitor"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	
	cfg, err := config.Load()
	if err != nil {
		log.WithError(err).Fatal("Failed to load configuration")
	}

	log.WithFields(logrus.Fields{
		"rpc_url":        cfg.EthereumRPC,
		"check_interval": cfg.CheckInterval,
		"metrics_port":   cfg.MetricsPort,
		"wallets_count":  len(cfg.Wallets),
	}).Info("Starting ETH Balance Watcher")

	ethClient, err := ethereum.NewClient(cfg.EthereumRPC)
	if err != nil {
		log.WithError(err).Fatal("Failed to create Ethereum client")
	}
	defer ethClient.Close()

	metricsCollector := metrics.NewCollector()
	
	balanceMonitor := monitor.New(
		ethClient,
		metricsCollector,
		cfg.Wallets,
		cfg.CheckInterval,
		log,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go balanceMonitor.Start(ctx)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/health", healthHandler)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.MetricsPort),
		Handler: mux,
	}

	go func() {
		log.Infof("Starting metrics server on port %d", cfg.MetricsPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("Failed to start metrics server")
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Info("Shutting down gracefully...")
	
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.WithError(err).Error("Failed to shutdown server gracefully")
	}

	cancel()
	
	time.Sleep(1 * time.Second)
	log.Info("Shutdown complete")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy"}`))
}