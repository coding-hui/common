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

package time

import (
	"fmt"
	"time"

	sqldriver "database/sql/driver"
)

const (
	defaultDateTimeFormat = "2006-01-02 15:04:05"
)

// Time format json time field by myself.
type Time struct {
	time.Time
}

// MarshalJSON on Time format Time field with %Y-%m-%d %H:%M:%S.
func (t Time) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", t.Format(defaultDateTimeFormat))
	return []byte(formatted), nil
}

// Value insert timestamp into mysql need this function.
func (t Time) Value() (sqldriver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan valueof time.Time.
func (t *Time) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = Time{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

// ToTime convert string to Time.
func ToTime(str string) (Time, error) {
	var jt Time
	loc, _ := time.LoadLocation("Local")
	value, err := time.ParseInLocation(defaultDateTimeFormat, str, loc)
	if err != nil {
		return jt, err
	}
	return Time{
		Time: value,
	}, nil
}

// Now returns the current time.
func Now() Time {
	return Time{
		Time: time.Now(),
	}
}
