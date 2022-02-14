// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import "fmt"

func ExampleShorthandLookup() {
	name := "verbose"
	short := 'v'

	Bool(name, false, "verbose output", OptShorthand(short))

	// len(short) must be == 1
	flag := ShorthandLookup(short)

	fmt.Println(flag.Name)
}

func ExampleFlagSet_ShorthandLookup() {
	name := "verbose"
	short := 'v'

	fs := NewFlagSet("Example", ContinueOnError)
	fs.Bool(name, false, "verbose output", OptShorthand(short))

	// len(short) must be == 1
	flag := fs.ShorthandLookup(short)

	fmt.Println(flag.Name)
}
