package letter

import (
	"encoding/json"
	"io"
)

const URL = "letters"

type Data struct {
	Cost    string      `json:"cost"`
	Created string      `json:"created"`
	Format  string      `json:"format"`
	ID      json.Number `json:"id"`
	PDFURL  string      `json:"pdf"`
	Status  string      `json:"status"`
}

type PDFRes struct {
	Contents io.ReadCloser
	Name     string
}

type RecipientDetails struct {
	Address1  string `json:"address1"`
	Address2  string `json:"address2"`
	Country   string `json:"country"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	State     string `json:"state"`
	Title     string `json:"title"`
	Town      string `json:"town"`
	Zipcode   string `json:"zipcode"`
}

type MergeVariables map[string]string

type SendReq struct {
	IdempotenceyKey string           `json:"idempotenceyKey"`
	MergeVariables  MergeVariables   `json:"mergeVariables"`
	Recipient       RecipientDetails `json:"recipient"`
	Template        string           `json:"template"`
}

type SendRes struct {
	Data    Data `json:"data"`
	Success bool `json:"success"`
}
