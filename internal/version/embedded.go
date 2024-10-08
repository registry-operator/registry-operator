// Copyright Registry Operator contributors
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

package version

import (
	"sigs.k8s.io/yaml"

	_ "embed"
)

//go:embed values.yaml
var valuesRaw []byte

func init() {
	err := yaml.Unmarshal(valuesRaw, &config)
	if err != nil {
		panic(err)
	}
}

type EmbeddedConfig struct {
	Registry EmbeddedRegistry `yaml:"registry"`
}

type EmbeddedRegistry struct {
	Image EmbeddedImage `yaml:"image"`
}

type EmbeddedImage struct {
	Repository string `yaml:"repository"`
	Tag        string `yaml:"tag"`
}

var config EmbeddedConfig

func GetRegistryImage() string {
	return config.Registry.Image.Repository + ":" + config.Registry.Image.Tag
}

func GetRegistryVersion() string {
	return config.Registry.Image.Tag
}
