// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import "fmt"

type errUnknownFlag struct {
	name string
}

func NewUnknownFlagError(name string) error {
	return errUnknownFlag{name: name}
}

func (e errUnknownFlag) Error() string {
	dash := "--"
	if len(e.name) == 1 {
		dash = "-"
	}

	return fmt.Sprintf("unknown flag: %s", dash+e.name)
}
