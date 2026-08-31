package providers

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

const (
	GithubPushEvent         Event = "push"
	GithubPullRequestEvent  Event = "pull_request"
	GithubIssueCommentEvent Event = "issue_comment"
)

// Header constants
const (
	XHubSignature = "X-Hub-Signature"
	// XHubSignature256 carries the HMAC-SHA256 digest. GitHub sends it on every
	// delivery alongside the SHA-1 header, and its own documentation says to
	// use it: "If you are using the X-Hub-Signature-256 header, you should use
	// the HMAC-SHA256 algorithm."
	XHubSignature256 = "X-Hub-Signature-256"
	XGitHubEvent     = "X-GitHub-Event"
	XGitHubDelivery  = "X-GitHub-Delivery"
)

const (
	SignaturePrefix = "sha1="
	// SignatureLength is len("sha1=") + 40 hex characters.
	SignatureLength = 45

	Signature256Prefix = "sha256="
	// Signature256Length is len("sha256=") + 64 hex characters.
	Signature256Length = 71

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

// GetOptionalHeaderKeys forwards whichever signature headers are present.
//
// BOTH SIGNATURE HEADERS ARE OPTIONAL AT THE PARSER, and enforcement moves to
// `Validate`, which is the only place that can express the real rule: at least
// ONE of them, and if the strong one is there it decides.
//
// Requiring either at the parser is wrong in a different direction each way.
// Requiring SHA-256 rejects GitLab and older GitHub Enterprise, which send
// SHA-1 alone — the very deliveries the fallback exists for. Requiring SHA-1
// rejects a sender that has dropped the legacy header, which is where GitHub is
// heading. Neither is a rule about the request; both are a rule about
// signatures, so both belong with the signature check.
//
// Nothing is weakened: a request with no signature at all still fails, one
// step later, in `Validate`.
func (p *GithubProvider) GetOptionalHeaderKeys() []string {
	if len(strings.TrimSpace(p.secret)) > 0 {
		return []string{XHubSignature, XHubSignature256}
	}
	return nil
}

func (p *GithubProvider) GetHeaderKeys() []string {
	if len(strings.TrimSpace(p.secret)) > 0 {
		return []string{
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

// Validate checks the delivery's signature against the configured secret.
//
// SHA-256 IS PREFERRED AND SHA-1 IS THE FALLBACK. GitHub sends both headers on
// every delivery, and its documentation is explicit that X-Hub-Signature-256 is
// the one to use — SHA-1 is retained only for compatibility with receivers that
// predate it. A proxy that validates only the SHA-1 header is doing weaker
// checking than the payload already permits, on the one process that is exposed
// to the internet by design.
//
// The fallback is deliberate rather than reluctant: GitLab and older GitHub
// Enterprise installations still send SHA-1 only, and dropping it would break
// working deployments for no security gain when SHA-256 is absent. When both
// headers are present the SHA-256 one decides, so an attacker cannot choose the
// weaker algorithm by omitting a header — omitting X-Hub-Signature-256 on a
// GitHub delivery that would carry it produces a SHA-1 check the attacker still
// cannot satisfy without the secret.
func (p *GithubProvider) Validate(hook Hook) bool {
	if signature := hook.Headers[XHubSignature256]; signature != "" {
		if len(signature) != Signature256Length ||
			!strings.HasPrefix(signature, Signature256Prefix) {
			return false
		}
		return IsValidPayload256(p.secret, signature[len(Signature256Prefix):], hook.Payload)
	}

	githubSignature := hook.Headers[XHubSignature]
	if len(githubSignature) != SignatureLength ||
		!strings.HasPrefix(githubSignature, SignaturePrefix) {
		return false
	}

	return IsValidPayload(p.secret, githubSignature[len(SignaturePrefix):], hook.Payload)
}

func (p *GithubProvider) GetProviderName() string {
	return GithubName
}

func (p *GithubProvider) GetCommitter(hook Hook) string {
	eventType := Event(hook.Headers[XGitHubEvent])
	var pushPayloadData GithubPushPayload
	var pullRequestPayloadData GithubPullRequestPayload
	var issueCommentPayloadData GithubIssueCommentPayload

	log.Printf("Received event type: %v", eventType)
	switch eventType {
	case GithubPushEvent:
		if err := json.Unmarshal(hook.Payload, &pushPayloadData); err != nil {
			log.Printf("Github payload unmarshaling failed for Push event: %v", err)
			return ""
		}
		return pushPayloadData.Sender.Login
	case GithubPullRequestEvent:
		if err := json.Unmarshal(hook.Payload, &pullRequestPayloadData); err != nil {
			log.Printf("Github payload unmarshaling failed for Pull Request event: %v", err)
			return ""
		}
		return pullRequestPayloadData.Sender.Login
	case GithubIssueCommentEvent:
		if err := json.Unmarshal(hook.Payload, &issueCommentPayloadData); err != nil {
			log.Printf("Github payload unmarshaling failed for issue comment event: %v", err)
			return ""
		}
		return issueCommentPayloadData.Comment.User.Login
	}

	log.Printf("Event type is not supported: %v", eventType)
	return ""
}

// IsValidPayload checks if the github payload's SHA-1 hash fits with
// the hash computed by GitHub sent as a header.
//
// The computed hash is NOT logged. It used to be, and that was a real leak: the
// digest is derived from the shared secret over a body an attacker can choose,
// so publishing it to logs hands out an oracle for free. There is nothing a
// reader can do with the value that they could not do with the header they
// already sent, and the failure it helps debug — a mismatched secret — is
// equally visible from the rejection itself.
func IsValidPayload(secret, headerHash string, payload []byte) bool {
	hash := HashPayload(secret, payload)
	return hmac.Equal(
		[]byte(hash),
		[]byte(headerHash),
	)
}

// IsValidPayload256 checks the payload against the HMAC-SHA256 digest GitHub
// sends in X-Hub-Signature-256.
func IsValidPayload256(secret, headerHash string, payload []byte) bool {
	hash := HashPayload256(secret, payload)
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

// HashPayload256 computes the HMAC-SHA256 of the payload body under the
// webhook's secret, returned as a hexadecimal string — the form GitHub sends
// after the `sha256=` prefix.
func HashPayload256(secret string, payloadBody []byte) string {
	hm := hmac.New(sha256.New, []byte(secret))
	hm.Write(payloadBody)
	sum := hm.Sum(nil)
	return fmt.Sprintf("%x", sum)
}
