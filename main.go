package main

import (
	"fmt"

	"log"
	 "time"

	"fabric-setup/hflib"
)

// Example usage of hflib
func main() {
	// Create config
	config := &hflib.Config{
		ConfigFile:  "./config.yaml",
		ChaincodeID: "base",
		ChannelID:   "mychannel",
		User:        "Admin",
		Org:         "Org1",
	}

	// Initalize context
	ctx := hflib.Init(config)

	// Log an event
	err := ctx.LogEvent("faucet", "1234", "6d61697961686565")
	if err != nil {
		log.Fatalf("LogEvent failed: %v", err)
	}
	err = ctx.LogEvent("faucet", "40", "6d61697961686565")
	if err != nil {
		log.Fatalf("LogEvent failed: %v", err)
	}
	err = ctx.LogEvent("fauci", "1234", "6d61697961686565")
	if err != nil {
		log.Fatalf("LogEvent failed: %v", err)
	}

	// Wait for event to complete
	time.Sleep(2 * time.Second)

	Query event
	k, err := ctx.QueryEvent("faucet", "40")
	if err != nil {
		log.Fatalf("QueryEvent failed: %v", err)
	}
	fmt.Printf(k)

	c, err := ctx.QueryDeviceByDateRange("41", "2000","fauci")
	if err != nil {
		log.Fatalf("QueryDeviceByDateRange failed: %v", err)
	}
	fmt.Printf(c)

	f, err := ctx.QueryAllByDateRange("41", "2000")
	if err != nil {
		log.Fatalf("QueryDeviceByDateRange failed: %v", err)
	}
	fmt.Printf(f)

}
