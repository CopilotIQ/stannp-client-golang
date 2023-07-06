package stannp

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/CopilotIQ/stannp-client-golang/address"
	"github.com/CopilotIQ/stannp-client-golang/letter"
	"github.com/CopilotIQ/stannp-client-golang/util"
)

const APIKeyQSP = "api_key"
const BaseURL = "https://us.stannp.com/api/v1"
const ContentTypeHeaderKey = "Content-Type"
const CreateURL = "create"
const PDFURLPrefix = "https://us.stannp.com/api/v1/storage"
const URLEncodedHeaderVal = "application/x-www-form-urlencoded"
const ValidateURL = "validate"
const XIdempotenceyHeaderKey = "X-Idempotency-Key"

type Stannp struct {
	apiKey         string
	baseUrl        string
	clearZone      bool
	client         *http.Client
	duplex         bool
	postUnverified bool
	test           bool
}

type APIOption func(*Stannp)

func WithTest(test bool) APIOption {
	return func(s *Stannp) {
		s.test = test
	}
}

func WithAPIKey(apiKey string) APIOption {
	return func(s *Stannp) {
		s.apiKey = apiKey
	}
}

func WithPostUnverified(postUnverified bool) APIOption {
	return func(s *Stannp) {
		s.postUnverified = postUnverified
	}
}

func WithClearZone(clearZone bool) APIOption {
	return func(s *Stannp) {
		s.clearZone = clearZone
	}
}

func WithDuplex(duplex bool) APIOption {
	return func(s *Stannp) {
		s.duplex = duplex
	}
}

func WithHTTPClient(hc *http.Client) APIOption {
	return func(s *Stannp) {
		s.client = hc
	}
}

func New(options ...APIOption) *Stannp {
	api := &Stannp{
		apiKey:         "test123456",
		baseUrl:        BaseURL,
		clearZone:      true,
		client:         http.DefaultClient,
		duplex:         true,
		postUnverified: false,
		test:           true,
	}

	for _, option := range options {
		option(api)
	}
	return api
}

func (s *Stannp) PostUnverified() bool {
	return s.postUnverified
}

func (s *Stannp) IsTest() bool {
	return s.test
}

func (s *Stannp) wrapAuth(inputURL string) (string, *util.APIError) {
	u, err := url.Parse(inputURL)
	if err != nil {
		return "", util.BuildError(500, fmt.Sprintf("error parsing inputURL [%s]", inputURL))
	}

	q := u.Query()
	q.Set(APIKeyQSP, s.apiKey)
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func (s *Stannp) post(ctx context.Context, inputReader io.Reader, inputURL, idempotenceyHeaderVal string) (*http.Response, *util.APIError) {
	authURL, wrapErr := s.wrapAuth(inputURL)
	if wrapErr != nil {
		return nil, wrapErr
	}

	req, err := http.NewRequestWithContext(ctx, "POST", authURL, inputReader)
	if err != nil {
		return nil, util.BuildError(500, fmt.Sprintf("error generating POST req [%+v] for req [%+v]", err, req))
	}

	req.Header.Set(ContentTypeHeaderKey, URLEncodedHeaderVal)

	if idempotenceyHeaderVal != "" {
		req.Header.Set(XIdempotenceyHeaderKey, idempotenceyHeaderVal)
	}

	res, err := s.client.Do(req)
	if err != nil {
		return nil, util.BuildError(500, fmt.Sprintf("error sending req [%+v]", req))
	}

	return res, nil
}

func (s *Stannp) GetPDFContents(ctx context.Context, pdfURL string) (*letter.PDFRes, *util.APIError) {
	if !strings.HasPrefix(pdfURL, PDFURLPrefix) {
		return nil, util.BuildError(400, fmt.Sprintf("pdfURL must begin with [%s]. your input was [%s]", PDFURLPrefix, pdfURL))
	}

	fileURL, err := url.Parse(pdfURL)
	if err != nil {
		return nil, util.BuildError(500, err.Error())
	}
	path := fileURL.Path
	urlSegments := strings.Split(path, "/")
	fileName := urlSegments[len(urlSegments)-1]

	pdfGetReq, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, pdfURL, nil)
	if reqErr != nil {
		return nil, util.BuildError(500, reqErr.Error())
	}

	resp, err := s.client.Do(pdfGetReq)
	if err != nil {
		return nil, util.BuildError(500, err.Error())
	}

	return &letter.PDFRes{
		Contents: resp.Body,
		Name:     fileName,
	}, nil
}

func (s *Stannp) SavePDFContents(pdfContents io.ReadCloser) (*os.File, *util.APIError) {
	tmpFile, err := os.CreateTemp("", "stannp_letter.*.pdf")
	if err != nil {
		return nil, util.BuildError(500, err.Error())
	}

	_, copyErr := io.Copy(tmpFile, pdfContents)
	if copyErr != nil {
		removeErr := os.Remove(tmpFile.Name())
		if removeErr != nil {
			return nil, util.BuildError(500, removeErr.Error())
		}

		return nil, util.BuildError(500, copyErr.Error())
	}

	return tmpFile, nil
}

func (s *Stannp) SendLetter(ctx context.Context, request *letter.SendReq) (*letter.SendRes, *util.APIError) {
	formData := url.Values{}
	formData.Set("clearzone", strconv.FormatBool(s.clearZone))
	formData.Set("duplex", strconv.FormatBool(s.duplex))
	formData.Set("post_unverified", strconv.FormatBool(s.postUnverified))
	formData.Set("recipient[address1]", request.Recipient.Address1)
	formData.Set("recipient[address2]", request.Recipient.Address2)
	formData.Set("recipient[country]", request.Recipient.Country)
	formData.Set("recipient[firstname]", request.Recipient.Firstname)
	formData.Set("recipient[lastname]", request.Recipient.Lastname)
	formData.Set("recipient[state]", request.Recipient.State)
	formData.Set("recipient[title]", request.Recipient.Title)
	formData.Set("recipient[town]", request.Recipient.Town)
	formData.Set("recipient[zipcode]", request.Recipient.Zipcode)
	formData.Set("template", request.Template)
	formData.Set("test", strconv.FormatBool(s.test))

	// set custom merge variables in the formData
	for key, value := range request.MergeVariables {
		formData.Set("recipient["+key+"]", value)
	}

	res, postErr := s.post(ctx, strings.NewReader(formData.Encode()), strings.Join([]string{s.baseUrl, letter.URL, CreateURL}, "/"), request.IdempotenceyKey)
	if postErr != nil {
		return nil, postErr
	}

	var letterRes letter.SendRes
	resErr := util.ResToType(res.StatusCode, res.Body, &letterRes)
	return &letterRes, resErr
}

func (s *Stannp) ValidateAddress(ctx context.Context, request *address.ValidateReq) (*address.ValidateRes, *util.APIError) {
	// Create URL values
	formData := url.Values{}
	formData.Set("company", request.Company)
	formData.Set("address1", request.Address1)
	formData.Set("address2", request.Address2)
	formData.Set("city", request.City)
	formData.Set("zipcode", request.Zipcode)
	formData.Set("country", request.Country)

	res, postErr := s.post(ctx, strings.NewReader(formData.Encode()), strings.Join([]string{s.baseUrl, address.URL, ValidateURL}, "/"), "")
	if postErr != nil {
		return nil, postErr
	}

	var addressRes address.ValidateRes
	resErr := util.ResToType(res.StatusCode, res.Body, &addressRes)
	return &addressRes, resErr
}
