//go:build go1.7
// +build go1.7

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
	"fmt"
	"testing"

	stderrors "errors"

	pkgerrors "github.com/pkg/errors"
)

func stdErrors(at, depth int) error {
	if at >= depth {
		return stderrors.New("no error")
	}
	return stdErrors(at+1, depth)
}

func pkgErrors(at, depth int) error {
	if at >= depth {
		return pkgerrors.New("no error")
	}
	return pkgErrors(at+1, depth)
}

func marmotErrors(at, depth int) error {
	if at >= depth {
		return WithCode(unknownCoder.Code(), "ye error")
	}
	return marmotErrors(at+1, depth)
}

// GlobalE is an exported global to store the result of benchmark results,
// preventing the compiler from optimising the benchmark functions away.
var GlobalE interface{}

func BenchmarkErrors(b *testing.B) {
	type run struct {
		stack int
		pkg   string
	}
	runs := []run{
		{10, "std"},
		{10, "pkg"},
		{10, "marmot"},
		{100, "std"},
		{100, "pkg"},
		{100, "marmot"},
		{1000, "std"},
		{1000, "pkg"},
		{1000, "marmot"},
	}
	for _, r := range runs {
		var part string
		var f func(at, depth int) error
		switch r.pkg {
		case "std":
			part = "errors"
			f = stdErrors
		case "pkg":
			part = "pkg/errors"
			f = pkgErrors
		case "marmot":
			part = "marmot/errors"
			f = marmotErrors
		default:
		}

		name := fmt.Sprintf("%s-stack-%d", part, r.stack)
		b.Run(name, func(b *testing.B) {
			var err error
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				err = f(0, r.stack)
			}
			b.StopTimer()
			GlobalE = err
		})
	}
}

func BenchmarkStackFormatting(b *testing.B) {
	type run struct {
		stack  int
		format string
	}
	runs := []run{
		{10, "%s"},
		{10, "%v"},
		{10, "%+v"},
		{30, "%s"},
		{30, "%v"},
		{30, "%+v"},
		{60, "%s"},
		{60, "%v"},
		{60, "%+v"},
	}

	var stackStr string
	for _, r := range runs {
		name := fmt.Sprintf("%s-stack-%d", r.format, r.stack)
		b.Run(name, func(b *testing.B) {
			err := stdErrors(0, r.stack)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				stackStr = fmt.Sprintf(r.format, err)
			}
			b.StopTimer()
		})
	}

	for _, r := range runs {
		name := fmt.Sprintf("%s-stacktrace-%d", r.format, r.stack)
		b.Run(name, func(b *testing.B) {
			err := stdErrors(0, r.stack)
			st := err.(*fundamental).stack.StackTrace()
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				stackStr = fmt.Sprintf(r.format, st)
			}
			b.StopTimer()
		})
	}
	GlobalE = stackStr
}

func BenchmarkCodeFormatting(b *testing.B) {
	type run struct {
		stack  int
		format string
	}
	runs := []run{
		{10, "%s"},
		{10, "%v"},
		{10, "%-v"},
		{10, "%+v"},
		{10, "%#v"},
		{10, "%#-v"},
		{10, "%#+v"},
		{30, "%s"},
		{30, "%v"},
		{30, "%-v"},
		{30, "%+v"},
		{30, "%#v"},
		{30, "%#-v"},
		{30, "%#+v"},
		{60, "%s"},
		{60, "%v"},
		{60, "%-v"},
		{60, "%+v"},
		{60, "%#v"},
		{60, "%#-v"},
		{60, "%#+v"},
	}

	var stackStr string
	for _, r := range runs {
		name := fmt.Sprintf("%s-stack-%d", r.format, r.stack)
		b.Run(name, func(b *testing.B) {
			err := stdErrors(0, r.stack)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				stackStr = fmt.Sprintf(r.format, err)
			}
			b.StopTimer()
		})
	}

	for _, r := range runs {
		name := fmt.Sprintf("%s-stacktrace-%d", r.format, r.stack)
		b.Run(name, func(b *testing.B) {
			err := stdErrors(0, r.stack)
			st := err.(*fundamental).stack.StackTrace()
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				stackStr = fmt.Sprintf(r.format, st)
			}
			b.StopTimer()
		})
	}
	GlobalE = stackStr
}
