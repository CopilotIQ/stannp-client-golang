package letter

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type MergeVariables map[string]string

type SendLetterRequest struct {
	ClearZone      bool             `json:"clearZone"`
	Duplex         bool             `json:"duplex"`
	MergeVariables MergeVariables   `json:"mergeVariables"`
	Recipient      RecipientDetails `json:"recipient"`
	Template       int              `json:"template"`
	Test           bool             `json:"test"`
}

type RecipientDetails struct {
	Address1  string `json:"address1"`
	Country   string `json:"country"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	State     string `json:"state"`
	Title     string `json:"title"`
	Town      string `json:"town"`
	Zipcode   string `json:"zipcode"`
}

type SendLetterResponseData struct {
	Pdf     string `json:"pdf"`
	Id      string `json:"id"`
	Created string `json:"created"`
	Format  string `json:"format"`
	Cost    string `json:"cost"`
	Status  string `json:"status"`
}

type SendLetterResponse struct {
	Success bool                   `json:"success"`
	Data    SendLetterResponseData `json:"data"`
}

func (a *API) SendLetter(request SendLetterRequest) (*SendLetterResponse, error) {
	formData := url.Values{}
	formData.Set("test", strconv.FormatBool(request.Test))
	formData.Set("template", strconv.Itoa(request.Template))
	formData.Set("clearzone", strconv.FormatBool(request.ClearZone))
	formData.Set("duplex", strconv.FormatBool(request.Duplex))
	formData.Set("post_unverified", strconv.FormatBool(a.postUnverified))
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

	req, err := http.NewRequest("POST", a.baseURL+"/letters/create", strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("apikey", a.apiKey)

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

	var response SendLetterResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
