package stannp

import (
	"github.com/jgroeneveld/trial/assert"
	"github.com/joho/godotenv"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	// Load .env file
	err := godotenv.Load()
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
