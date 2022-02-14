// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import "strconv"

// -- uint8 Value
type uint8Value uint8

func newUint8Value(val uint8, p *uint8) *uint8Value {
	*p = val
	return (*uint8Value)(p)
}

func (i *uint8Value) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 8)
	*i = uint8Value(v)
	return err
}

func (i *uint8Value) Get() interface{} {
	return uint8(*i)
}

func (i *uint8Value) Type() string {
	return "uint8"
}

func (i *uint8Value) String() string { return strconv.FormatUint(uint64(*i), 10) }

// GetUint8 return the uint8 value of a flag with the given name
func (f *FlagSet) GetUint8(name string) (uint8, error) {
	val, err := f.getFlagType(name, "uint8")
	if err != nil {
		return 0, err
	}
	return val.(uint8), nil
}

// MustGetUint8 is like GetUint8, but panics on error.
func (f *FlagSet) MustGetUint8(name string) uint8 {
	val, err := f.GetUint8(name)
	if err != nil {
		panic(err)
	}
	return val
}

// Uint8Var defines an uint8 flag with specified name, default value, and usage string.
// The argument p points to an uint8 variable in which to store the value of the flag.
func (f *FlagSet) Uint8Var(p *uint8, name string, value uint8, usage string, opts ...Opt) {
	f.Var(newUint8Value(value, p), name, usage, opts...)
}

// Uint8Var defines an uint8 flag with specified name, default value, and usage string.
// The argument p points to an uint8 variable in which to store the value of the flag.
func Uint8Var(p *uint8, name string, value uint8, usage string, opts ...Opt) {
	CommandLine.Uint8Var(p, name, value, usage, opts...)
}

// Uint8 defines an uint8 flag with specified name, default value, and usage string.
// The return value is the address of an uint8 variable that stores the value of the flag.
func (f *FlagSet) Uint8(name string, value uint8, usage string, opts ...Opt) *uint8 {
	var p uint8
	f.Uint8Var(&p, name, value, usage, opts...)
	return &p
}

// Uint8 defines an uint8 flag with specified name, default value, and usage string.
// The return value is the address of an uint8 variable that stores the value of the flag.
func Uint8(name string, value uint8, usage string, opts ...Opt) *uint8 {
	return CommandLine.Uint8(name, value, usage, opts...)
}
