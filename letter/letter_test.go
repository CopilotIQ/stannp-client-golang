package letter

import (
	"copilotiq/stannp-client-golang/stannp"
	"encoding/json"
	"testing"

	"github.com/nsf/jsondiff"
)

func TestSendLetter(t *testing.T) {
	// Initialize API with random api key
	api := stannp.New(
		stannp.WithAPIKey("random-api-key"),
		stannp.WithBaseURL("https://us.stannp.com/api/v1"),
		stannp.WithPostUnverified(true),
	)

	// Call SendLetter with a new instance of Request
	request := stannp.SendLetterRequest{
		Test:      true,
		Template:  305202,
		Clearzone: true,
		Duplex:    true,
		Recipient: stannp.RecipientDetails{
			Title:     "Mr.",
			Firstname: "John",
			Lastname:  "Doe",
			Address1:  "123 Random St",
			Town:      "Townsville",
			Zipcode:   "12345",
			State:     "Stateville",
			Country:   "US",
			Testvalue: 99.99,
			Honorific: "Mr.",
			Date:      "2025-01-01",
		},
	}

	// Note: This call is not actually sending a request.
	response, err := api.SendLetter(request)
	require.NoError(t, err)

	// Generate random expected response
	expectedResponse := &stannp.SendLetterResponse{
		Success: true,
		Data: struct {
			Pdf     string `json:"pdf"`
			Id      string `json:"id"`
			Created string `json:"created"`
			Format  string `json:"format"`
			Cost    string `json:"cost"`
			Status  string `json:"status"`
		}{
			Pdf:     "https://random.pdf",
			Id:      "0",
			Created: "2023-05-17T21:06:48+00:00",
			Format:  "US-LETTER",
			Cost:    "0.81",
			Status:  "test",
		},
	}

	// Convert to JSON
	responseJSON, err := json.Marshal(response)
	require.NoError(t, err)
	expectedJSON, err := json.Marshal(expectedResponse)
	require.NoError(t, err)

	// Compare using jsondiff
	differ := jsondiff.New()
	_, differences := differ.Compare(responseJSON, expectedJSON)
	if differences != jsondiff.FullMatch {
		t.Errorf("Response does not match expected: %s", differences)
	}
}
