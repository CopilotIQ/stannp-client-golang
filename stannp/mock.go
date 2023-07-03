package stannp

import (
	"github.com/CopilotIQ/stannp-client-golang/address"
	"github.com/CopilotIQ/stannp-client-golang/letter"
	"github.com/CopilotIQ/stannp-client-golang/util"
)

// Client interface is for mocking / testing. Implement it however you wish!
type Client interface {
	SendLetter(request *letter.SendReq) (*letter.SendRes, *util.APIError)
	ValidateAddress(request *address.ValidateReq) (*address.ValidateRes, *util.APIError)
}

type MockOption func(*MockClient)

type MockClient struct {
	invalidNext     bool
	codeNext        int
	errNext         string
	addressFailNext bool
	letterFailNext  bool
}

func WithAddressFailNext(failNext bool) MockOption {
	return func(c *MockClient) {
		c.addressFailNext = failNext
	}
}

func WithLetterFailNext(failNext bool) MockOption {
	return func(c *MockClient) {
		c.letterFailNext = failNext
	}
}

func WithInvalidNext(invalidNext bool) MockOption {
	return func(c *MockClient) {
		c.invalidNext = invalidNext
	}
}

func WithCodeNext(codeNext int) MockOption {
	return func(c *MockClient) {
		c.codeNext = codeNext
	}
}

func WithErrNext(errNext string) MockOption {
	return func(c *MockClient) {
		c.errNext = errNext
	}
}

func NewMockClient(opts ...MockOption) *MockClient {
	client := &MockClient{}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func (mc *MockClient) SendLetter(_ *letter.SendReq) (*letter.SendRes, *util.APIError) {
	if mc.letterFailNext {
		apiErr := &util.APIError{
			Code:    500,
			Error:   "letterFailNext is true",
			Success: false,
		}

		if mc.codeNext != 0 {
			apiErr.Code = mc.codeNext
		}

		if mc.errNext != "" {
			apiErr.Error = mc.errNext
		}

		return nil, apiErr
	}

	return &letter.SendRes{
		Data: letter.Data{
			Cost:    util.RandomString(10),
			Created: util.RandomString(10),
			Format:  util.RandomString(10),
			Id:      "0",
			PDF:     util.RandomString(10),
			Status:  "received",
		},
		Success: true,
	}, nil
}

func (mc *MockClient) ValidateAddress(_ *address.ValidateReq) (*address.ValidateRes, *util.APIError) {
	if mc.addressFailNext {
		apiErr := &util.APIError{
			Code:    500,
			Error:   "addressFailNext is true",
			Success: false,
		}

		if mc.codeNext != 0 {
			apiErr.Code = mc.codeNext
		}

		if mc.errNext != "" {
			apiErr.Error = mc.errNext
		}

		return nil, apiErr
	}

	validateRes := &address.ValidateRes{
		Data: address.Data{
			IsValid: true,
		},
		Success: true,
	}

	if mc.invalidNext {
		validateRes.Data.IsValid = false
	}

	return validateRes, nil
}
