package proxy

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stakater/JenkinsProxy/pkg/providers"
)

const (
	proxyGitlabTestSecret = "testSecret"
	proxyGitlabTestEvent  = "testEvent"
	proxyGitlabTestBody   = "testBody"
)

func TestProxy_isPathAllowed(t *testing.T) {
	type fields struct {
		provider     string
		upstreamURL  string
		allowedPaths []string
		secret       string
	}
	type args struct {
		path string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "isPathAllowedWithValidMultipleAllowedPaths",
			fields: fields{
				provider:     providers.GithubProviderKind,
				upstreamURL:  "https://dummyurl.com",
				allowedPaths: []string{"/path1", "/path2"},
				secret:       "secret",
			},
			args: args{
				path: "/path2",
			},
			want: true,
		},
		{
			name: "isPathAllowedWithValidOneAllowedPaths",
			fields: fields{
				provider:     providers.GithubProviderKind,
				upstreamURL:  "https://dummyurl.com",
				allowedPaths: []string{"/path1"},
				secret:       "secret",
			},
			args: args{
				path: "/path1",
			},
			want: true,
		},
		{
			name: "isPathAllowedWithInvalidPath",
			fields: fields{
				provider:     providers.GithubProviderKind,
				upstreamURL:  "https://dummyurl.com",
				allowedPaths: []string{"/path1", "/path2"},
				secret:       "secret",
			},
			args: args{
				path: "/path3",
			},
			want: false,
		},
		{
			name: "isPathAllowedWithEmtpyPathArg",
			fields: fields{
				provider:     providers.GithubProviderKind,
				upstreamURL:  "https://dummyurl.com",
				allowedPaths: []string{"/path1", "/path2"},
				secret:       "secret",
			},
			args: args{
				path: "",
			},
			want: false,
		},
		{
			name: "isPathAllowedWithAllPathsAllowedAndEmptyPathArg",
			fields: fields{
				provider:     providers.GithubProviderKind,
				upstreamURL:  "https://dummyurl.com",
				allowedPaths: []string{},
				secret:       "secret",
			},
			args: args{
				path: "",
			},
			want: true,
		},
		{
			name: "isPathAllowedWithAllPathsAllowedAndRootEmptyPathArg",
			fields: fields{
				provider:     providers.GithubProviderKind,
				upstreamURL:  "https://dummyurl.com",
				allowedPaths: []string{},
				secret:       "secret",
			},
			args: args{
				path: "/",
			},
			want: true,
		},
		{
			name: "isPathAllowedWithAllPathsAllowedAndNonEmptyPathArg",
			fields: fields{
				provider:     providers.GithubProviderKind,
				upstreamURL:  "https://dummyurl.com",
				allowedPaths: []string{},
				secret:       "secret",
			},
			args: args{
				path: "/path1",
			},
			want: true,
		},
		{
			name: "isPathAllowedWithSomePathsAllowedAndRootPathArg",
			fields: fields{
				provider:     providers.GithubProviderKind,
				upstreamURL:  "https://dummyurl.com",
				allowedPaths: []string{"/path1", "/path2"},
				secret:       "secret",
			},
			args: args{
				path: "/",
			},
			want: false,
		},
		{
			name: "isPathAllowedWithSomePathsAllowedAndSubPathArg",
			fields: fields{
				provider:     providers.GithubProviderKind,
				upstreamURL:  "https://dummyurl.com",
				allowedPaths: []string{"/path1", "/path2"},
				secret:       "secret",
			},
			args: args{
				path: "/path2/path3",
			},
			want: false,
		},
		{
			name: "isPathAllowedWithSubPathsAllowedAndSubPathArg",
			fields: fields{
				provider:     providers.GithubProviderKind,
				upstreamURL:  "https://dummyurl.com",
				allowedPaths: []string{"/path1", "/path2/path3"},
				secret:       "secret",
			},
			args: args{
				path: "/path2/path3",
			},
			want: true,
		},
		{
			name: "isPathAllowedWithSubPathsAllowedAndPathArg",
			fields: fields{
				provider:     providers.GithubProviderKind,
				upstreamURL:  "https://dummyurl.com",
				allowedPaths: []string{"/path1", "/path2/path3"},
				secret:       "secret",
			},
			args: args{
				path: "/path2",
			},
			want: false,
		},
		{
			name: "isPathAllowedWithAllowedPathTrailingSlashAndNotInArg",
			fields: fields{
				provider:     providers.GithubProviderKind,
				upstreamURL:  "https://dummyurl.com",
				allowedPaths: []string{"/path1", "/path2/"},
				secret:       "secret",
			},
			args: args{
				path: "/path2",
			},
			want: true,
		},
		{
			name: "isPathAllowedWithSimpleAllowedPathAndTrailingSlashInArg",
			fields: fields{
				provider:     providers.GithubProviderKind,
				upstreamURL:  "https://dummyurl.com",
				allowedPaths: []string{"/path1", "/path2"},
				secret:       "secret",
			},
			args: args{
				path: "/path2/",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Proxy{
				provider:     tt.fields.provider,
				upstreamURL:  tt.fields.upstreamURL,
				allowedPaths: tt.fields.allowedPaths,
				secret:       tt.fields.secret,
			}
			if got := p.isPathAllowed(tt.args.path); got != tt.want {
				t.Errorf("Proxy.isPathAllowed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func createGitlabHook(tokenHeader string, tokenEvent string, body string) *providers.Hook {
	return &providers.Hook{
		Headers: map[string]string{
			providers.XGitlabToken: tokenHeader,
			providers.XGitlabEvent: tokenEvent,
		},
		Payload: []byte(body),
	}
}

func TestProxy_redirect(t *testing.T) {
	type fields struct {
		provider     string
		upstreamURL  string
		allowedPaths []string
		secret       string
	}
	type args struct {
		hook *providers.Hook
		path string
	}
	tests := []struct {
		name               string
		fields             fields
		args               args
		wantStatusCode     int
		wantRedirectedHost string // Only Host not complete URL
		wantErr            bool
	}{
		{
			name: "TestRedirectWithValidValues",
			fields: fields{
				provider:     "gitlab",
				upstreamURL:  "https://httpbin.org",
				allowedPaths: []string{},
				secret:       "dummy",
			},
			args: args{
				path: "/post",
				hook: createGitlabHook(proxyGitlabTestSecret, proxyGitlabTestEvent, proxyGitlabTestBody),
			},
			wantStatusCode:     http.StatusOK,
			wantRedirectedHost: "httpbin.org",
		},
		//TODO: With get url and different values
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Proxy{
				provider:     tt.fields.provider,
				upstreamURL:  tt.fields.upstreamURL,
				allowedPaths: tt.fields.allowedPaths,
				secret:       tt.fields.secret,
			}
			gotResp, gotErrors := p.redirect(tt.args.hook, tt.args.path)

			if (gotErrors != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", gotErrors, tt.wantErr)
				return
			}

			if gotResp.StatusCode != tt.wantStatusCode {
				t.Errorf("Proxy.redirect() got StatusCode in response= %v, want %v",
					gotResp.StatusCode, tt.wantStatusCode)
			}

			if gotResp.Request.Host != tt.wantRedirectedHost {
				t.Errorf("Proxy.redirect() got Redirected Host in response= %v, want Redirected Host= %v",
					gotResp.Request.Host, tt.wantRedirectedHost)
			}
		})
	}
}

func TestProxy_proxyRequest(t *testing.T) {
	type fields struct {
		provider     string
		upstreamURL  string
		allowedPaths []string
		secret       string
	}
	type args struct {
		httpMethod string
		path       string
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantStatusCode int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Proxy{
				provider:     tt.fields.provider,
				upstreamURL:  tt.fields.upstreamURL,
				allowedPaths: tt.fields.allowedPaths,
				secret:       tt.fields.secret,
			}
			router := httprouter.New()
			router.POST("/*path", p.proxyRequest)

			req, err := http.NewRequest(tt.args.httpMethod, tt.args.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.wantStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.wantStatusCode)
			}

		})
	}
}

func TestProxy_health(t *testing.T) {
	type fields struct {
		provider     string
		upstreamURL  string
		allowedPaths []string
		secret       string
	}
	type args struct {
		httpMethod string
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantStatusCode int
	}{
		{
			name: "TestHealthCheckGet",
			args: args{
				httpMethod: "GET",
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "TestHealthCheckPost",
			args: args{
				httpMethod: "POST",
			},
			wantStatusCode: http.StatusMethodNotAllowed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Proxy{
				provider:     tt.fields.provider,
				upstreamURL:  tt.fields.upstreamURL,
				allowedPaths: tt.fields.allowedPaths,
				secret:       tt.fields.secret,
			}
			router := httprouter.New()
			router.GET("/health", p.health)

			req, err := http.NewRequest(tt.args.httpMethod, "/health", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.wantStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.wantStatusCode)
			}
		})
	}
}

func TestNewProxy(t *testing.T) {
	type args struct {
		listenAddress string
		upstreamURL   string
		allowedPaths  []string
		provider      string
		secret        string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewProxy(tt.args.listenAddress, tt.args.upstreamURL, tt.args.allowedPaths, tt.args.provider, tt.args.secret)
		})
	}
}
