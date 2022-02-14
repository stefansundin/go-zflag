// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import "strconv"

// -- uint Value
type uintValue uint

func newUintValue(val uint, p *uint) *uintValue {
	*p = val
	return (*uintValue)(p)
}

func (i *uintValue) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 64)
	*i = uintValue(v)
	return err
}

func (i *uintValue) Get() interface{} {
	return uint(*i)
}

func (i *uintValue) Type() string {
	return "uint"
}

func (i *uintValue) String() string { return strconv.FormatUint(uint64(*i), 10) }

// GetUint return the uint value of a flag with the given name
func (f *FlagSet) GetUint(name string) (uint, error) {
	val, err := f.getFlagType(name, "uint")
	if err != nil {
		return 0, err
	}
	return val.(uint), nil
}

// MustGetUint is like GetUint, but panics on error.
func (f *FlagSet) MustGetUint(name string) uint {
	val, err := f.GetUint(name)
	if err != nil {
		panic(err)
	}
	return val
}

// UintVar defines an uint flag with specified name, default value, and usage string.
// The argument p points to an uint variable in which to store the value of the flag.
func (f *FlagSet) UintVar(p *uint, name string, value uint, usage string, opts ...Opt) {
	f.Var(newUintValue(value, p), name, usage, opts...)
}

// UintVar defines an uint flag with specified name, default value, and usage string.
// The argument p points to an uint variable in which to store the value of the flag.
func UintVar(p *uint, name string, value uint, usage string, opts ...Opt) {
	CommandLine.UintVar(p, name, value, usage, opts...)
}

// Uint defines an uint flag with specified name, default value, and usage string.
// The return value is the address of an uint variable that stores the value of the flag.
func (f *FlagSet) Uint(name string, value uint, usage string, opts ...Opt) *uint {
	var p uint
	f.UintVar(&p, name, value, usage, opts...)
	return &p
}

// Uint defines an uint flag with specified name, default value, and usage string.
// The return value is the address of an uint variable that stores the value of the flag.
func Uint(name string, value uint, usage string, opts ...Opt) *uint {
	return CommandLine.Uint(name, value, usage, opts...)
}
