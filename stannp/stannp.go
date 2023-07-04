package stannp

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/CopilotIQ/stannp-client-golang/address"
	"github.com/CopilotIQ/stannp-client-golang/letter"
	"github.com/CopilotIQ/stannp-client-golang/util"
	"github.com/google/uuid"
)

const PDFURLPrefix = "https://us.stannp.com/api/v1/storage"
const BaseURL = "https://us.stannp.com/api/v1"
const CreateURL = "create"
const ValidateURL = "validate"

type IdempotencyFunc func() string

var DefaultIdemFunc = func() string {
	guid, uuidErr := uuid.NewUUID()
	if uuidErr != nil {
		log.Fatalf("cannot generate UUID [%+v]", uuidErr)
	}

	return guid.String()
}

type Stannp struct {
	apiKey         string
	baseUrl        string
	clearZone      bool
	client         *http.Client
	duplex         bool
	idemFunc       IdempotencyFunc
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

func WithIdempotencyFunc(f IdempotencyFunc) APIOption {
	return func(s *Stannp) {
		s.idemFunc = f
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
	q.Set("api_key", s.apiKey)
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func (s *Stannp) post(inputReader io.Reader, inputURL string) (*http.Response, *util.APIError) {
	authURL, wrapErr := s.wrapAuth(inputURL)
	if wrapErr != nil {
		return nil, wrapErr
	}

	req, err := http.NewRequest("POST", authURL, inputReader)
	if err != nil {
		return nil, util.BuildError(500, fmt.Sprintf("error generating POST req [%+v] for req [%+v]", err, req))
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if s.idemFunc != nil {
		req.Header.Set("X-Idempotency-Key", s.idemFunc())
	}

	res, err := s.client.Do(req)
	if err != nil {
		return nil, util.BuildError(500, fmt.Sprintf("error sending req [%+v]", req))
	}

	return res, nil
}

func (s *Stannp) DownloadPDF(urlInput string) (*letter.PDFRes, *util.APIError) {
	if !strings.HasPrefix(urlInput, PDFURLPrefix) {
		return nil, util.BuildError(400, fmt.Sprintf("urlInput must begin with [%s]. your input was [%s]", PDFURLPrefix, urlInput))
	}

	// Build fileName from fullPath
	fileURL, err := url.Parse(urlInput)
	if err != nil {
		log.Fatal(err)
	}
	path := fileURL.Path
	segments := strings.Split(path, "/")
	fileName := segments[len(segments)-1]

	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	resp, err := client.Get(urlInput)
	if err != nil {
		log.Fatal(err)
	}

	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, util.BuildError(500, err.Error())
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, util.BuildError(500, err.Error())
	}

	return &letter.PDFRes{
		Bytes: byteArray,
		Len:   len(byteArray),
		Name:  fileName,
	}, nil
}

func (s *Stannp) BytesToPDF(data []byte) (*os.File, error) {
	tmpFile, err := ioutil.TempFile("", "example.*.pdf")
	if err != nil {
		return nil, err
	}

	if _, writeErr := tmpFile.Write(data); writeErr != nil {
		closeErr := tmpFile.Close()
		if closeErr != nil {
			return nil, closeErr
		}
		return nil, writeErr
	}

	return tmpFile, nil
}

func (s *Stannp) SendLetter(request *letter.SendReq) (*letter.SendRes, *util.APIError) {
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

	res, postErr := s.post(strings.NewReader(formData.Encode()), strings.Join([]string{s.baseUrl, letter.URL, CreateURL}, "/"))
	if postErr != nil {
		return nil, postErr
	}

	var letterRes letter.SendRes
	resErr := util.ResToType(res.StatusCode, res.Body, &letterRes)
	return &letterRes, resErr
}

func (s *Stannp) ValidateAddress(request *address.ValidateReq) (*address.ValidateRes, *util.APIError) {
	// Create URL values
	formData := url.Values{}
	formData.Set("company", request.Company)
	formData.Set("address1", request.Address1)
	formData.Set("address2", request.Address2)
	formData.Set("city", request.City)
	formData.Set("zipcode", request.Zipcode)
	formData.Set("country", request.Country)

	res, postErr := s.post(strings.NewReader(formData.Encode()), strings.Join([]string{s.baseUrl, address.URL, ValidateURL}, "/"))
	if postErr != nil {
		return nil, postErr
	}

	var addressRes address.ValidateRes
	resErr := util.ResToType(res.StatusCode, res.Body, &addressRes)
	return &addressRes, resErr
}
