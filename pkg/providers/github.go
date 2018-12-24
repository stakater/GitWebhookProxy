package providers

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

const (
	GithubPushEvent Event = "push"
	GithubPullRequestEvent Event = "pull_request"
)

// Header constants
const (
	XHubSignature   = "X-Hub-Signature"
	XGitHubEvent    = "X-GitHub-Event"
	XGitHubDelivery = "X-GitHub-Delivery"
)

const (
	SignaturePrefix = "sha1="
	SignatureLength = 45
	GithubName = "github"
)

type GithubProvider struct {
	secret string
}

func NewGithubProvider(secret string) (*GithubProvider, error) {
	return &GithubProvider{
		secret: secret,
	}, nil
}

func (p *GithubProvider) GetHeaderKeys() []string {
	if (len(strings.TrimSpace(p.secret)) > 0) {
		return []string {
			XHubSignature,
			XGitHubDelivery,
			XGitHubEvent,
			ContentTypeHeader,
		}
	}

	return []string{
		XGitHubDelivery,
		XGitHubEvent,
		ContentTypeHeader,
	}
}

// TODO: Update implementation and tests
// Github Signature Validation:
func (p *GithubProvider) Validate(hook Hook) bool {

	githubSignature := hook.Headers[XHubSignature]
	if len(githubSignature) != SignatureLength ||
		!strings.HasPrefix(githubSignature, SignaturePrefix) {
		return false
	}

	return IsValidPayload(p.secret, githubSignature[len(SignaturePrefix):], hook.Payload)
}

func (p *GithubProvider) GetProviderName() string {
	return GithubName;
}

func (p *GithubProvider) GetCommitter(hook Hook) string {
	var pushPayloadData GithubPushPayload
	var pullRequestPayloadData GithubPullRequestPayload
	if err := json.Unmarshal(hook.Payload, &pushPayloadData); err != nil {
		log.Printf("Github payload unmarshaling failed for Push event: %v", err)
		log.Printf("Now trying to unmarshal for pull request event")
		if err = json.Unmarshal(hook.Payload, &pullRequestPayloadData); err != nil {
			log.Printf("Github payload unmarshaling failed for pull request event: %v", err)
			return "error"
		}
	}

	eventType := Event(hook.Headers[XGitHubEvent])
	switch eventType {
	case GithubPushEvent:
		return pushPayloadData.Sender.Login
	case GithubPullRequestEvent:
		return pullRequestPayloadData.Sender.Login
	}
	log.Printf("Not a pull request or push event: %v", eventType)
	return ""
}

// IsValidPayload checks if the github payload's hash fits with
// the hash computed by GitHub sent as a header
func IsValidPayload(secret, headerHash string, payload []byte) bool {
	hash := HashPayload(secret, payload)
	log.Printf("Calculated Hash: %s", hash)
	return hmac.Equal(
		[]byte(hash),
		[]byte(headerHash),
	)
}

// HashPayload computes the hash of payload's body according to the webhook's secret token
// see https://developer.github.com/webhooks/securing/#validating-payloads-from-github
// returning the hash as a hexadecimal string
func HashPayload(secret string, playloadBody []byte) string {
	hm := hmac.New(sha1.New, []byte(secret))
	hm.Write(playloadBody)
	sum := hm.Sum(nil)
	return fmt.Sprintf("%x", sum)
}
