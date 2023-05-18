package stannp

const BaseURL = "https://us.stannp.com/api/v1"
const CreateURL = "create"

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
	}

	for _, option := range options {
		option(api)
	}
	return api
}

func (s *Stannp) APIKey() string {
	return s.apiKey
}

func (s *Stannp) PostUnverified() bool {
	return s.postUnverified
}

func (s *Stannp) IsTest() bool {
	return s.test
}
