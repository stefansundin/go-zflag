// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pflag

import "fmt"

func ExampleShorthandLookup() {
	name := "verbose"
	short := name[:1]

	BoolP(name, short, false, "verbose output")

	// len(short) must be == 1
	flag := ShorthandLookup(short)

	fmt.Println(flag.Name)
}

func ExampleFlagSet_ShorthandLookup() {
	name := "verbose"
	short := name[:1]

	fs := NewFlagSet("Example", ContinueOnError)
	fs.BoolP(name, short, false, "verbose output")

	// len(short) must be == 1
	flag := fs.ShorthandLookup(short)

	fmt.Println(flag.Name)
}
