package letter

const URL = "letters"

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

type MergeVariables map[string]string

type Request struct {
	MergeVariables MergeVariables   `json:"mergeVariables"`
	Recipient      RecipientDetails `json:"recipient"`
	Template       string           `json:"template"`
}

type Data struct {
	Cost    string `json:"cost"`
	Created string `json:"created"`
	Format  string `json:"format"`
	Id      string `json:"id"`
	Pdf     string `json:"pdf"`
	Status  string `json:"status"`
}

type Response struct {
	Data    Data `json:"data"`
	Success bool `json:"success"`
}
