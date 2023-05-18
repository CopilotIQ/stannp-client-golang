package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type APIError struct {
	Code    int    `json:"code"`
	Error   string `json:"error"`
	Success bool   `json:"success"`
}

func BuildError(code int, errorMessage string, success bool) *APIError {
	return &APIError{
		Code:    code,
		Error:   errorMessage,
		Success: success,
	}
}

func ResToType(code int, reader io.Reader, successType interface{}) *APIError {
	if code < http.StatusOK || (code < http.StatusBadRequest && code >= http.StatusMultipleChoices) {
		return BuildError(500, fmt.Sprintf("unexpected status code [%d]", code), false)
	}

	resBody, err := io.ReadAll(reader)
	if err != nil {
		return BuildError(500, fmt.Sprintf("error reading response body [%+v] with err [%+v]", string(resBody), err), false)
	}

	var jsonErr error
	var doReturnError bool
	var serverErr APIError
	if code >= http.StatusBadRequest {
		doReturnError = true
		jsonErr = json.Unmarshal(resBody, &serverErr)
		serverErr.Code = code
	} else {
		jsonErr = json.Unmarshal(resBody, &successType)
	}

	if jsonErr != nil {
		return BuildError(500, fmt.Sprintf("error unmarshalling res [%+v]", string(resBody)), false)
	}

	if doReturnError {
		return &serverErr
	}

	return nil
}
