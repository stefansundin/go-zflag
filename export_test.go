// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pflag

import "os"

// Additional routines compiled into the package only during testing.

var DefaultUsage = Usage

// ResetForTesting clears all flag state and sets the usage function as directed.
// After calling ResetForTesting, parse errors in flag handling will not
// exit the program.
func ResetForTesting(usage func()) {
	CommandLine = NewFlagSet(os.Args[0], ContinueOnError)
	CommandLine.Usage = DefaultUsage
	Usage = usage
}
