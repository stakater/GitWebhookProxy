package providers

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"strings"
)

// Header constants
const (
	XGitHubSignature = "X-GitHub-Signature"
	XGitHubEvent     = "X-GitHub-Event"
	XGitHubDelivery  = "X-GitHub-Delivery"
)

const (
	SignaturePrefix = "sha1="
	SignatureLength = 45
)

type GithubProvider struct {
	secret string
}

func NewGithubProvider(secret string) (*GithubProvider, error) {
	if len(strings.TrimSpace(secret)) == 0 {
		return nil, errors.New("Cannot create github provider with empty secret")
	}

	return &GithubProvider{
		secret: secret,
	}, nil
}

func (p *GithubProvider) GetHeaderKeys() []string {
	return []string{
		XGitHubDelivery,
		XGitHubEvent,
		XGitHubSignature,
	}
}

// TODO: Update implementation and tests

// Github Signature Validation:
// https://developer.github.com/webhooks/securing/#validating-payloads-from-github
func (p *GithubProvider) Validate(hook Hook) bool {

	githubSignature := hook.Headers[XGitHubSignature]
	if len(githubSignature) != SignatureLength ||
		!strings.HasPrefix(githubSignature, SignaturePrefix) {
		return false
	}

	// decodedSignature := make([]byte, 20)
	// hex.Decode(decodedSignature, []byte(githubSignature[len(SignaturePrefix):]))

	decodedSignature, err := hex.DecodeString(githubSignature[len(SignaturePrefix):])
	if err != nil {
		panic("error decoding")
	}

	//TODO: Return this
	hmac.Equal(signBody([]byte(p.secret), hook.Payload), decodedSignature)
	return true
}

func signBody(secret []byte, body []byte) []byte {
	computed := hmac.New(sha1.New, secret)
	computed.Write(body)
	return computed.Sum(nil)
}
