package main

import (
	"fmt"
	"os"
)

func debugAPI() {
	company := os.Getenv("BAMBOOHR_COMPANY")
	apiKey := os.Getenv("BAMBOOHR_API_KEY")

	if company == "" || apiKey == "" {
		fmt.Println("Please set BAMBOOHR_COMPANY and BAMBOOHR_API_KEY environment variables")
		return
	}

	client := NewBambooHRClient(company, apiKey)

	// Test the time-off balance call
	fmt.Println("Testing time-off balance API call...")
	balances, err := client.GetTimeOffBalance(157) // Keith's ID
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Success! Got %d balance records\n", len(balances))
		for i, balance := range balances {
			fmt.Printf("Balance %d: %+v\n", i, balance)
		}
	}

	// Test the time-off requests call
	fmt.Println("\nTesting time-off requests API call...")
	requests, err := client.GetTimeOffRequests(157, "", "")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Success! Got %d request records\n", len(requests))
		for i, request := range requests {
			fmt.Printf("Request %d: %+v\n", i, request)
		}
	}
}
