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

func TestGithubProvider_Validate(t *testing.T) {
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
		// {
		// 	name: "TestValidateWithEmptySignatureValue",
		// 	fields: fields{
		// 		secret: githubTestSecret,
		// 	},
		// 	args: args{
		// 		hook: Hook{
		// 			Headers: map[string]string{
		// 				XHubSignature: "",
		// 			},
		// 		},
		// 	},
		// 	want: false,
		// },
		// {
		// 	name: "TestValidateWithEmptyHeaders",
		// 	fields: fields{
		// 		secret: githubTestSecret,
		// 	},
		// 	args: args{
		// 		hook: Hook{
		// 			Headers: map[string]string{},
		// 		},
		// 	},
		// 	want: false,
		// },
		// {
		// 	name: "TestValidateWithWrongSignatureValue",
		// 	fields: fields{
		// 		secret: githubTestSecret,
		// 	},
		// 	args: args{
		// 		hook: Hook{
		// 			Headers: map[string]string{
		// 				XHubSignature: "IncorrectSecret",
		// 			},
		// 			Payload: nil,
		// 		},
		// 	},
		// 	want: false,
		// },
		// {
		// 	name: "TestValidateWithCorrectTokenValue",
		// 	fields: fields{
		// 		secret: githubTestSecret,
		// 	},
		// 	args: args{
		// 		hook: Hook{
		// 			Headers: map[string]string{
		// 				XHubSignature: githubTestSecret,
		// 			},
		// 			Payload: nil,
		// 		},
		// 	},
		// 	want: true,
		// },
		// {
		// 	name: "TestValidateWithWrongHeaderKey",
		// 	fields: fields{
		// 		secret: githubTestSecret,
		// 	},
		// 	args: args{
		// 		hook: Hook{
		// 			Headers: map[string]string{
		// 				"X-Wrong-Signature": "IncorrectSecret",
		// 			},
		// 			Payload: nil,
		// 		},
		// 	},
		// 	want: false,
		// },
		// {
		// 	name: "TestValidateWithNilHeaders",
		// 	fields: fields{
		// 		secret: githubTestSecret,
		// 	},
		// 	args: args{
		// 		hook: Hook{
		// 			Headers: nil,
		// 			Payload: nil,
		// 		},
		// 	},
		// 	want: false,
		// },
		// {
		// 	name: "TestValidateWithNoHookArg",
		// 	fields: fields{
		// 		secret: githubTestSecret,
		// 	},
		// 	args: args{},
		// 	want: false,
		// },
		// {
		// 	name: "TestValidateWithWrongSecretInProxy",
		// 	fields: fields{
		// 		secret: "WrongSecret",
		// 	},
		// 	args: args{
		// 		hook: Hook{
		// 			Headers: map[string]string{
		// 				XHubSignature: githubTestSecret,
		// 			},
		// 			Payload: nil,
		// 		},
		// 	},
		// 	want: false,
		// },
		// {
		// 	name: "TestValidateWithEmptySecretInProxy",
		// 	fields: fields{
		// 		secret: "",
		// 	},
		// 	args: args{
		// 		hook: Hook{
		// 			Headers: map[string]string{
		// 				XHubSignature: githubTestSecret,
		// 			},
		// 			Payload: nil,
		// 		},
		// 	},
		// 	want: false,
		// },
		// {
		// 	name:   "TestValidateWithNoSecretInProxy",
		// 	fields: fields{},
		// 	args: args{
		// 		hook: Hook{
		// 			Headers: map[string]string{
		// 				XHubSignature: githubTestSecret,
		// 			},
		// 			Payload: nil,
		// 		},
		// 	},
		// 	want: false,
		// },
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
