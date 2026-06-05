package main

import (
	"os"
	"strings"
	"testing"

	"github.com/spf13/pflag"
)

func TestEnvironmentVariables(t *testing.T) {
	// Save original args and env
	origArgs := os.Args
	origEnv := os.Environ()
	origCommandLine := pflag.CommandLine
	defer func() {
		os.Args = origArgs
		os.Clearenv()
		for _, e := range origEnv {
			parts := strings.SplitN(e, "=", 2)
			if len(parts) == 2 {
				os.Setenv(parts[0], parts[1])
			}
		}
		pflag.CommandLine = origCommandLine
	}()

	tests := []struct {
		name     string
		envVars  map[string]string
		args     []string
		wantURL  string
		wantFail bool
	}{
		{
			name: "env var uppercase",
			envVars: map[string]string{
				"GWP_UPSTREAMURL": "https://jenkins.example.com",
			},
			args:    []string{"program"},
			wantURL: "https://jenkins.example.com",
		},
		{
			name: "env var mixed case",
			envVars: map[string]string{
				"GWP_upstreamURL": "https://jenkins.example.com",
			},
			args:    []string{"program"},
			wantURL: "https://jenkins.example.com",
		},
		{
			name: "command line flag takes precedence",
			envVars: map[string]string{
				"GWP_UPSTREAMURL": "https://jenkins.example.com",
			},
			args:    []string{"program", "--upstreamURL=https://jenkins2.example.com"},
			wantURL: "https://jenkins2.example.com",
		},
		{
			name: "multiple env vars",
			envVars: map[string]string{
				"GWP_UPSTREAMURL":  "https://jenkins.example.com",
				"GWP_ALLOWEDPATHS": "/path1,/path2",
				"GWP_PROVIDER":     "gitlab",
			},
			args:    []string{"program"},
			wantURL: "https://jenkins.example.com",
		},
		{
			name:     "no url specified",
			envVars:  map[string]string{},
			args:     []string{"program"},
			wantFail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset environment and args
			os.Clearenv()
			os.Args = tt.args

			// Create new FlagSet and define flags
			pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError)
			upstreamURL = pflag.String("upstreamURL", "", "URL to which the proxy requests will be forwarded (required)")
			secret = pflag.String("secret", "", "Secret of the Webhook API")
			provider = pflag.String("provider", "github", "Git Provider which generates the Webhook")
			allowedPaths = pflag.String("allowedPaths", "", "Comma-Separated String List of allowed paths")
			ignoredUsers = pflag.String("ignoredUsers", "", "Comma-Separated String List of users to ignore")
			allowedUsers = pflag.String("allowedUser", "", "Comma-Separated String List of users to allow")
			listenAddress = pflag.String("listen", ":8080", "Address on which the proxy listens")

			// Set up test environment variables
			for k, v := range tt.envVars {
				if err := os.Setenv(k, v); err != nil {
					t.Fatal(err)
				}
			}

			// Parse command line flags first
			pflag.CommandLine.ParseErrorsWhitelist.UnknownFlags = true
			if err := pflag.CommandLine.Parse(tt.args[1:]); err != nil {
				t.Fatal(err)
			}

			// Then process environment variables
			for _, env := range os.Environ() {
				if strings.HasPrefix(env, "GWP_") {
					parts := strings.SplitN(env, "=", 2)
					if len(parts) != 2 {
						continue
					}

					envName := strings.TrimPrefix(parts[0], "GWP_")
					value := parts[1]

					// Map environment variables to flag names
					var name string
					switch strings.ToUpper(envName) {
					case "UPSTREAMURL":
						name = "upstreamURL"
					case "ALLOWEDPATHS":
						name = "allowedPaths"
					case "IGNOREDUSERS":
						name = "ignoredUsers"
					case "ALLOWEDUSER":
						name = "allowedUser"
					case "PROVIDER":
						name = "provider"
					case "SECRET":
						name = "secret"
					case "LISTEN":
						name = "listen"
					default:
						name = strings.ToLower(envName)
					}

					if flag := pflag.Lookup(name); flag != nil && !flag.Changed {
						if err := flag.Value.Set(value); err != nil {
							t.Errorf("Error setting flag %s: %v", name, err)
						}
					}
				}
			}

			if tt.wantFail {
				defer func() {
					if r := recover(); r == nil {
						t.Error("Expected panic but got none")
					}
				}()
				validateRequiredFlags()
			} else {
				validateRequiredFlags()
				if got := *upstreamURL; got != tt.wantURL {
					t.Errorf("got upstreamURL = %v, want %v", got, tt.wantURL)
				}
			}
		})
	}
}
