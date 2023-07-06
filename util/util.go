package util

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type APIError struct {
	Code         int    `json:"code"`
	ErrorMessage string `json:"error"`
	Success      bool   `json:"success"`
}

func (apiError *APIError) Error() string {
	return fmt.Sprintf("Code [%d] ErrorMessage [%s] Success [%t]", apiError.Code, apiError.ErrorMessage, apiError.Success)
}

func (apiError *APIError) String() string {
	return fmt.Sprintf("Stannp Client APIError: Code [%d] Success [%t] ErrorMessage [%s]", apiError.Code, apiError.Success, apiError.ErrorMessage)
}

func BuildError(code int, errorMessage string) *APIError {
	return &APIError{
		Code:         code,
		ErrorMessage: errorMessage,
		Success:      false,
	}
}

func RandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, n)

	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}

	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}

	return string(bytes)
}

func ResToType(code int, reader io.Reader, successType interface{}) *APIError {
	if code < http.StatusOK || (code < http.StatusBadRequest && code >= http.StatusMultipleChoices) {
		return BuildError(500, fmt.Sprintf("unexpected status code [%d]", code))
	}

	resBody, err := io.ReadAll(reader)
	if err != nil {
		return BuildError(500, fmt.Sprintf("error reading response body [%+v] with err [%+v]", string(resBody), err))
	}

	var jsonErr error
	var serverErr *APIError
	if code >= http.StatusBadRequest {
		jsonErr = json.Unmarshal(resBody, serverErr)
		serverErr.Code = code
	} else {
		jsonErr = json.Unmarshal(resBody, &successType)
	}

	if jsonErr != nil {
		return BuildError(500, fmt.Sprintf("error unmarshalling res [%+v]", string(resBody)))
	}

	return serverErr
}
