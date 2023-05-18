package stannp

import (
	"copilotiq/stannp-client-golang/letter"
	"copilotiq/stannp-client-golang/util"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const BaseURL = "https://us.stannp.com/api/v1"
const CreateURL = "create"

type Stannp struct {
	apiKey         string
	baseUrl        string
	client         *http.Client
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

func New(options ...APIOption) *Stannp {
	api := &Stannp{
		apiKey:         "test123456",
		baseUrl:        BaseURL,
		postUnverified: false,
		test:           true,
		client:         http.DefaultClient,
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
		return "", util.BuildError(500, fmt.Sprintf("error parsing inputURL [%s]", inputURL), false)
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
		return nil, util.BuildError(500, fmt.Sprintf("error generating POST req [%+v] for req [%+v]", err, req), false)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := s.client.Do(req)
	if err != nil {
		return nil, util.BuildError(500, fmt.Sprintf("error sending req [%+v]", req), false)
	}

	return res, nil
}

func (s *Stannp) SendLetter(request letter.Request) (*letter.Response, *util.APIError) {
	formData := url.Values{}
	formData.Set("test", strconv.FormatBool(s.test))
	formData.Set("template", strconv.Itoa(request.Template))
	formData.Set("clearzone", strconv.FormatBool(request.ClearZone))
	formData.Set("duplex", strconv.FormatBool(request.Duplex))
	formData.Set("post_unverified", strconv.FormatBool(s.postUnverified))
	formData.Set("recipient[title]", request.Recipient.Title)
	formData.Set("recipient[firstname]", request.Recipient.Firstname)
	formData.Set("recipient[lastname]", request.Recipient.Lastname)
	formData.Set("recipient[address1]", request.Recipient.Address1)
	formData.Set("recipient[town]", request.Recipient.Town)
	formData.Set("recipient[zipcode]", request.Recipient.Zipcode)
	formData.Set("recipient[state]", request.Recipient.State)
	formData.Set("recipient[country]", request.Recipient.Country)

	// set custom merge variables in the formData
	for key, value := range request.MergeVariables {
		formData.Set("recipient["+key+"]", value)
	}

	res, postErr := s.post(strings.NewReader(formData.Encode()), strings.Join([]string{s.baseUrl, letter.URL, CreateURL}, "/"))
	if postErr != nil {
		return nil, postErr
	}

	var letterRes letter.Response
	resErr := util.ResToType(res.StatusCode, res.Body, &letterRes)
	if resErr != nil {
		return nil, resErr
	}

	return &letterRes, nil
}
