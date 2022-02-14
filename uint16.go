// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import "strconv"

// -- uint16 value
type uint16Value uint16

func newUint16Value(val uint16, p *uint16) *uint16Value {
	*p = val
	return (*uint16Value)(p)
}

func (i *uint16Value) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 16)
	*i = uint16Value(v)
	return err
}

func (i *uint16Value) Get() interface{} {
	return uint16(*i)
}

func (i *uint16Value) Type() string {
	return "uint16"
}

func (i *uint16Value) String() string { return strconv.FormatUint(uint64(*i), 10) }

// GetUint16 return the uint16 value of a flag with the given name
func (f *FlagSet) GetUint16(name string) (uint16, error) {
	val, err := f.getFlagType(name, "uint16")
	if err != nil {
		return 0, err
	}
	return val.(uint16), nil
}

// MustGetUint16 is like GetUint16, but panics on error.
func (f *FlagSet) MustGetUint16(name string) uint16 {
	val, err := f.GetUint16(name)
	if err != nil {
		panic(err)
	}
	return val
}

// Uint16Var defines an uint16 flag with specified name, default value, and usage string.
// The argument p points to an uint16 variable in which to store the value of the flag.
func (f *FlagSet) Uint16Var(p *uint16, name string, value uint16, usage string, opts ...Opt) {
	f.Var(newUint16Value(value, p), name, usage, opts...)
}

// Uint16Var defines an uint16 flag with specified name, default value, and usage string.
// The argument p points to an uint16 variable in which to store the value of the flag.
func Uint16Var(p *uint16, name string, value uint16, usage string, opts ...Opt) {
	CommandLine.Uint16Var(p, name, value, usage, opts...)
}

// Uint16 defines an uint16 flag with specified name, default value, and usage string.
// The return value is the address of an uint16 variable that stores the value of the flag.
func (f *FlagSet) Uint16(name string, value uint16, usage string, opts ...Opt) *uint16 {
	var p uint16
	f.Uint16Var(&p, name, value, usage, opts...)
	return &p
}

// Uint16 defines an uint16 flag with specified name, default value, and usage string.
// The return value is the address of an uint16 variable that stores the value of the flag.
func Uint16(name string, value uint16, usage string, opts ...Opt) *uint16 {
	return CommandLine.Uint16(name, value, usage, opts...)
}
