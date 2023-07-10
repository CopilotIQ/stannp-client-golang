package stannp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/CopilotIQ/stannp-client-golang/address"
	"github.com/CopilotIQ/stannp-client-golang/letter"
	"github.com/CopilotIQ/stannp-client-golang/util"
	"github.com/jgroeneveld/trial/assert"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
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
		WithClearZone(false),
		WithDuplex(false),
		WithPostUnverified(false),
		WithTest(true),
	)

	if !TestClient.IsTest() {
		log.Fatalf("Cannot proceed when API key is live [%s]", apiKey)
	}
}

//goland:noinspection GoBoolExpressions
func TestNew(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatal("error loading .env file")
	}

	envAPIKey, exists := os.LookupEnv("STANNP_API_KEY")
	if !exists {
		t.Fatal("STANNP_API_KEY not set in .env file")
	}

	api := New(
		WithAPIKey(envAPIKey),
		WithClearZone(false),
		WithDuplex(false),
		WithPostUnverified(true),
		WithTest(true),
	)

	assert.Equal(t, envAPIKey, api.apiKey)
	assert.Equal(t, BaseURL, api.baseUrl)
	assert.Equal(t, false, api.clearZone)
	assert.Equal(t, false, api.duplex)
	assert.Equal(t, true, api.postUnverified)
	assert.Equal(t, true, api.test)
}

func TestPost(t *testing.T) {
	apiKey := util.RandomString(10)
	ctx := context.Background()
	inputBody := util.RandomString(10)
	inputReader := bytes.NewBufferString(inputBody)
	idempotenceyKey := util.RandomString(10)
	testURL := "/dashboard/u/1?docId=" + util.RandomString(10)
	urlParts := strings.Split(testURL, "?")

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		// make sure the original query string parameters are preserved AND the api key is injected
		assert.Equal(t, r.URL.String(), urlParts[0]+"?"+APIKeyQSP+"="+apiKey+"&"+urlParts[1])
		assert.True(t, reflect.DeepEqual(r.Header[ContentTypeHeaderKey], []string{URLEncodedHeaderVal}))
		assert.True(t, reflect.DeepEqual(r.Header[XIdempotenceyHeaderKey], []string{idempotenceyKey}))

		// in golang this closes the response writer automagically
		w.Header().Set("Content-Type", "application/json")
	}

	ts := httptest.NewServer(http.HandlerFunc(handler))

	api := New(
		WithAPIKey(apiKey),
		WithHTTPClient(ts.Client()),
	)

	if api.test != true {
		t.FailNow()
	}

	res, apiErr := api.post(ctx, inputReader, ts.URL+testURL, idempotenceyKey)
	assert.True(t, reflect.ValueOf(apiErr).IsNil())
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestStannp(t *testing.T) {
	t.Run("test SendLetter and verify the response is correct", func(t *testing.T) {
		request := &letter.SendReq{
			Template:        "307051",
			IdempotenceyKey: util.RandomString(10),
			Recipient: letter.RecipientDetails{
				Address1:  "9355 Burton Way",
				Address2:  "Courthouse",
				Country:   "United States",
				Firstname: "Judge",
				Lastname:  "Judy",
				State:     "CA",
				Title:     "Mrs.",
				Town:      "Beverly Hills",
				Zipcode:   "90210",
			},
		}

		response, sendErr := TestClient.SendLetter(context.Background(), request)
		assert.True(t, reflect.ValueOf(sendErr).IsNil())
		assert.NotNil(t, response)

		dateString := time.Now().Format("2006-01-02")

		assert.Equal(t, response.Data.Cost, "0.84")
		assert.Equal(t, response.Data.Format, "US-LETTER")
		assert.Equal(t, response.Data.Id.String(), "0")
		assert.Equal(t, response.Data.Status, "test")
		assert.True(t, response.Success)
		assert.True(t, strings.HasPrefix(response.Data.Created, dateString))
		assert.True(t, strings.HasPrefix(response.Data.PDFURL, "https://us.stannp.com/api/v1/storage/get/"))

		t.Run("test GetPDFContents and verify the response is correct", func(t *testing.T) {
			pdfRes, getPDFErr := TestClient.GetPDFContents(context.Background(), response.Data.PDFURL)
			assert.True(t, reflect.ValueOf(getPDFErr).IsNil())
			assert.NotNil(t, pdfRes)

			t.Run("test SavePDFContents and verify the response is correct", func(t *testing.T) {
				fileRes, fileErr := TestClient.SavePDFContents(pdfRes.Contents)
				assert.True(t, reflect.ValueOf(fileErr).IsNil())
				defer func() {
					removeErr := os.Remove(fileRes.Name())
					if removeErr != nil {
						panic(fmt.Sprintf("Error deleting temp file while testing. Please verify [%s] does not exist on your local disk.", fileRes.Name()))
					}
				}()

				content, err := os.ReadFile(fileRes.Name())
				assert.Nil(t, err)

				// Compare the length of the original data versus the expected return result to see if the PDF changed
				assert.Equal(t, 624899, len(content))
			})
		})
	})

	t.Run("test ValidateAddress and verify is_valid is false for fake data", func(t *testing.T) {
		request := &address.ValidateReq{
			Address1: "9354444445 Burton Way",
			City:     "Beverly Hills",
			Company:  "Beverly Hills Courthouse",
			Country:  "US",
			State:    "CA",
			Zipcode:  "90210",
		}

		validateRes, apiErr := TestClient.ValidateAddress(context.Background(), request)
		assert.True(t, reflect.ValueOf(apiErr).IsNil())
		assert.False(t, validateRes.Data.IsValid)
		assert.True(t, validateRes.Success)
	})
	t.Run("test ValidateAddress and verify is_valid is true for real data", func(t *testing.T) {
		request := &address.ValidateReq{
			Address1: "9355 Burton Way",
			City:     "Beverly Hills",
			Company:  "Beverly Hills Courthouse",
			Country:  "US",
			State:    "CA",
			Zipcode:  "90210",
		}

		validateRes, apiErr := TestClient.ValidateAddress(context.Background(), request)
		assert.True(t, reflect.ValueOf(apiErr).IsNil())
		assert.True(t, validateRes.Data.IsValid)
		assert.True(t, validateRes.Success)
	})
}

func TestJSONValuesUnmarshalWithCorrectFlexibility(t *testing.T) {
	t.Run("verify when Id is an int", func(t *testing.T) {
		rawJSON := `
{
  "data": {
    "cost": "10.99",
    "created": "2023-06-22",
    "format": "A4",
    "id": 12345,
    "pdf": "https://example.com/document.pdf",
    "status": "completed"
  },
  "success": true
}`
		var letterRes letter.SendRes
		jsonErr := json.Unmarshal([]byte(rawJSON), &letterRes)
		assert.Nil(t, jsonErr)

		int64Val, err := letterRes.Data.Id.Int64()
		assert.Nil(t, err)
		assert.Equal(t, "12345", letterRes.Data.Id.String())
		assert.Equal(t, int64(12345), int64Val)
	})
	t.Run("verify when Id is a string", func(t *testing.T) {
		rawJSON := `
{
  "data": {
    "cost": "10.99",
    "created": "2023-06-22",
    "format": "A4",
    "id": "12345",
    "pdf": "https://example.com/document.pdf",
    "status": "completed"
  },
  "success": true
}`

		var letterRes letter.SendRes
		jsonErr := json.Unmarshal([]byte(rawJSON), &letterRes)
		assert.Nil(t, jsonErr)

		int64Val, err := letterRes.Data.Id.Int64()
		assert.Nil(t, err)
		assert.Equal(t, "12345", letterRes.Data.Id.String())
		assert.Equal(t, int64(12345), int64Val)
	})
}
