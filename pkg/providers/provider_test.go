package providers

import (
	"reflect"
	"testing"
)

func TestNewProvider(t *testing.T) {
	type args struct {
		provider string
		secret   string
	}
	tests := []struct {
		name    string
		args    args
		want    Provider
		wantErr bool
	}{
		{
			name: "TestNewProviderWithCorrectGithubProviderAndSecret",
			args: args{
				provider: GithubProviderKind,
				secret:   githubTestSecret,
			},
			want: &GithubProvider{
				secret: githubTestSecret,
			},
		},
		{
			name: "TestNewProviderWithEmptyProviderAndSecret",
			args: args{
				provider: "",
				secret:   githubTestSecret,
			},
			wantErr: true,
		},
		{
			name: "TestNewProviderWithGithubProviderAndEmptySecret",
			args: args{
				provider: GithubProviderKind,
				secret:   "",
			},
			want: &GithubProvider{
				secret: "",
			},
		},
		{
			name: "TestNewProviderWithEmptyGithubProviderAndEmptySecret",
			args: args{
				provider: "",
				secret:   "",
			},
			wantErr: true,
		},
		{
			name: "TestNewProviderWithGitlabProviderSecret",
			args: args{
				provider: GitlabProviderKind,
				secret:   gitlabTestSecret,
			},
			want: &GitlabProvider{
				secret: gitlabTestSecret,
			},
		},
		{
			name: "TestNewProviderWithIncorrectProviderKind",
			args: args{
				provider: "incorrectprovider",
				secret:   gitlabTestSecret,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewProvider(tt.args.provider, tt.args.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewProvider() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewProvider() = %v, want %v", got, tt.want)
			}
		})
	}
}
