/*
 * MIT License
 *
 * Copyright (c) 2023 WeCoding.
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

package validation

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Base is the interface for all configs used in Aptomi (e.g. client config, server config).
type Base interface {
	IsDebug() bool
}

type testStruct struct {
	Host     string `validate:"required,hostname|ip"`
	Port     int    `validate:"required,min=1,max=65535"`
	SomeDir  string `validate:"required,dir"`
	SomeFile string `validate:"omitempty,file"`
}

// writeTempFile creates a temporary file, writes given data into it and returns its name.
// It's up to a caller to delete the created temporary file by calling os.Remove() on its name.
func writeTempFile(prefix string, data []byte) string {
	tmpFile, err := ioutil.TempFile("", "aptomi-"+prefix)
	if err != nil {
		panic("Failed to create temp file")
	}
	defer tmpFile.Close()

	_, err = tmpFile.Write(data)
	if err != nil {
		panic("Failed to write to temp file")
	}

	return tmpFile.Name()
}

func (t *testStruct) IsDebug() bool {
	return false
}

func displayErrorMessages() bool {
	return false
}

func TestConfigValidation(t *testing.T) {
	tmpFile := writeTempFile("unittest", []byte("unittest"))
	defer os.Remove(tmpFile)

	tests := []struct {
		config Base
		result bool
	}{
		{
			&testStruct{
				Host:    "0.0.0.0",
				Port:    80,
				SomeDir: "/tmp",
			},
			true,
		},
		{
			&testStruct{
				Host:    "0.0.0.0",
				Port:    80,
				SomeDir: "",
			},
			false,
		},
		{
			&testStruct{
				Host:    "0.0.0.0",
				Port:    80,
				SomeDir: "/nonexistingdirectoryinroot",
			},
			false,
		},
		{
			&testStruct{
				Host:    "127.0.0.1",
				Port:    8080,
				SomeDir: "/tmp",
			},
			true,
		},
		{
			&testStruct{
				Host:    "10.20.30.40",
				Port:    65080,
				SomeDir: "/tmp",
			},
			true,
		},
		{
			&testStruct{
				Host:    "demo.aptomi.io",
				Port:    65080,
				SomeDir: "/tmp",
			},
			true,
		},
		{
			&testStruct{
				Host:    "0.0.0.0",
				Port:    0,
				SomeDir: "/tmp",
			},
			false,
		},
		{
			&testStruct{
				Host:    "0.0.0.0",
				Port:    -1,
				SomeDir: "/tmp",
			},
			false,
		},
		{
			&testStruct{
				Host:    "",
				Port:    80,
				SomeDir: "/tmp",
			},
			false,
		},
		{
			&testStruct{
				Host:     "0.0.0.0",
				Port:     80,
				SomeDir:  "/tmp",
				SomeFile: tmpFile,
			},
			true,
		},
		{
			&testStruct{
				Host:     "0.0.0.0",
				Port:     80,
				SomeDir:  "/tmp",
				SomeFile: tmpFile + ".non-existing",
			},
			false,
		},
	}
	for _, test := range tests {
		val := NewValidator(test.config)
		err := val.Validate()
		failed := !assert.Equal(t, test.result, err == nil, "Validation test case failed: %s", test.config)
		if err != nil {
			msg := err.ToAggregate().Error()
			if displayErrorMessages() || failed {
				t.Log(msg)
			}
		}
	}
}
