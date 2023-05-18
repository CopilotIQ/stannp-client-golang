package stannp

const BaseUrl = "https://us.stannp.com/api/v1"

type ErrorDetails struct {
	Error   string `json:"error"`
	Success bool   `json:"success"`
}

type Stannp struct {
	apiKey         string
	baseUrl        string
	postUnverified bool
	test           bool
}

type APIOption func(*Stannp)

func WithTest(test bool) APIOption {
	return func(a *Stannp) {
		a.test = test
	}
}

func WithAPIKey(apiKey string) APIOption {
	return func(a *Stannp) {
		a.apiKey = apiKey
	}
}

func WithPostUnverified(postUnverified bool) APIOption {
	return func(a *Stannp) {
		a.postUnverified = postUnverified
	}
}

func New(options ...APIOption) *Stannp {
	api := &Stannp{
		apiKey:         "test123456",
		baseUrl:        BaseUrl,
		postUnverified: false,
		test:           true,
	}

	for _, option := range options {
		option(api)
	}
	return api
}

func (api *Stannp) APIKey() string {
	return api.apiKey
}

func (api *Stannp) IsTest() bool {
	return api.test
}
