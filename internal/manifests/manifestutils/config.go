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

	"github.com/distribution/distribution/v3/configuration"

	registryv1alpha1 "github.com/registry-operator/registry-operator/api/v1alpha1"
)

func GenerateConfig(registry registryv1alpha1.RegistrySpec) *configuration.Configuration {
	// Source: https://github.com/distribution/distribution/blob/v3.0.0-rc.3/cmd/registry/config-dev.yml
	return &configuration.Configuration{
		Version: "0.1",
		Log: configuration.Log{
			Level: "debug",
			Fields: map[string]interface{}{
				"service":     "registry",
				"environment": "operator-default",
			},
		},
		Storage: configuration.Storage{
			"delete": configuration.Parameters{
				"enabled": true,
			},
			"cache": configuration.Parameters{
				"blobdescriptor": "inmemory",
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
		HTTP: configuration.HTTP{
			Addr: ":5000",
			Debug: configuration.Debug{
				Addr: ":5001",
				Prometheus: configuration.Prometheus{
					Enabled: true,
					Path:    "/metrics",
				},
			},
			Headers: http.Header{
				"X-Content-Type-Options": []string{"nosniff"},
			},
		},
		Health: configuration.Health{
			StorageDriver: configuration.StorageDriver{
				Enabled:   true,
				Interval:  time.Second * 10,
				Threshold: 3,
			},
		},
	}
}
