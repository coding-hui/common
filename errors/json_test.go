/*
 * MIT License
 *
 * Copyright (c) 2023.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package errors

import (
	"encoding/json"
	"regexp"
	"testing"
)

func TestFrameMarshalText(t *testing.T) {
	var tests = []struct {
		Frame
		want string
	}{{
		initpc,
		`^github.com/coding-hui/common/errors\.init(\.ializers)? .+/github\.com/marmotedu/errors/stack_test.go:\d+$`,
	}, {
		0,
		`^unknown$`,
	}}
	for i, tt := range tests {
		got, err := tt.Frame.MarshalText()
		if err != nil {
			t.Fatal(err)
		}
		if !regexp.MustCompile(tt.want).Match(got) {
			t.Errorf("test %d: MarshalJSON:\n got %q\n want %q", i+1, string(got), tt.want)
		}
	}
}

func TestFrameMarshalJSON(t *testing.T) {
	var tests = []struct {
		Frame
		want string
	}{{
		initpc,
		`^"github\.com/marmotedu/errors\.init(\.ializers)? .+/github\.com/marmotedu/errors/stack_test.go:\d+"$`,
	}, {
		0,
		`^"unknown"$`,
	}}
	for i, tt := range tests {
		got, err := json.Marshal(tt.Frame)
		if err != nil {
			t.Fatal(err)
		}
		if !regexp.MustCompile(tt.want).Match(got) {
			t.Errorf("test %d: MarshalJSON:\n got %q\n want %q", i+1, string(got), tt.want)
		}
	}
}
