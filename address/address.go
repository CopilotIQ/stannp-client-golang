package address

const URL = "addresses"

type ValidateReq struct {
	Address1 string `json:"address1"`
	Address2 string `json:"address2"`
	City     string `json:"city"`
	Company  string `json:"company"`
	Country  string `json:"country"`
	Zipcode  string `json:"zipcode"`
}

type Data struct {
	IsValid bool `json:"is_valid"`
}

type ValidateRes struct {
	Data    Data `json:"data"`
	Success bool `json:"success"`
}
