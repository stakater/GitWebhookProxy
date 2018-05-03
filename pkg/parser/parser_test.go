package parser

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stakater/GitWebhookProxy/pkg/providers"
)

const (
	parserGitlabTestSecret = "testSecret"
	parserGitlabTestEvent  = "testEvent"
	parserGitlabTestBody   = "testBody"
)

func createGitlabRequest(method string, path string, tokenHeader string,
	eventHeader string, body string) *http.Request {
	req := httptest.NewRequest(method, path, bytes.NewReader([]byte(body)))
	req.Header.Add(providers.XGitlabToken, tokenHeader)
	req.Header.Add(providers.XGitlabEvent, eventHeader)
	req.Header.Add(providers.ContentTypeHeader, providers.DefaultContentTypeHeaderValue)
	return req
}

func createRequestWithWrongHeaders(method string, path string, tokenHeader string,
	eventHeader string, body string) *http.Request {
	req := httptest.NewRequest(method, path, bytes.NewReader([]byte(body)))
	req.Header.Add("X-Wrong-Token", tokenHeader)
	req.Header.Add("X-Wrong-Event", eventHeader)
	return req
}

func createGitlabProvider(secret string) providers.Provider {
	provider, _ := providers.NewGitlabProvider(secret)
	return provider
}

func createGitlabHook(tokenHeader string, tokenEvent string, body string, method string) *providers.Hook {
	return &providers.Hook{
		Headers: map[string]string{
			providers.XGitlabToken:      tokenHeader,
			providers.XGitlabEvent:      tokenEvent,
			providers.ContentTypeHeader: providers.DefaultContentTypeHeaderValue,
		},
		Payload:       []byte(body),
		RequestMethod: method,
	}
}

func TestParse(t *testing.T) {
	type args struct {
		req      *http.Request
		provider providers.Provider
	}

	tests := []struct {
		name    string
		args    args
		want    *providers.Hook
		wantErr bool
	}{
		{
			name: "TestParseWithCorrectRequestValues",
			args: args{
				req: createGitlabRequest("post", "/dummy", parserGitlabTestSecret,
					parserGitlabTestEvent, parserGitlabTestBody),
				provider: createGitlabProvider(parserGitlabTestSecret),
			},
			want: createGitlabHook(parserGitlabTestSecret, parserGitlabTestEvent, parserGitlabTestBody, "post"),
		},
		{
			name: "TestParseWithEmptyTokenHeaderValue",
			args: args{
				req: createGitlabRequest("post", "/dummy", "",
					parserGitlabTestEvent, parserGitlabTestBody),
				provider: createGitlabProvider(parserGitlabTestSecret),
			},
			wantErr: true,
		},
		{
			name: "TestParseWithNoEventHeaderValue",
			args: args{
				req: createGitlabRequest("post", "/dummy", parserGitlabTestSecret,
					"", parserGitlabTestBody),
				provider: createGitlabProvider(parserGitlabTestSecret),
			},
			wantErr: true,
		},
		{
			name: "TestParseWithNoBody",
			args: args{
				req: createGitlabRequest("post", "/dummy", parserGitlabTestSecret,
					parserGitlabTestEvent, ""),
				provider: createGitlabProvider(parserGitlabTestSecret),
			},
			want: createGitlabHook(parserGitlabTestSecret, parserGitlabTestEvent, "", "post"),
		},
		{
			name: "TestParseWithNoHeaders",
			args: args{
				req:      httptest.NewRequest("post", "/dummy", bytes.NewReader([]byte(parserGitlabTestBody))),
				provider: createGitlabProvider(parserGitlabTestSecret),
			},
			wantErr: true,
		},
		{
			name: "TestParseWithWrongHeaderKeys",
			args: args{
				req: createRequestWithWrongHeaders("post", "/dummy", parserGitlabTestSecret,
					parserGitlabTestEvent, parserGitlabTestBody),
				provider: createGitlabProvider(parserGitlabTestSecret),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.req, tt.args.provider)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
