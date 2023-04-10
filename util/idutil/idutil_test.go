// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package idutil

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUUID36(t *testing.T) {
	fmt.Println(GetUUID36(""))
}

func TestRandString(t *testing.T) {
	str := randString(Alphabet62, 50)
	assert.Equal(t, 50, len(str))
	t.Log(str)

	str = randString(Alphabet62, 255)
	assert.Equal(t, 255, len(str))
	t.Log(str)
}
