package config_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/SudiMango/stratus/internal/config"
)

func TestProxyValidate(t *testing.T) {
	certificatePath, keyPath := writeCertificatePair(t)
	validCertificate := config.Certificate{
		Host: "app.example.com",
		Cert: certificatePath,
		Key:  keyPath,
	}

	tests := []struct {
		name        string
		proxy       config.Proxy
		wantMessage string
	}{
		{
			name: "HTTP only is valid",
			proxy: config.Proxy{
				Http: config.Http{Enabled: true, Port: 8080},
			},
		},
		{
			name: "HTTPS only is valid",
			proxy: config.Proxy{
				Https: config.Https{
					Enabled: true,
					Port:    8443,
				},
				Certificates: config.Certificates{validCertificate},
			},
		},
		{
			name: "multiple HTTPS certificates are valid",
			proxy: config.Proxy{
				Https: config.Https{Enabled: true, Port: 8443},
				Certificates: config.Certificates{
					validCertificate,
					{Host: "api.example.com", Cert: certificatePath, Key: keyPath},
				},
			},
		},
		{
			name:        "at least one listener is required",
			proxy:       config.Proxy{},
			wantMessage: "at least one listener",
		},
		{
			name: "redirect requires HTTPS",
			proxy: config.Proxy{
				Http: config.Http{Enabled: true, Port: 8080, RedirectToHttps: true},
			},
			wantMessage: "cannot redirect",
		},
		{
			name: "redirect requires HTTP",
			proxy: config.Proxy{
				Http: config.Http{Enabled: false, RedirectToHttps: true},
				Https: config.Https{
					Enabled: true,
					Port:    8443,
				},
				Certificates: config.Certificates{validCertificate},
			},
			wantMessage: "redirect",
		},
		{
			name: "listeners cannot share a port",
			proxy: config.Proxy{
				Http: config.Http{Enabled: true, Port: 8443},
				Https: config.Https{
					Enabled: true,
					Port:    8443,
				},
				Certificates: config.Certificates{validCertificate},
			},
			wantMessage: "same port",
		},
		{
			name: "HTTP port must be positive",
			proxy: config.Proxy{
				Http: config.Http{Enabled: true, Port: 0},
			},
			wantMessage: "HTTP port",
		},
		{
			name: "HTTPS port cannot exceed 65535",
			proxy: config.Proxy{
				Https: config.Https{
					Enabled: true,
					Port:    65536,
				},
				Certificates: config.Certificates{validCertificate},
			},
			wantMessage: "HTTPS port",
		},
		{
			name: "negative timeout is invalid",
			proxy: config.Proxy{
				Http:        config.Http{Enabled: true, Port: 8080},
				ReadTimeout: -time.Second,
			},
			wantMessage: "timeout",
		},
		{
			name: "negative maximum header bytes is invalid",
			proxy: config.Proxy{
				Http:           config.Http{Enabled: true, Port: 8080},
				MaxHeaderBytes: -1,
			},
			wantMessage: "header",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.proxy.Validate()
			assertValidationResult(t, err, test.wantMessage)
		})
	}
}

func TestProxyValidateCertificates(t *testing.T) {
	certificatePath, keyPath := writeCertificatePair(t)
	directory := t.TempDir()
	invalidKeyPath := filepath.Join(directory, "invalid-key.pem")
	if err := os.WriteFile(invalidKeyPath, []byte("not a private key"), 0o600); err != nil {
		t.Fatalf("write invalid key fixture: %v", err)
	}

	tests := []struct {
		name        string
		certificate config.Certificate
		wantMessage string
	}{
		{
			name: "matching certificate and key",
			certificate: config.Certificate{
				Host: "app.example.com",
				Cert: certificatePath,
				Key:  keyPath,
			},
		},
		{
			name: "certificate file does not exist",
			certificate: config.Certificate{
				Host: "app.example.com",
				Cert: filepath.Join(t.TempDir(), "missing-cert.pem"),
				Key:  keyPath,
			},
			wantMessage: "certificate",
		},
		{
			name: "private key is malformed",
			certificate: config.Certificate{
				Host: "app.example.com",
				Cert: certificatePath,
				Key:  invalidKeyPath,
			},
			wantMessage: "key",
		},
		{
			name: "host is required",
			certificate: config.Certificate{
				Cert: certificatePath,
				Key:  keyPath,
			},
			wantMessage: "host",
		},
		{
			name: "certificate path is required",
			certificate: config.Certificate{
				Host: "app.example.com",
				Key:  keyPath,
			},
			wantMessage: "certificate path",
		},
		{
			name: "key path is required",
			certificate: config.Certificate{
				Host: "app.example.com",
				Cert: certificatePath,
			},
			wantMessage: "key path",
		},
		{
			name: "certificate must match host",
			certificate: config.Certificate{
				Host: "unrelated.test",
				Cert: certificatePath,
				Key:  keyPath,
			},
			wantMessage: "not valid for host",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			proxy := config.Proxy{
				Https:        config.Https{Enabled: true, Port: 8443},
				Certificates: config.Certificates{test.certificate},
			}
			err := proxy.Validate()
			assertValidationResult(t, err, test.wantMessage)
		})
	}
}

func assertValidationResult(t *testing.T, err error, wantMessage string) {
	t.Helper()

	if wantMessage == "" {
		if err != nil {
			t.Fatalf("Validate() error = %v, want nil", err)
		}
		return
	}

	if err == nil {
		t.Fatalf("Validate() error = nil, want message containing %q", wantMessage)
	}
	if !strings.Contains(strings.ToLower(err.Error()), strings.ToLower(wantMessage)) {
		t.Errorf("Validate() error = %q, want message containing %q", err, wantMessage)
	}
}
