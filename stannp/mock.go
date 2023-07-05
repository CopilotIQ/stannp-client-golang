package stannp

import (
	"github.com/CopilotIQ/stannp-client-golang/address"
	"github.com/CopilotIQ/stannp-client-golang/letter"
	"github.com/CopilotIQ/stannp-client-golang/util"
	"os"
)

// Client interface is for mocking / testing. Implement it however you wish!
type Client interface {
	BytesToPDF(data []byte) (*os.File, *util.APIError)
	DownloadPDF(urlInput string) (*letter.PDFRes, *util.APIError)
	SendLetter(request *letter.SendReq) (*letter.SendRes, *util.APIError)
	ValidateAddress(request *address.ValidateReq) (*address.ValidateRes, *util.APIError)
}

type MockOption func(*MockClient)

type MockClient struct {
	addressFailNext     bool
	bytesToPDFFailNext  bool
	codeNext            int
	downloadPDFFailNext bool
	errNext             string
	invalidNext         bool
	letterFailNext      bool
}

func WithAddressFailNext(failNext bool) MockOption {
	return func(c *MockClient) {
		c.addressFailNext = failNext
	}
}

func WithBytesToPDFFailNext(failNext bool) MockOption {
	return func(c *MockClient) {
		c.bytesToPDFFailNext = failNext
	}
}

func WithCodeNext(codeNext int) MockOption {
	return func(c *MockClient) {
		c.codeNext = codeNext
	}
}

func WithDownloadPDFFailNext(failNext bool) MockOption {
	return func(c *MockClient) {
		c.downloadPDFFailNext = failNext
	}
}

func WithErrNext(errNext string) MockOption {
	return func(c *MockClient) {
		c.errNext = errNext
	}
}

func WithInvalidNext(invalidNext bool) MockOption {
	return func(c *MockClient) {
		c.invalidNext = invalidNext
	}
}

func WithLetterFailNext(failNext bool) MockOption {
	return func(c *MockClient) {
		c.letterFailNext = failNext
	}
}

func NewMockClient(opts ...MockOption) *MockClient {
	client := &MockClient{}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func (mc *MockClient) BytesToPDF(_ []byte) (*os.File, *util.APIError) {
	if mc.bytesToPDFFailNext {
		apiErr := util.BuildError(500, "bytesToPDFFailNext is true")

		if mc.codeNext != 0 {
			apiErr.Code = mc.codeNext
		}

		if mc.errNext != "" {
			apiErr.Error = mc.errNext
		}

		return nil, apiErr
	}
	return &os.File{}, nil
}

func (mc *MockClient) DownloadPDF(_ string) (*letter.PDFRes, *util.APIError) {
	if mc.downloadPDFFailNext {
		apiErr := util.BuildError(500, "downloadPDFFailNext is true")

		if mc.codeNext != 0 {
			apiErr.Code = mc.codeNext
		}

		if mc.errNext != "" {
			apiErr.Error = mc.errNext
		}

		return nil, apiErr
	}

	return &letter.PDFRes{
		Bytes: []byte("hi sean"),
		Len:   len([]byte("hi sean")),
		Name:  "hi sean.pdf",
	}, nil
}

func (mc *MockClient) SendLetter(_ *letter.SendReq) (*letter.SendRes, *util.APIError) {
	if mc.letterFailNext {
		apiErr := util.BuildError(500, "letterFailNext is true")

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
			PDFURL:  util.RandomString(10),
			Status:  "received",
		},
		Success: true,
	}, nil
}

func (mc *MockClient) ValidateAddress(_ *address.ValidateReq) (*address.ValidateRes, *util.APIError) {
	if mc.addressFailNext {
		apiErr := util.BuildError(500, "addressFailNext is true")

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
