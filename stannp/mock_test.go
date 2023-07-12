package stannp

import (
	"context"
	"io"
	"reflect"
	"testing"

	"github.com/CopilotIQ/stannp-client-golang/address"
	"github.com/CopilotIQ/stannp-client-golang/letter"
	"github.com/CopilotIQ/stannp-client-golang/util"
	"github.com/jgroeneveld/trial/assert"
)

func TestNewMockClient(t *testing.T) {
	tests := []struct {
		name   string
		opts   []MockOption
		expect MockClient
	}{
		{
			name:   "no options",
			expect: MockClient{},
		},
		{
			name: "with validateAddressFailNext",
			opts: []MockOption{
				WithValidateAddressFailNext(true),
			},
			expect: MockClient{validateAddressFailNext: true},
		},
		{
			name: "with savePDFContentsFailNext",
			opts: []MockOption{
				WithSavePDFContentsFailNext(true),
			},
			expect: MockClient{savePDFContentsFailNext: true},
		},
		{
			name: "with codeNext",
			opts: []MockOption{
				WithCodeNext(400),
			},
			expect: MockClient{codeNext: 400},
		},
		{
			name: "with getPDFContentsFailNext",
			opts: []MockOption{
				WithGetPDFContentsFailNext(true),
			},
			expect: MockClient{getPDFContentsFailNext: true},
		},
		{
			name: "with errorMessageNext",
			opts: []MockOption{
				WithErrorMessageNext("error"),
			},
			expect: MockClient{errorMessageNext: "error"},
		},
		{
			name: "with addressInvalidNext",
			opts: []MockOption{
				WithAddressInvalidNext(true),
			},
			expect: MockClient{addressInvalidNext: true},
		},
		{
			name: "with letterFailNext",
			opts: []MockOption{
				WithLetterFailNext(true),
			},
			expect: MockClient{letterFailNext: true},
		},
		{
			name: "with all options",
			opts: []MockOption{
				WithValidateAddressFailNext(true),
				WithSavePDFContentsFailNext(true),
				WithCodeNext(400),
				WithGetPDFContentsFailNext(true),
				WithErrorMessageNext("simulated error"),
				WithAddressInvalidNext(true),
				WithLetterFailNext(true),
			},
			expect: MockClient{
				validateAddressFailNext: true,
				savePDFContentsFailNext: true,
				codeNext:                400,
				getPDFContentsFailNext:  true,
				errorMessageNext:        "simulated error",
				addressInvalidNext:      true,
				letterFailNext:          true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewMockClient(tt.opts...)
			assert.Equal(t, tt.expect, *client)
		})
	}
}

func TestMockBytesToPDF(t *testing.T) {
	tests := []struct {
		name              string
		mockClientOptions []MockOption
		expectedSuccess   bool
		expectedError     *util.APIError
	}{
		{
			name:              "success expected err not expected",
			mockClientOptions: []MockOption{},
			expectedSuccess:   true,
			expectedError:     nil,
		},
		{
			name:              "success not expected err expected",
			mockClientOptions: []MockOption{WithSavePDFContentsFailNext(true)},
			expectedSuccess:   false,
			expectedError:     util.BuildError(500, "savePDFContentsFailNext is true"),
		},
		{
			name: "err expected code expected custom err expected",
			mockClientOptions: []MockOption{
				WithCodeNext(404),
				WithErrorMessageNext("custom message"),
				WithSavePDFContentsFailNext(true),
			},
			expectedSuccess: false,
			expectedError:   util.BuildError(404, "custom message"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := NewMockClient(tt.mockClientOptions...)
			localFile, apiErr := mockClient.SavePDFContents(nil)

			if tt.expectedError != nil {
				assert.NotNil(t, apiErr)
				assert.Equal(t, *tt.expectedError, *apiErr)
				assert.True(t, reflect.ValueOf(localFile).IsNil())
			} else {
				assert.True(t, reflect.ValueOf(apiErr).IsNil())
				assert.NotNil(t, localFile)
			}
		})
	}
}

func TestMockDownloadPDF(t *testing.T) {
	tests := []struct {
		name              string
		mockClientOptions []MockOption
		expectedSuccess   bool
		expectedError     *util.APIError
	}{
		{
			name:              "success expected err not expected",
			mockClientOptions: []MockOption{},
			expectedSuccess:   true,
			expectedError:     nil,
		},
		{
			name:              "success not expected err expected",
			mockClientOptions: []MockOption{WithGetPDFContentsFailNext(true)},
			expectedSuccess:   false,
			expectedError:     util.BuildError(500, "getPDFContentsFailNext is true"),
		},
		{
			name: "err expected code expected custom err expected",
			mockClientOptions: []MockOption{
				WithCodeNext(404),
				WithErrorMessageNext("custom message"),
				WithGetPDFContentsFailNext(true),
			},
			expectedSuccess: false,
			expectedError:   util.BuildError(404, "custom message"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := NewMockClient(tt.mockClientOptions...)
			pdfURL := util.RandomString(10)
			downloadPDFRes, apiErr := mockClient.GetPDFContents(context.Background(), pdfURL)

			if tt.expectedError != nil {
				assert.NotNil(t, apiErr)
				assert.Equal(t, *tt.expectedError, *apiErr)
				assert.True(t, reflect.ValueOf(downloadPDFRes).IsNil())
			} else {
				assert.True(t, reflect.ValueOf(apiErr).IsNil())
				assert.NotNil(t, downloadPDFRes)

				assert.Equal(t, downloadPDFRes.Name, pdfURL)

				byteArray, err := io.ReadAll(downloadPDFRes.Contents)
				assert.Nil(t, err)

				err = downloadPDFRes.Contents.Close()
				assert.Nil(t, err)

				assert.Equal(t, len([]byte(pdfURL)), len(byteArray))
			}
		})
	}
}

func TestMockSendLetter(t *testing.T) {
	tests := []struct {
		name              string
		mockClientOptions []MockOption
		expectedSuccess   bool
		expectedError     *util.APIError
	}{
		{
			name:              "success expected err not expected",
			mockClientOptions: []MockOption{},
			expectedSuccess:   true,
			expectedError:     nil,
		},
		{
			name:              "success not expected err expected",
			mockClientOptions: []MockOption{WithLetterFailNext(true)},
			expectedSuccess:   false,
			expectedError:     util.BuildError(500, "letterFailNext is true"),
		},
		{
			name: "err expected code expected custom err expected",
			mockClientOptions: []MockOption{
				WithCodeNext(404),
				WithErrorMessageNext("custom message"),
				WithLetterFailNext(true),
			},
			expectedSuccess: false,
			expectedError:   util.BuildError(404, "custom message"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := NewMockClient(tt.mockClientOptions...)
			sendLetterRes, apiErr := mockClient.SendLetter(context.Background(), &letter.SendReq{})

			if tt.expectedError != nil {
				assert.NotNil(t, apiErr)
				assert.Equal(t, *tt.expectedError, *apiErr)
				assert.True(t, reflect.ValueOf(sendLetterRes).IsNil())
			} else {
				assert.True(t, reflect.ValueOf(apiErr).IsNil())
				assert.NotNil(t, sendLetterRes)
				assert.Equal(t, tt.expectedSuccess, sendLetterRes.Success)

				assert.Equal(t, sendLetterRes.Data.Status, "received")
				assert.True(t, sendLetterRes.Data.Cost != "")
				assert.True(t, sendLetterRes.Data.Created != "")
				assert.True(t, sendLetterRes.Data.Format != "")
				assert.True(t, sendLetterRes.Data.Id == "0")
				assert.True(t, sendLetterRes.Data.PDFURL != "")
				assert.True(t, sendLetterRes.Data.Status != "")
			}
		})
	}
}

func TestMockValidateAddress(t *testing.T) {
	tests := []struct {
		name              string
		mockClientOptions []MockOption
		isValidExpected   bool
		errExpected       *util.APIError
	}{
		{
			name:              "valid expected err not expected",
			mockClientOptions: []MockOption{},
			isValidExpected:   true,
			errExpected:       nil,
		},
		{
			name:              "valid not expected err not expected",
			mockClientOptions: []MockOption{WithAddressInvalidNext(true)},
			isValidExpected:   false,
			errExpected:       nil,
		},
		{
			name:              "err expected",
			mockClientOptions: []MockOption{WithValidateAddressFailNext(true)},
			isValidExpected:   false,
			errExpected:       util.BuildError(500, "validateAddressFailNext is true"),
		},
		{
			name: "fail next code next err next",
			mockClientOptions: []MockOption{
				WithValidateAddressFailNext(true),
				WithCodeNext(400),
				WithErrorMessageNext("custom message"),
			},
			isValidExpected: false,
			errExpected:     util.BuildError(400, "custom message"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := NewMockClient(tt.mockClientOptions...)
			validateAddressRes, apiErr := mockClient.ValidateAddress(context.Background(), &address.ValidateReq{})

			if tt.errExpected != nil {
				assert.NotNil(t, apiErr)
				assert.Equal(t, *tt.errExpected, *apiErr)
			} else {
				assert.True(t, reflect.ValueOf(apiErr).IsNil())

				if tt.isValidExpected {
					assert.True(t, validateAddressRes.Data.IsValid)
				} else {
					assert.False(t, validateAddressRes.Data.IsValid)
				}
			}
		})
	}
}

func TestInterface(t *testing.T) {
	newReal := func() Client {
		return New()
	}

	newFake := func() Client {
		return NewMockClient()
	}

	assert.NotNil(t, newReal)
	assert.NotNil(t, newFake())
}
