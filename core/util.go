// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
)

type funcID int

func (id funcID) String() string {
	return fmt.Sprint(int(id))
}

func parseFuncID(s string) (id funcID, _ error) {
	n, err := fmt.Sscanf(s, "%d", &id)
	if n == 1 {
		return id, nil
	}
	return 0, err
}

func imin(a, b int) int {
	if a < b {
		return a
	}
	return b
}
