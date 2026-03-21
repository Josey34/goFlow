package main

import (
	"context"
	"flag"
	"fmt"
	"goflow/cmd"
	"goflow/config"
	"goflow/factory"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.Parse()

	f, err := factory.New(cfg)
	if err != nil {
		log.Fatal("Failed to initialize factory:", err)
	}

	if len(flag.Args()) == 0 {
		fmt.Println("Usage: goflow [health|process|stats]")
		return
	}

	command := flag.Args()[0]

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		fmt.Printf("\nReceived signal: %v", sig)
		fmt.Println("Shutting down gracefully")
		cancel()
	}()

	switch command {
	case "health":
		if err := cmd.Health(ctx, f); err != nil {
			log.Fatal(err)
		}
	case "process":
		if err := cmd.Process(ctx, f); err != nil {
			log.Fatal(err)
		}
	case "stats":
		if err := cmd.Stats(f); err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Println("Unknown command:", command)
	}

	fmt.Println("Shutdown complete")
}
