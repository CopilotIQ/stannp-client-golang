# Stannp Go API Client

This is a Go client for the Stannp API. It provides a simple way to interact with the Stannp API, allowing you to send letters programmatically.

## Features

- Initialize a new API client with your API key, base URL, and post unverified preference.
- Send a letter via the Stannp API, with recipient details and other parameters specified in a `SendLetterRequest` struct.
- Receive the response from the Stannp API as a `SendLetterResponse` struct, which includes details like the PDF URL, ID, creation time, format, cost, and status of the letter.

## Usage

To create a new API client, use the `New` function with the `WithAPIKey`, `WithBaseURL`, and `WithPostUnverified` options:

```go
package main

import "github.com/copilotiq/stannp"

func main() {
	api := stannp.New(
		stannp.WithAPIKey("your-api-key"),
		stannp.WithBaseURL("https://us.stannp.com/api/v1"),
		stannp.WithPostUnverified(true),
	)

	request := stannp.SendLetterRequest{
		Test:      true,
		Template:  305202,
		Clearzone: true,
		Duplex:    true,
		Recipient: stannp.RecipientDetails{
			Title:      "Mr.",
			Firstname:  "First",
			Lastname:   "Last",
			Address1:   "Address",
			Town:       "City",
			Zipcode:    "Zip",
			State:      "State",
			Country:    "US",
		},
	}
	response, err := api.SendLetter(request)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Response: %+v\n", response)
	
}
```

This README provides a brief overview of the client's features and shows how to initialize a new API client and send a letter. Always make sure to replace `"your-api-key"` with your actual Stannp API key when initializing the API client.
