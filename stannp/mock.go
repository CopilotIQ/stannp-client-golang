package stannp

import (
	"bytes"
	"context"
	"io"

	"github.com/copilotiq/stannp-client-golang/address"
	"github.com/copilotiq/stannp-client-golang/letter"
	"github.com/copilotiq/stannp-client-golang/util"
	"os"
)

type MockOption func(*MockClient)

type MockClient struct {
	addressInvalidNext      bool
	codeNext                int
	errorMessageNext        string
	getPDFContentsFailNext  bool
	getPDFResponseNext      *letter.PDFRes
	savePDFContentsFailNext bool
	sendLetterFailNext      bool
	sendLetterResponseNext  *letter.SendRes
	validateAddressFailNext bool
}

var _ Client = (*MockClient)(nil)

func WithAddressInvalidNext(invalidNext bool) MockOption {
	return func(c *MockClient) {
		c.addressInvalidNext = invalidNext
	}
}

func WithCodeNext(codeNext int) MockOption {
	return func(c *MockClient) {
		c.codeNext = codeNext
	}
}

func WithErrorMessageNext(errNext string) MockOption {
	return func(c *MockClient) {
		c.errorMessageNext = errNext
	}
}

func WithGetPDFContentsFailNext(failNext bool) MockOption {
	return func(c *MockClient) {
		c.getPDFContentsFailNext = failNext
	}
}

func WithGetPDFResponseNext(res *letter.PDFRes) MockOption {
	return func(c *MockClient) {
		c.getPDFResponseNext = res
	}
}

func WithSendLetterFailNext(failNext bool) MockOption {
	return func(c *MockClient) {
		c.sendLetterFailNext = failNext
	}
}

func WithSavePDFContentsFailNext(failNext bool) MockOption {
	return func(c *MockClient) {
		c.savePDFContentsFailNext = failNext
	}
}

func WithSendLetterResponseNext(res *letter.SendRes) MockOption {
	return func(c *MockClient) {
		c.sendLetterResponseNext = res
	}
}

func WithValidateAddressFailNext(failNext bool) MockOption {
	return func(c *MockClient) {
		c.validateAddressFailNext = failNext
	}
}

func NewMockClient(opts ...MockOption) *MockClient {
	client := &MockClient{}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func (mc *MockClient) GetPDFContents(_ context.Context, pdfURL string) (*letter.PDFRes, *util.APIError) {
	if mc.getPDFContentsFailNext {
		apiErr := util.BuildError(500, "getPDFContentsFailNext is true")

		if mc.codeNext != 0 {
			apiErr.Code = mc.codeNext
		}

		if mc.errorMessageNext != "" {
			apiErr.ErrorMessage = mc.errorMessageNext
		}

		return nil, apiErr
	}

	if mc.getPDFResponseNext != nil {
		return mc.getPDFResponseNext, nil
	}

	return &letter.PDFRes{
		Contents: io.NopCloser(bytes.NewBufferString(pdfURL)), // give the caller something to read if they want to
		Name:     pdfURL,
	}, nil
}

func (mc *MockClient) SavePDFContents(_ io.Reader) (*os.File, *util.APIError) {
	if mc.savePDFContentsFailNext {
		apiErr := util.BuildError(500, "savePDFContentsFailNext is true")

		if mc.codeNext != 0 {
			apiErr.Code = mc.codeNext
		}

		if mc.errorMessageNext != "" {
			apiErr.ErrorMessage = mc.errorMessageNext
		}

		return nil, apiErr
	}
	return &os.File{}, nil
}

func (mc *MockClient) SendLetter(_ context.Context, _ *letter.SendReq) (*letter.SendRes, *util.APIError) {
	if mc.sendLetterFailNext {
		apiErr := util.BuildError(500, "sendLetterFailNext is true")

		if mc.codeNext != 0 {
			apiErr.Code = mc.codeNext
		}

		if mc.errorMessageNext != "" {
			apiErr.ErrorMessage = mc.errorMessageNext
		}

		return nil, apiErr
	}

	if mc.sendLetterResponseNext != nil {
		return mc.sendLetterResponseNext, nil
	}

	return &letter.SendRes{
		Data: letter.Data{
			Cost:    util.RandomString(10),
			Created: util.RandomString(10),
			Format:  util.RandomString(10),
			ID:      "0",
			PDFURL:  util.RandomString(10),
			Status:  "received",
		},
		Success: true,
	}, nil
}

func (mc *MockClient) ValidateAddress(_ context.Context, _ *address.ValidateReq) (*address.ValidateRes, *util.APIError) {
	if mc.validateAddressFailNext {
		apiErr := util.BuildError(500, "validateAddressFailNext is true")

		if mc.codeNext != 0 {
			apiErr.Code = mc.codeNext
		}

		if mc.errorMessageNext != "" {
			apiErr.ErrorMessage = mc.errorMessageNext
		}

		return nil, apiErr
	}

	validateRes := &address.ValidateRes{
		Data: address.Data{
			IsValid: true,
		},
		Success: true,
	}

	if mc.addressInvalidNext {
		validateRes.Data.IsValid = false
	}

	return validateRes, nil
}
