// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package v1

import "strings"

// isSensitiveData returns whether the input string contains sensitive information
func isSensitiveData(key string, sensitiveKeys []string) bool {
	for _, v := range sensitiveKeys {
		if strings.Contains(strings.ToLower(key), v) {
			return true
		}
	}
	return false
}

// desensitize returns the desensitized data
func desensitize(data map[string]interface{}, sensitiveKeys []string) map[string]interface{} {
	output := make(map[string]interface{})
	for k, v := range data {
		if isSensitiveData(k, sensitiveKeys) {
			continue
		}
		switch v := v.(type) {
		case map[interface{}]interface{}:
			output[k] = desensitize(convert(v), sensitiveKeys)
		default:
			output[k] = v
		}
	}
	return output
}

// convert returns formatted data
func convert(m map[interface{}]interface{}) map[string]interface{} {
	output := make(map[string]interface{})
	for k, v := range m {
		switch k := k.(type) {
		case string:
			output[k] = v
		}
	}
	return output
}
