package stannp

import (
	"copilotiq/stannp-client-golang/letter"
	"encoding/json"
	"github.com/jgroeneveld/trial/assert"
	"github.com/joho/godotenv"
	"github.com/nsf/jsondiff"
	"log"
	"os"
	"testing"
)

const ApiKeyEnvKey = "STANNP_API_KEY"

var TestClient *Stannp

func TestMain(m *testing.M) {
	setup()
	m.Run()
}

func setup() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Unable to load .env file: %s", err)
	}

	apiKey := os.Getenv(ApiKeyEnvKey)
	if apiKey == "" {
		log.Fatalf("Cannot proceed when apiKey is the empty string [%s]", apiKey)
	}

	// Initialize Stannp with test data
	TestClient = New(
		WithAPIKey(apiKey),
		WithPostUnverified(false),
		WithTest(true),
	)

	if !TestClient.IsTest() {
		log.Fatalf("Cannot proceed when API key is live [%s]", apiKey)
	}
}

//goland:noinspection GoBoolExpressions
func TestNew(t *testing.T) {
	// Load .env file
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatal("Error loading .env file")
	}

	// Get API key from environment variable
	envAPIKey, exists := os.LookupEnv("STANNP_API_KEY")
	if !exists {
		t.Fatal("STANNP_API_KEY not set in .env file")
	}

	// Test data
	postUnverified := false
	test := true

	// Initialize Stannp with test data
	api := New(
		WithAPIKey(envAPIKey),
		WithPostUnverified(postUnverified),
		WithTest(test),
	)

	// Assert that the Stannp client has been initialized with the correct values
	assert.Equal(t, envAPIKey, api.apiKey, "APIKey does not match expected")
	assert.Equal(t, BaseURL, api.baseUrl, "BaseURL does not match expected")
	assert.Equal(t, postUnverified, api.postUnverified, "PostUnverified does not match expected")
	assert.Equal(t, test, api.test, "Test does not match expected")
}

func TestSendLetter(t *testing.T) {
	// Call SendLetter with a new instance of Request
	request := letter.Request{
		Test:      true,
		Template:  305202,
		ClearZone: true,
		Duplex:    true,
		Recipient: letter.RecipientDetails{
			Title:     "Mr.",
			Firstname: "John",
			Lastname:  "Doe",
			Address1:  "123 Random St",
			Town:      "Townsville",
			Zipcode:   "12345",
			State:     "Stateville",
			Country:   "US",
		},
	}

	// Note: This call is not actually sending a request.
	response, err := TestClient.SendLetter(request)
	assert.Nil(t, err)

	expected := &letter.Response{
		Success: true,
		Data: letter.Data{
			Pdf:     "https://random.pdf",
			Id:      "0",
			Created: "2023-05-17T21:06:48+00:00",
			Format:  "US-LETTER",
			Cost:    "0.81",
			Status:  "test",
		},
	}

	// Convert to JSON
	responseJSON, jsonErr := json.Marshal(response)
	assert.Nil(t, jsonErr)
	expectedJSON, jsonErr := json.Marshal(expected)
	assert.Nil(t, jsonErr)

	difference, differences := jsondiff.Compare(responseJSON, expectedJSON, &jsondiff.Options{})
	if difference != jsondiff.FullMatch {
		t.Errorf("Response does not match expected: %s", differences)
	}
}
