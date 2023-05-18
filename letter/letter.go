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
	ClearZone      bool             `json:"clearZone"`
	Duplex         bool             `json:"duplex"`
	MergeVariables MergeVariables   `json:"mergeVariables"`
	Recipient      RecipientDetails `json:"recipient"`
	Template       int              `json:"template"`
}

type Data struct {
	Pdf     string `json:"pdf"`
	Id      string `json:"id"`
	Created string `json:"created"`
	Format  string `json:"format"`
	Cost    string `json:"cost"`
	Status  string `json:"status"`
}

type Response struct {
	Success bool `json:"success"`
	Data    Data `json:"data"`
}
