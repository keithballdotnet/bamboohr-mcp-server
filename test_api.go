package main

import (
	"fmt"
	"os"
)

func testTimeOffRequests() {
	company := os.Getenv("BAMBOOHR_COMPANY")
	apiKey := os.Getenv("BAMBOOHR_API_KEY")

	if company == "" || apiKey == "" {
		fmt.Println("Please set BAMBOOHR_COMPANY and BAMBOOHR_API_KEY environment variables")
		return
	}

	client := NewBambooHRClient(company, apiKey)

	// Test the time-off requests call with debug output
	fmt.Println("Testing time-off requests API call...")
	requests, err := client.GetTimeOffRequests(157, "", "")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Success! Got %d request records\n", len(requests))
		for i, request := range requests[:3] { // Show first 3 requests
			fmt.Printf("Request %d: ID=%s, Start=%s, End=%s, Type=%s, Amount=%v\n",
				i, request.ID, request.Start, request.End, request.Type.Name, request.Amount.Amount)
		}
	}
}

func testMain() {
	if len(os.Args) > 1 && os.Args[1] == "test-requests" {
		testTimeOffRequests()
		return
	}

	// Original main function logic would go here
	fmt.Println("Use 'go run *.go test-requests' to test the API")
}
