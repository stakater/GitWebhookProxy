package providers

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
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

// func (p *GithubProvider) GetTokenHeaderKey() string {
// 	return XGitHubSignature
// }

func (p *GithubProvider) Validate(hook Hook) bool {
	// githubSignature := hook.Headers[XGitHubSignature]
	// if len(githubSignature) != SignatureLength ||
	// 	!strings.HasPrefix(githubSignature, SignaturePrefix) {
	// 	return false
	// }
	// // decodedSignature := make([]byte, 20)

	// // hex.Decode(decodedSignature, []byte(githubSignature[len(SignaturePrefix):]))

	// decodedSignature, err := hex.DecodeString(githubSignature[len(SignaturePrefix):])
	// if err != nil {
	// 	panic("error decoding")
	// }
	// fmt.Printf("decodedSignature: %s\n", decodedSignature)
	// fmt.Printf("githubSignature: %s\n", githubSignature)
	// fmt.Printf("githubSignatureTrim: %s\n", githubSignature[len(SignaturePrefix):])
	// fmt.Printf("p.secret: %s\n", p.secret)
	// return hmac.Equal(signBody([]byte(p.secret), hook.Payload), decodedSignature)
	return true
}

func signBody(secret []byte, body []byte) []byte {
	computed := hmac.New(sha1.New, secret)
	computed.Write(body)
	fmt.Printf("Computed to string: %s ", hex.EncodeToString(computed.Sum(nil)))
	return computed.Sum(nil)
}
