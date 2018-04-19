package providers

import (
	"net/http"
)

type GithubProvider struct {
	secret string
}

func NewGithubProvider(secret string) (*GithubProvider, error) {
	// if secret == nil {
	// 	return nil, fmt.Errorf("blah ")
	// }

	return &GithubProvider{
		secret: secret,
	}, nil
}

func (p *GithubProvider) GetHook(req *http.Request) (Hook, error) {

	return Hook{Id: "lol1"}, nil
}

func (p *GithubProvider) validate() bool {
	return false
}
