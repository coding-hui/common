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

package field

import "testing"

func TestPath(t *testing.T) {
	testCases := []struct {
		op       func(*Path) *Path
		expected string
	}{
		{
			func(p *Path) *Path { return p },
			"root",
		},
		{
			func(p *Path) *Path { return p.Child("first") },
			"root.first",
		},
		{
			func(p *Path) *Path { return p.Child("second") },
			"root.first.second",
		},
		{
			func(p *Path) *Path { return p.Index(0) },
			"root.first.second[0]",
		},
		{
			func(p *Path) *Path { return p.Child("third") },
			"root.first.second[0].third",
		},
		{
			func(p *Path) *Path { return p.Index(93) },
			"root.first.second[0].third[93]",
		},
		{
			func(p *Path) *Path { return p.parent },
			"root.first.second[0].third",
		},
		{
			func(p *Path) *Path { return p.parent },
			"root.first.second[0]",
		},
		{
			func(p *Path) *Path { return p.Key("key") },
			"root.first.second[0][key]",
		},
	}

	root := NewPath("root")
	p := root
	for i, tc := range testCases {
		p = tc.op(p)
		if p.String() != tc.expected {
			t.Errorf("[%d] Expected %q, got %q", i, tc.expected, p.String())
		}
		if p.Root() != root {
			t.Errorf("[%d] Wrong root: %#v", i, p.Root())
		}
	}
}

func TestPathMultiArg(t *testing.T) {
	testCases := []struct {
		op       func(*Path) *Path
		expected string
	}{
		{
			func(p *Path) *Path { return p },
			"root.first",
		},
		{
			func(p *Path) *Path { return p.Child("second", "third") },
			"root.first.second.third",
		},
		{
			func(p *Path) *Path { return p.Index(0) },
			"root.first.second.third[0]",
		},
		{
			func(p *Path) *Path { return p.parent },
			"root.first.second.third",
		},
		{
			func(p *Path) *Path { return p.parent },
			"root.first.second",
		},
		{
			func(p *Path) *Path { return p.parent },
			"root.first",
		},
		{
			func(p *Path) *Path { return p.parent },
			"root",
		},
	}

	root := NewPath("root", "first")
	p := root
	for i, tc := range testCases {
		p = tc.op(p)
		if p.String() != tc.expected {
			t.Errorf("[%d] Expected %q, got %q", i, tc.expected, p.String())
		}
		if p.Root() != root.Root() {
			t.Errorf("[%d] Wrong root: %#v", i, p.Root())
		}
	}
}
