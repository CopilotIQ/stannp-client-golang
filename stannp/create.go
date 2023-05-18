package stannp

import (
	"copilotiq/stannp-client-golang/letter"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func (s *Stannp) SendLetter(request letter.Request) (*letter.Response, error) {
	formData := url.Values{}
	formData.Set("test", strconv.FormatBool(request.Test))
	formData.Set("template", strconv.Itoa(request.Template))
	formData.Set("clearzone", strconv.FormatBool(request.ClearZone))
	formData.Set("duplex", strconv.FormatBool(request.Duplex))
	formData.Set("post_unverified", strconv.FormatBool(s.))
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

	reqUrl := strings.Join([]string{s.baseUrl, letter.URL, CreateURL}, "/")

	req, err := http.NewRequest("POST", reqUrl, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("apikey", s.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response letter.Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
