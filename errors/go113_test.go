//go:build go1.13
// +build go1.13

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
	stderrors "errors"
	"fmt"
	"reflect"
	"testing"
)

func TestErrorChainCompat(t *testing.T) {
	err := stderrors.New("error that gets wrapped")
	wrapped := Wrap(err, "wrapped up")
	if !stderrors.Is(wrapped, err) {
		t.Errorf("Wrap does not support Go 1.13 error chains")
	}
}

func TestIs(t *testing.T) {
	err := New("test")

	type args struct {
		err    error
		target error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "with stack",
			args: args{
				err:    WithStack(err),
				target: err,
			},
			want: true,
		},
		{
			name: "with message",
			args: args{
				err:    WithMessage(err, "test"),
				target: err,
			},
			want: true,
		},
		{
			name: "with message format",
			args: args{
				err:    WithMessagef(err, "%s", "test"),
				target: err,
			},
			want: true,
		},
		{
			name: "std errors compatibility",
			args: args{
				err:    fmt.Errorf("wrap it: %w", err),
				target: err,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Is(tt.args.err, tt.args.target); got != tt.want {
				t.Errorf("Is() = %v, want %v", got, tt.want)
			}
		})
	}
}

type customErr struct {
	msg string
}

func (c customErr) Error() string { return c.msg }

func TestAs(t *testing.T) {
	var err = customErr{msg: "test message"}

	type args struct {
		err    error
		target interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "with stack",
			args: args{
				err:    WithStack(err),
				target: new(customErr),
			},
			want: true,
		},
		{
			name: "with message",
			args: args{
				err:    WithMessage(err, "test"),
				target: new(customErr),
			},
			want: true,
		},
		{
			name: "with message format",
			args: args{
				err:    WithMessagef(err, "%s", "test"),
				target: new(customErr),
			},
			want: true,
		},
		{
			name: "std errors compatibility",
			args: args{
				err:    fmt.Errorf("wrap it: %w", err),
				target: new(customErr),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := As(tt.args.err, tt.args.target); got != tt.want {
				t.Errorf("As() = %v, want %v", got, tt.want)
			}

			ce := tt.args.target.(*customErr)
			if !reflect.DeepEqual(err, *ce) {
				t.Errorf("set target error failed, target error is %v", *ce)
			}
		})
	}
}

func TestUnwrap(t *testing.T) {
	err := New("test")

	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "with stack",
			args: args{err: WithStack(err)},
			want: err,
		},
		{
			name: "with message",
			args: args{err: WithMessage(err, "test")},
			want: err,
		},
		{
			name: "with message format",
			args: args{err: WithMessagef(err, "%s", "test")},
			want: err,
		},
		{
			name: "std errors compatibility",
			args: args{err: fmt.Errorf("wrap: %w", err)},
			want: err,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Unwrap(tt.args.err); !reflect.DeepEqual(err, tt.want) {
				t.Errorf("Unwrap() error = %v, want %v", err, tt.want)
			}
		})
	}
}
