// Copyright 2025 The Registry Operator Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package manifestutils

import (
	"net/http"
	"time"

	"github.com/docker/distribution/configuration"

	registryv1alpha1 "github.com/registry-operator/registry-operator/api/v1alpha1"
)

func GenerateConfig(registry registryv1alpha1.RegistrySpec) *configuration.Configuration {
	// yes, this is messy, because the config in distribution 2.8.3 is inlinded...
	// Source: https://github.com/distribution/distribution/blob/v2.8.3/cmd/registry/config-dev.yml
	return &configuration.Configuration{
		Version: "0.1",
		Log: struct {
			AccessLog struct {
				Disabled bool "yaml:\"disabled,omitempty\""
			} "yaml:\"accesslog,omitempty\""
			Level     configuration.Loglevel  "yaml:\"level,omitempty\""
			Formatter string                  "yaml:\"formatter,omitempty\""
			Fields    map[string]interface{}  "yaml:\"fields,omitempty\""
			Hooks     []configuration.LogHook "yaml:\"hooks,omitempty\""
		}{
			Level: "debug",
			Fields: map[string]interface{}{
				"service":     "registry",
				"environment": "operator-default",
			},
			Hooks: []configuration.LogHook{
				{
					Type:     "mail",
					Disabled: true,
					Levels: []string{
						"panic",
					},
					MailOptions: configuration.MailOptions{
						From: "sender@example.com",
						To: []string{
							"errors@example.com",
						},
						SMTP: struct {
							Addr     string "yaml:\"addr,omitempty\""
							Username string "yaml:\"username,omitempty\""
							Password string "yaml:\"password,omitempty\""
							Insecure bool   "yaml:\"insecure,omitempty\""
						}{
							Addr:     "mail.example.com:25",
							Username: "mailuser",
							Password: "password",
							Insecure: true,
						},
					},
				},
			},
		},
		Storage: configuration.Storage{
			"delete": configuration.Parameters{
				"enabled": true,
			},
			"cache": configuration.Parameters{
				"blobdescriptor": "redis",
			},
			"filesystem": configuration.Parameters{
				"rootdirectory": "/var/lib/registry",
			},
			"maintenance": configuration.Parameters{
				"uploadpurging": map[string]interface{}{
					"enabled": false,
				},
			},
		},
		HTTP: struct {
			Addr         string        "yaml:\"addr,omitempty\""
			Net          string        "yaml:\"net,omitempty\""
			Host         string        "yaml:\"host,omitempty\""
			Prefix       string        "yaml:\"prefix,omitempty\""
			Secret       string        "yaml:\"secret,omitempty\""
			RelativeURLs bool          "yaml:\"relativeurls,omitempty\""
			DrainTimeout time.Duration "yaml:\"draintimeout,omitempty\""
			TLS          struct {
				Certificate  string   "yaml:\"certificate,omitempty\""
				Key          string   "yaml:\"key,omitempty\""
				ClientCAs    []string "yaml:\"clientcas,omitempty\""
				MinimumTLS   string   "yaml:\"minimumtls,omitempty\""
				CipherSuites []string "yaml:\"ciphersuites,omitempty\""
				LetsEncrypt  struct {
					CacheFile string   "yaml:\"cachefile,omitempty\""
					Email     string   "yaml:\"email,omitempty\""
					Hosts     []string "yaml:\"hosts,omitempty\""
				} "yaml:\"letsencrypt,omitempty\""
			} "yaml:\"tls,omitempty\""
			Headers http.Header "yaml:\"headers,omitempty\""
			Debug   struct {
				Addr       string "yaml:\"addr,omitempty\""
				Prometheus struct {
					Enabled bool   "yaml:\"enabled,omitempty\""
					Path    string "yaml:\"path,omitempty\""
				} "yaml:\"prometheus,omitempty\""
			} "yaml:\"debug,omitempty\""
			HTTP2 struct {
				Disabled bool "yaml:\"disabled,omitempty\""
			} "yaml:\"http2,omitempty\""
		}{
			Addr: ":5000",
			Debug: struct {
				Addr       string "yaml:\"addr,omitempty\""
				Prometheus struct {
					Enabled bool   "yaml:\"enabled,omitempty\""
					Path    string "yaml:\"path,omitempty\""
				} "yaml:\"prometheus,omitempty\""
			}{
				Addr: ":5001",
				Prometheus: struct {
					Enabled bool   "yaml:\"enabled,omitempty\""
					Path    string "yaml:\"path,omitempty\""
				}{
					Enabled: true,
					Path:    "/metrics",
				},
			},
			Headers: http.Header{
				"X-Content-Type-Options": []string{"nosniff"},
			},
		},
		Redis: struct {
			Addr         string        "yaml:\"addr,omitempty\""
			Password     string        "yaml:\"password,omitempty\""
			DB           int           "yaml:\"db,omitempty\""
			DialTimeout  time.Duration "yaml:\"dialtimeout,omitempty\""
			ReadTimeout  time.Duration "yaml:\"readtimeout,omitempty\""
			WriteTimeout time.Duration "yaml:\"writetimeout,omitempty\""
			Pool         struct {
				MaxIdle     int           "yaml:\"maxidle,omitempty\""
				MaxActive   int           "yaml:\"maxactive,omitempty\""
				IdleTimeout time.Duration "yaml:\"idletimeout,omitempty\""
			} "yaml:\"pool,omitempty\""
		}{
			Addr: "localhost:6379",
			Pool: struct {
				MaxIdle     int           "yaml:\"maxidle,omitempty\""
				MaxActive   int           "yaml:\"maxactive,omitempty\""
				IdleTimeout time.Duration "yaml:\"idletimeout,omitempty\""
			}{
				MaxIdle:     16,
				MaxActive:   64,
				IdleTimeout: time.Second * 300,
			},
			DialTimeout:  time.Millisecond * 10,
			ReadTimeout:  time.Millisecond * 10,
			WriteTimeout: time.Millisecond * 10,
		},
		Notifications: configuration.Notifications{
			EventConfig: configuration.Events{
				IncludeReferences: true,
			},
			Endpoints: []configuration.Endpoint{
				{
					Name: "local-5003",
					URL:  "http://localhost:5003/callback",
					Headers: http.Header{
						"Authorization": []string{"Bearer <an example token>"},
					},
					Timeout:   time.Second * 1,
					Threshold: 10,
					Backoff:   time.Second * 1,
					Disabled:  true,
				},
				{
					Name:      "local-8083",
					URL:       "http://localhost:8083/callback",
					Timeout:   time.Second * 1,
					Threshold: 10,
					Backoff:   time.Second * 1,
					Disabled:  true,
				},
			},
		},
		Health: configuration.Health{
			StorageDriver: struct {
				Enabled   bool          "yaml:\"enabled,omitempty\""
				Interval  time.Duration "yaml:\"interval,omitempty\""
				Threshold int           "yaml:\"threshold,omitempty\""
			}{
				Enabled:   true,
				Interval:  time.Second * 10,
				Threshold: 3,
			},
		},
	}
}
