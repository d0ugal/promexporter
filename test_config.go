package main

import (
	"fmt"
	"log"

	"github.com/d0ugal/promexporter/config"
)

func main() {
	// Test with explicit values
	cfg1, err := config.Load("test-config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Test 1 - Explicit values:\n")
	fmt.Printf("  Web UI Enabled: %v\n", cfg1.Server.IsWebUIEnabled())
	fmt.Printf("  Health Enabled: %v\n", cfg1.Server.IsHealthEnabled())

	// Test with defaults
	cfg2, err := config.Load("test-config-defaults.yaml")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nTest 2 - Default values:\n")
	fmt.Printf("  Web UI Enabled: %v\n", cfg2.Server.IsWebUIEnabled())
	fmt.Printf("  Health Enabled: %v\n", cfg2.Server.IsHealthEnabled())
}
