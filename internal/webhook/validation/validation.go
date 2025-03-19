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
	"fmt"
	"reflect"
)

// HasAtMostOne checks if the provided value is either a struct (or pointer to struct)
// or a map (or pointer to map) and that it has at most one non-zero member.
// For any other type, it panics.
func HasAtMostOne(val any) bool {
	return PopulatedFields(val) <= 1
}

// PopulatedFields checks if the provided value is either a struct (or pointer to struct)
// or a map (or pointer to map) and then it counts non-zero members.
// For any other type, it panics.
func PopulatedFields(val any) int {
	if val == nil {
		return 0
	}
	v := reflect.ValueOf(val)
	// Dereference pointers
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return 0
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Struct:
		count := 0
		numField := v.NumField()
		for i := 0; i < numField; i++ {
			// Only consider exported fields (optional; remove if you want all fields)
			field := v.Field(i)
			if !field.IsZero() {
				count++
			}
		}
		return count

	case reflect.Map:
		count := 0
		for _, key := range v.MapKeys() {
			elem := v.MapIndex(key)
			if !elem.IsZero() {
				count++
			}
		}
		return count

	default:
		panic(fmt.Sprintf("unsupported type: %s", v.Type()))
	}
}
