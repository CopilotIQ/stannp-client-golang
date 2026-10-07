package stannp

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/copilotiq/stannp-client-golang/address"
	"github.com/copilotiq/stannp-client-golang/letter"
	"github.com/jgroeneveld/trial/assert"
)

const credentialTestKey = "fake-stannp-key-for-tests"

// receivedRequest is what the fake Stannp server saw.
type receivedRequest struct {
	path, rawQuery, form, basicUser, basicPass string
	basicOK                                    bool
}

func fakeStannp(t *testing.T, got *receivedRequest) *httptest.Server {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		got.path, got.rawQuery, got.form = r.URL.Path, r.URL.RawQuery, r.PostForm.Encode()
		got.basicUser, got.basicPass, got.basicOK = r.BasicAuth()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"id":0,"is_valid":true}}`))
	}))
	t.Cleanup(ts.Close)
	return ts
}

func TestStannpCalls_SendKeyAsBasicAuthNeverInURLOrBody(t *testing.T) {
	calls := []struct {
		name     string
		wantPath string
		call     func(*Stannp) error
	}{
		{"SendLetter", "/letters/create", func(s *Stannp) error {
			_, err := s.SendLetter(context.Background(), &letter.SendReq{Template: "307051", IdempotenceyKey: "idem-1"})
			if err != nil {
				return err
			}
			return nil
		}},
		{"ValidateAddress", "/addresses/validate", func(s *Stannp) error {
			_, err := s.ValidateAddress(context.Background(), &address.ValidateReq{Address1: "9355 Burton Way", Country: "US"})
			if err != nil {
				return err
			}
			return nil
		}},
	}

	for _, tc := range calls {
		t.Run(tc.name, func(t *testing.T) {
			var got receivedRequest
			ts := fakeStannp(t, &got)
			api := New(WithAPIKey(credentialTestKey), WithHTTPClient(ts.Client()))
			api.baseUrl = ts.URL

			err := tc.call(api)

			assert.Nil(t, err)
			assert.Equal(t, tc.wantPath, got.path)
			assert.Equal(t, "", got.rawQuery, "no query string at all: the key must not ride in the URL")
			assert.False(t, strings.Contains(got.form, credentialTestKey), "the key must not be in the form body")
			assert.True(t, got.basicOK, "the request must carry Basic auth")
			assert.Equal(t, credentialTestKey, got.basicUser)
			assert.Equal(t, "", got.basicPass)
		})
	}
}

func TestSendLetter_TransportErrorNamesPathButNeverTheKey(t *testing.T) {
	ts := httptest.NewServer(http.NotFoundHandler())
	baseURL := ts.URL
	ts.Close() // nothing is listening any more: the client gets a connection error
	api := New(WithAPIKey(credentialTestKey))
	api.baseUrl = baseURL

	_, apiErr := api.SendLetter(context.Background(), &letter.SendReq{Template: "307051"})

	assert.NotNil(t, apiErr)
	msg := apiErr.String()
	encoded := base64.StdEncoding.EncodeToString([]byte(credentialTestKey + ":"))
	for _, secret := range []string{credentialTestKey, "api_key", encoded, "Authorization"} {
		assert.False(t, strings.Contains(msg, secret), "error text is logged and stored by callers; it leaked %q: %s", secret, msg)
	}
	assert.True(t, strings.Contains(msg, "POST /letters/create"), "the error must still say which call failed: %s", msg)
	assert.Equal(t, 500, apiErr.Code)
}

func TestSendLetter_HostileKeysTravelOnlyInBasicAuth(t *testing.T) {
	keys := []struct {
		name string
		key  string
	}{
		{"empty key", ""},
		{"key smuggling a query parameter", "k&api_key=x&test=false"},
		{"unicode key", "ключ-🔑-clé"},
		{"5k-char key", strings.Repeat("f", 5000)},
	}

	for _, tc := range keys {
		t.Run(tc.name, func(t *testing.T) {
			var got receivedRequest
			ts := fakeStannp(t, &got)
			api := New(WithAPIKey(tc.key), WithHTTPClient(ts.Client()))
			api.baseUrl = ts.URL

			_, apiErr := api.SendLetter(context.Background(), &letter.SendReq{Template: "307051"})

			assert.True(t, apiErr == nil, "unexpected error: %v", apiErr)
			assert.Equal(t, "", got.rawQuery)
			assert.Equal(t, tc.key, got.basicUser)
		})
	}
}
