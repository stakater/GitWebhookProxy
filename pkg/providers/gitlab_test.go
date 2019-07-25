package providers

import (
	"reflect"
	"testing"
)

const (
	gitlabTestSecret = "myGitLabTestSecret"
)

func TestNewGitlabProvider(t *testing.T) {
	type args struct {
		secret string
	}
	tests := []struct {
		name    string
		args    args
		want    *GitlabProvider
		wantErr bool
	}{
		{
			name: "TestNewGitlabProviderWithCorrectSecret",
			args: args{
				secret: gitlabTestSecret,
			},
			want: &GitlabProvider{
				secret: gitlabTestSecret,
			},
			wantErr: false,
		},
		{
			name: "TestNewGitlabProviderWithEmptySecret",
			args: args{
				secret: "",
			},
			want: &GitlabProvider{
				secret: "",
			},
			wantErr: false,
		},
		{
			name:    "TestNewGitlabProviderWithNoSecret",
			args:    args{},
			want:    &GitlabProvider{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewGitlabProvider(tt.args.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGitlabProvider() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGitlabProvider() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGitlabProvider_GetHeaderKeys(t *testing.T) {
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
			want: []string{XGitlabEvent, ContentTypeHeader},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &GitlabProvider{
				secret: tt.fields.secret,
			}
			if got := p.GetHeaderKeys(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GitlabProvider.GetHeaderKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGitlabProvider_Validate(t *testing.T) {
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
			name: "TestValidateWithEmptyTokenValue",
			fields: fields{
				secret: gitlabTestSecret,
			},
			args: args{
				hook: Hook{
					Headers: map[string]string{
						XGitlabToken: "",
					},
				},
			},
			want: false,
		},
		{
			name: "TestValidateWithEmptyHeaders",
			fields: fields{
				secret: gitlabTestSecret,
			},
			args: args{
				hook: Hook{
					Headers: map[string]string{},
				},
			},
			want: false,
		},
		{
			name: "TestValidateWithWrongTokenValue",
			fields: fields{
				secret: gitlabTestSecret,
			},
			args: args{
				hook: Hook{
					Headers: map[string]string{
						XGitlabToken: "IncorrectSecret",
					},
					Payload: nil,
				},
			},
			want: false,
		},
		{
			name: "TestValidateWithCorrectTokenValue",
			fields: fields{
				secret: gitlabTestSecret,
			},
			args: args{
				hook: Hook{
					Headers: map[string]string{
						XGitlabToken: gitlabTestSecret,
					},
					Payload: nil,
				},
			},
			want: true,
		},
		{
			name: "TestValidateWithWrongHeaderKey",
			fields: fields{
				secret: gitlabTestSecret,
			},
			args: args{
				hook: Hook{
					Headers: map[string]string{
						"X-Wrong-Token": "IncorrectSecret",
					},
					Payload: nil,
				},
			},
			want: false,
		},
		{
			name: "TestValidateWithNilHeaders",
			fields: fields{
				secret: gitlabTestSecret,
			},
			args: args{
				hook: Hook{
					Headers: nil,
					Payload: nil,
				},
			},
			want: false,
		},
		{
			name: "TestValidateWithNoHookArg",
			fields: fields{
				secret: gitlabTestSecret,
			},
			args: args{},
			want: false,
		},
		{
			name: "TestValidateWithWrongSecretInProxy",
			fields: fields{
				secret: "WrongSecret",
			},
			args: args{
				hook: Hook{
					Headers: map[string]string{
						XGitlabToken: gitlabTestSecret,
					},
					Payload: nil,
				},
			},
			want: false,
		},
		{
			name: "TestValidateWithEmptySecretInProxy",
			fields: fields{
				secret: "",
			},
			args: args{
				hook: Hook{
					Headers: map[string]string{
						XGitlabToken: gitlabTestSecret,
					},
					Payload: nil,
				},
			},
			want: false,
		},
		{
			name:   "TestValidateWithNoSecretInProxy",
			fields: fields{},
			args: args{
				hook: Hook{
					Headers: map[string]string{
						XGitlabToken: gitlabTestSecret,
					},
					Payload: nil,
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &GitlabProvider{
				secret: tt.fields.secret,
			}
			if got := p.Validate(tt.args.hook); got != tt.want {
				t.Errorf("GitlabProvider.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
