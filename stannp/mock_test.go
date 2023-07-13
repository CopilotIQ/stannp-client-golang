package stannp

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"reflect"
	"testing"

	"github.com/CopilotIQ/stannp-client-golang/address"
	"github.com/CopilotIQ/stannp-client-golang/letter"
	"github.com/CopilotIQ/stannp-client-golang/util"
	"github.com/jgroeneveld/trial/assert"
)

func TestMockClient_GetPDFContents(t *testing.T) {
	tests := []struct {
		name              string
		mockClientOptions []MockOption
		expectedSuccess   bool
		expectedError     *util.APIError
	}{
		{
			name: "success expected err not expected",
			mockClientOptions: []MockOption{WithGetPDFResponseNext(
				&letter.PDFRes{
					Contents: io.NopCloser(bytes.NewBufferString("success expected err not expected")),
					Name:     "success expected err not expected",
				},
			)},
			expectedSuccess: true,
			expectedError:   nil,
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

				assert.Equal(t, downloadPDFRes.Name, "success expected err not expected")

				b1, err := io.ReadAll(downloadPDFRes.Contents)
				assert.Nil(t, err)

				err = downloadPDFRes.Contents.Close()
				assert.Nil(t, err)

				b2, err := io.ReadAll(bytes.NewBufferString("success expected err not expected"))
				assert.Nil(t, err)

				assert.True(t, reflect.DeepEqual(b1, b2))
			}
		})
	}
}

func TestMockClient_NewMockClient(t *testing.T) {
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
			name: "with addressInvalidNext",
			opts: []MockOption{
				WithAddressInvalidNext(true),
			},
			expect: MockClient{addressInvalidNext: true},
		},
		{
			name: "with codeNext",
			opts: []MockOption{
				WithCodeNext(400),
			},
			expect: MockClient{codeNext: 400},
		},
		{
			name: "with errorMessageNext",
			opts: []MockOption{
				WithErrorMessageNext("error"),
			},
			expect: MockClient{errorMessageNext: "error"},
		},
		{
			name: "with getPDFContentsFailNext",
			opts: []MockOption{
				WithGetPDFContentsFailNext(true),
			},
			expect: MockClient{getPDFContentsFailNext: true},
		},
		{
			name: "with getPDFResponseNext",
			opts: []MockOption{
				WithGetPDFResponseNext(
					&letter.PDFRes{
						Contents: io.NopCloser(bytes.NewBufferString("with getPDFResponseNext")),
						Name:     "with getPDFResponseNext",
					},
				),
			},
			expect: MockClient{getPDFResponseNext: &letter.PDFRes{
				Contents: io.NopCloser(bytes.NewBufferString("with getPDFResponseNext")),
				Name:     "with getPDFResponseNext",
			}},
		},
		{
			name: "with sendLetterFailNext",
			opts: []MockOption{
				WithSendLetterFailNext(true),
			},
			expect: MockClient{sendLetterFailNext: true},
		},
		{
			name: "with sendLetterResponseNext",
			opts: []MockOption{
				WithSendLetterResponseNext(
					&letter.SendRes{
						Data: letter.Data{
							Cost:    "1",
							Created: "2",
							Format:  "3",
							ID:      "4",
							PDFURL:  "5",
							Status:  "6",
						},
						Success: true,
					},
				),
			},
			expect: MockClient{
				sendLetterResponseNext: &letter.SendRes{
					Data: letter.Data{
						Cost:    "1",
						Created: "2",
						Format:  "3",
						ID:      "4",
						PDFURL:  "5",
						Status:  "6",
					},
					Success: true,
				},
			},
		},
		{
			name: "with savePDFContentsFailNext",
			opts: []MockOption{
				WithSavePDFContentsFailNext(true),
			},
			expect: MockClient{savePDFContentsFailNext: true},
		},
		{
			name: "with validateAddressFailNext",
			opts: []MockOption{
				WithValidateAddressFailNext(true),
			},
			expect: MockClient{validateAddressFailNext: true},
		},
		{
			name: "with all options",
			opts: []MockOption{
				WithAddressInvalidNext(true),
				WithCodeNext(400),
				WithErrorMessageNext("simulated error"),
				WithGetPDFContentsFailNext(true),
				WithGetPDFResponseNext(
					&letter.PDFRes{
						Contents: io.NopCloser(bytes.NewBufferString("with getPDFResponseNext")),
						Name:     "with getPDFResponseNext",
					},
				),
				WithSendLetterFailNext(true),
				WithSendLetterResponseNext(
					&letter.SendRes{
						Data: letter.Data{
							Cost:    "7",
							Created: "8",
							Format:  "9",
							ID:      "10",
							PDFURL:  "11",
							Status:  "12",
						},
						Success: true,
					}),
				WithSavePDFContentsFailNext(true),
				WithValidateAddressFailNext(true),
			},
			expect: MockClient{
				addressInvalidNext:     true,
				codeNext:               400,
				errorMessageNext:       "simulated error",
				getPDFContentsFailNext: true,
				getPDFResponseNext: &letter.PDFRes{
					Contents: io.NopCloser(bytes.NewBufferString("with getPDFResponseNext")),
					Name:     "with getPDFResponseNext",
				},
				sendLetterFailNext: true,
				sendLetterResponseNext: &letter.SendRes{
					Data: letter.Data{
						Cost:    "7",
						Created: "8",
						Format:  "9",
						ID:      "10",
						PDFURL:  "11",
						Status:  "12",
					},
					Success: true,
				},
				savePDFContentsFailNext: true,
				validateAddressFailNext: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewMockClient(tt.opts...)
			assert.Equal(t, tt.expect.addressInvalidNext, client.addressInvalidNext)
			assert.Equal(t, tt.expect.codeNext, client.codeNext)
			assert.Equal(t, tt.expect.errorMessageNext, client.errorMessageNext)
			assert.Equal(t, tt.expect.getPDFContentsFailNext, client.getPDFContentsFailNext)
			assert.Equal(t, tt.expect.savePDFContentsFailNext, client.savePDFContentsFailNext)
			assert.Equal(t, tt.expect.sendLetterFailNext, client.sendLetterFailNext)
			assert.Equal(t, tt.expect.validateAddressFailNext, client.validateAddressFailNext)

			if tt.expect.sendLetterResponseNext != nil {
				assert.True(t, reflect.DeepEqual(*tt.expect.sendLetterResponseNext, *client.sendLetterResponseNext))
			}

			if tt.expect.getPDFResponseNext != nil {
				b1, readErr := io.ReadAll(tt.expect.getPDFResponseNext.Contents)
				assert.Nil(t, readErr)

				b2, readErr := io.ReadAll(client.getPDFResponseNext.Contents)
				assert.Nil(t, readErr)

				assert.True(t, reflect.DeepEqual(b1, b2))
			}
		})
	}
}

func TestMockClient_SavePDFContents(t *testing.T) {
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

func TestMockClient_SendLetter(t *testing.T) {
	tests := []struct {
		name              string
		mockClientOptions []MockOption
		expectedSuccess   bool
		expectedError     *util.APIError
	}{
		{
			name:              "success not expected err expected",
			mockClientOptions: []MockOption{WithSendLetterFailNext(true)},
			expectedSuccess:   false,
			expectedError:     util.BuildError(500, "sendLetterFailNext is true"),
		},
		{
			name: "success expected with letter res pre-defined",
			mockClientOptions: []MockOption{WithSendLetterResponseNext(
				&letter.SendRes{
					Data: letter.Data{
						Cost:    "1",
						Created: "2",
						Format:  "3",
						ID:      "4",
						PDFURL:  "5",
						Status:  "6",
					},
					Success: true,
				},
			),
			},
			expectedSuccess: true,
			expectedError:   nil,
		},
		{
			name: "err expected code expected custom err expected",
			mockClientOptions: []MockOption{
				WithCodeNext(404),
				WithErrorMessageNext("custom message"),
				WithSendLetterFailNext(true),
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

				assert.Equal(t, sendLetterRes.Data.Cost, "1")
				assert.Equal(t, sendLetterRes.Data.Created, "2")
				assert.Equal(t, sendLetterRes.Data.Format, "3")
				assert.Equal(t, sendLetterRes.Data.ID, json.Number("4"))
				assert.Equal(t, sendLetterRes.Data.PDFURL, "5")
				assert.Equal(t, sendLetterRes.Data.Status, "6")
			}
		})
	}
}

func TestMockClient_ValidateAddress(t *testing.T) {
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
