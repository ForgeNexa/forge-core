package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ForgeNexa/forge-core/internal/concurrency"
)

const Version = "0.1.0"

func main() {
	fmt.Printf("[forge-core] ForgeNexa Systems Core v%s initializing...\n", Version)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Initialize background worker pool
	pool := concurrency.NewPool(ctx, 8, 32)
	defer pool.Stop()

	fmt.Println("[forge-core] Systems Core running. Ready for IPC/Service coordination.")

	<-sigChan
	fmt.Println("[forge-core] Shutdown signal received. Performing graceful resource cleanup...")
}
