package providers

import (
	"reflect"
	"testing"
)

const (
	githubTestSecret = "MyGithubTestSecret"
)

func TestNewGithubProvider(t *testing.T) {
	type args struct {
		secret string
	}
	tests := []struct {
		name    string
		args    args
		want    *GithubProvider
		wantErr bool
	}{
		{
			name: "TestNewGithubProviderWithCorrectSecret",
			args: args{
				secret: githubTestSecret,
			},
			want: &GithubProvider{
				secret: githubTestSecret,
			},
			wantErr: false,
		},
		{
			name: "TestNewGithubProviderWithEmptySecret",
			args: args{
				secret: "",
			},
			want: &GithubProvider{
				secret: "",
			},
			wantErr: false,
		},
		{
			name:    "TestNewGithubProviderWithNoSecret",
			args:    args{},
			want:    &GithubProvider{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewGithubProvider(tt.args.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGithubProvider() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGithubProvider() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGithubProvider_GetHeaderKeys(t *testing.T) {
	type fields struct {
		secret string
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "TestGetHeaderKeysWithCorrectValues",
			want: []string{XGitHubDelivery, XGitHubEvent, ContentTypeHeader},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &GithubProvider{
				secret: tt.fields.secret,
			}
			if got := p.GetHeaderKeys(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GithubProvider.GetHeaderKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

// A REAL PAYLOAD AND REAL DIGESTS, COMPUTED RATHER THAN PASTED. Every case
// below derives its signature from `githubTestSecret` and `testPayload` through
// the same helpers the provider uses, so a change to either helper fails these
// tests instead of quietly agreeing with itself.
var testPayload = []byte(`{"zen":"Non-blocking is better than blocking."}`)

func TestGithubProvider_Validate(t *testing.T) {
	validSHA1 := SignaturePrefix + HashPayload(githubTestSecret, testPayload)
	validSHA256 := Signature256Prefix + HashPayload256(githubTestSecret, testPayload)

	type fields struct {
		secret string
	}
	type args struct {
		hook Hook
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			// The header GitHub's own documentation tells receivers to use.
			name:   "TestValidateWithCorrectSHA256Signature",
			fields: fields{secret: githubTestSecret},
			args: args{hook: Hook{
				Headers: map[string]string{XHubSignature256: validSHA256},
				Payload: testPayload,
			}},
			want: true,
		},
		{
			// GitLab and older GitHub Enterprise send SHA-1 only. Dropping the
			// fallback would break working deployments for no security gain.
			name:   "TestValidateWithCorrectSHA1SignatureOnly",
			fields: fields{secret: githubTestSecret},
			args: args{hook: Hook{
				Headers: map[string]string{XHubSignature: validSHA1},
				Payload: testPayload,
			}},
			want: true,
		},
		{
			// THE ONE THAT MATTERS MOST. GitHub sends both headers, so a
			// receiver that fell back on a BAD SHA-256 would let an attacker
			// pick the weaker algorithm by sending a valid SHA-1 beside a
			// forged SHA-256. When the strong header is present it decides,
			// full stop.
			name:   "TestValidateRejectsBadSHA256EvenWithValidSHA1",
			fields: fields{secret: githubTestSecret},
			args: args{hook: Hook{
				Headers: map[string]string{
					XHubSignature:    validSHA1,
					XHubSignature256: Signature256Prefix + "00000000000000000000000000000000000000000000000000000000000000ff",
				},
				Payload: testPayload,
			}},
			want: false,
		},
		{
			name:   "TestValidateWithBothHeadersValid",
			fields: fields{secret: githubTestSecret},
			args: args{hook: Hook{
				Headers: map[string]string{
					XHubSignature:    validSHA1,
					XHubSignature256: validSHA256,
				},
				Payload: testPayload,
			}},
			want: true,
		},
		{
			// A signature over a different body must not validate; this is the
			// whole point of the check.
			name:   "TestValidateWithTamperedPayload",
			fields: fields{secret: githubTestSecret},
			args: args{hook: Hook{
				Headers: map[string]string{XHubSignature256: validSHA256},
				Payload: []byte(`{"zen":"Anything added dilutes everything else."}`),
			}},
			want: false,
		},
		{
			name:   "TestValidateWithWrongSecret",
			fields: fields{secret: "NotTheSecret"},
			args: args{hook: Hook{
				Headers: map[string]string{XHubSignature256: validSHA256},
				Payload: testPayload,
			}},
			want: false,
		},
		{
			// Length is checked before the prefix, so a truncated digest is
			// refused rather than compared.
			name:   "TestValidateWithTruncatedSHA256",
			fields: fields{secret: githubTestSecret},
			args: args{hook: Hook{
				Headers: map[string]string{XHubSignature256: Signature256Prefix + "abc123"},
				Payload: testPayload,
			}},
			want: false,
		},
		{
			// `sha1=` on the SHA-256 header is a malformed delivery, not a
			// reason to try the other algorithm.
			name:   "TestValidateWithWrongPrefixOnSHA256Header",
			fields: fields{secret: githubTestSecret},
			args: args{hook: Hook{
				Headers: map[string]string{
					XHubSignature256: SignaturePrefix + HashPayload(githubTestSecret, testPayload),
				},
				Payload: testPayload,
			}},
			want: false,
		},
		{
			name:   "TestValidateWithNoSignatureHeaders",
			fields: fields{secret: githubTestSecret},
			args:   args{hook: Hook{Headers: map[string]string{}, Payload: testPayload}},
			want:   false,
		},
		{
			name:   "TestValidateWithNilHeaders",
			fields: fields{secret: githubTestSecret},
			args:   args{hook: Hook{Headers: nil, Payload: testPayload}},
			want:   false,
		},
		{
			// An empty SHA-256 header must fall through to SHA-1 rather than
			// short-circuit, or a delivery carrying `X-Hub-Signature-256: ""`
			// would be rejected despite a valid SHA-1.
			name:   "TestValidateWithEmptySHA256FallsBackToSHA1",
			fields: fields{secret: githubTestSecret},
			args: args{hook: Hook{
				Headers: map[string]string{
					XHubSignature256: "",
					XHubSignature:    validSHA1,
				},
				Payload: testPayload,
			}},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &GithubProvider{
				secret: tt.fields.secret,
			}
			if got := p.Validate(tt.args.hook); got != tt.want {
				t.Errorf("GithubProvider.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
