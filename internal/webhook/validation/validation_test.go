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
