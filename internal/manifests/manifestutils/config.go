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
	"context"
	"errors"
	"fmt"
	"maps"
	"net/http"
	"net/url"
	"time"

	"github.com/distribution/distribution/v3/configuration"

	registryv1alpha1 "github.com/registry-operator/registry-operator/api/v1alpha1"
	"github.com/registry-operator/registry-operator/internal/manifests"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GenerateConfig(
	registry registryv1alpha1.RegistrySpec,
	s3 configuration.Parameters,
) (*configuration.Configuration, error) {
	storage := configuration.Storage{
		"delete": configuration.Parameters{
			"enabled": true,
		},
		"cache": configuration.Parameters{
			"blobdescriptor": "inmemory",
		},
		"maintenance": configuration.Parameters{
			"uploadpurging": map[string]interface{}{
				"enabled": false,
			},
		},
		"tag": configuration.Parameters{
			"concurrencylimit": 8,
		},
	}

	switch {
	case len(s3) > 0:
		storage["s3"] = configuration.Parameters{
			"rootdirectory": "/registry",
		}
		maps.Insert(storage["s3"], maps.All(s3))

	default:
		storage["filesystem"] = configuration.Parameters{
			"rootdirectory": "/var/lib/registry",
		}
	}

	return &configuration.Configuration{
		Version: "0.1",
		Log: configuration.Log{
			Level: "debug",
			Fields: map[string]interface{}{
				"service":     "registry",
				"environment": "development",
			},
		},
		Storage: storage,
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
	}, nil
}

func NewS3Config(ctx context.Context, params manifests.Params) (configuration.Parameters, error) {
	s3 := params.Registry.Spec.Storage.S3
	if s3 == nil {
		return nil, nil
	}

	s3c := configuration.Parameters{}

	nn := client.ObjectKey{
		Namespace: params.Registry.GetNamespace(),
	}

	var (
		key  string
		err  error
		errs error
	)

	nn.Name = s3.BucketName.Name
	key = s3.BucketName.Key
	s3c["bucket"], err = getDataFromSecret(ctx, params.Client, nn, key)
	if err != nil {
		errs = errors.Join(errs, err)
	}

	nn.Name = s3.Region.Name
	key = s3.Region.Key
	s3c["region"], err = getDataFromSecret(ctx, params.Client, nn, key)
	if err != nil {
		errs = errors.Join(errs, err)
	}

	if opt := s3.AccessKey; opt != nil {
		nn.Name = s3.AccessKey.Name
		key = s3.AccessKey.Key
		s3c["accesskey"], err = getDataFromSecret(ctx, params.Client, nn, key)
		if err != nil {
			errs = errors.Join(errs, err)
		}
	}

	if opt := s3.SecretKey; opt != nil {
		nn.Name = s3.SecretKey.Name
		key = s3.SecretKey.Key
		s3c["secretkey"], err = getDataFromSecret(ctx, params.Client, nn, key)
		if err != nil {
			errs = errors.Join(errs, err)
		}
	}

	if opt := s3.EndpointURL; opt != nil {
		nn.Name = s3.EndpointURL.Name
		key = s3.EndpointURL.Key
		var endpoint string
		endpoint, err = getDataFromSecret(ctx, params.Client, nn, key)
		if err != nil {
			errs = errors.Join(errs, err)
		} else {
			u, err := url.Parse(endpoint)
			if err != nil {
				return nil, fmt.Errorf("invalid endpoint URL: %w", err)
			}

			if u.Scheme == "" {
				u.Scheme = "https"
				endpoint = u.String()
			} else if u.Scheme == "http" {
				s3c["secure"] = false
			}

			s3c["regionendpoint"] = endpoint
		}
	}

	return s3c, errs
}

func getDataFromSecret(
	ctx context.Context,
	cli client.Client,
	nn client.ObjectKey,
	key string,
) (string, error) {
	sec := &corev1.Secret{}

	if err := cli.Get(ctx, nn, sec); err != nil {
		return "", fmt.Errorf("failed to fetch secret %v: %w", nn, err)
	}

	val, ok := sec.Data[key]
	if !ok {
		return "", fmt.Errorf("value for %s not found in %s", key, nn)
	}

	return string(val), nil
}
