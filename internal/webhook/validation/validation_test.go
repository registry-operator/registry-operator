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

package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"k8s.io/utils/ptr"
)

func TestHasAtMostOne(t *testing.T) {
	t.Parallel()

	type input = struct {
		Foo *int
		Bar *int
	}

	for name, tc := range map[string]struct {
		input    any
		expected bool
	}{
		"nil": {
			expected: true,
		},
		"zero": {
			input:    input{},
			expected: true,
		},
		"one": {
			input:    input{Foo: ptr.To(1)},
			expected: true,
		},
		"two": {
			input: input{Foo: ptr.To(1), Bar: ptr.To(2)},
		},
		"ptr zero": {
			input:    &input{},
			expected: true,
		},
		"ptr one": {
			input:    &input{Foo: ptr.To(1)},
			expected: true,
		},
		"ptr two": {
			input: &input{Foo: ptr.To(1), Bar: ptr.To(2)},
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			actual := HasAtMostOne(tc.input)
			assert.Equal(t, tc.expected, actual)
		})
	}

}
