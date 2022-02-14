// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import "strconv"

// -- uint64 Value
type uint64Value uint64

func newUint64Value(val uint64, p *uint64) *uint64Value {
	*p = val
	return (*uint64Value)(p)
}

func (i *uint64Value) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 64)
	*i = uint64Value(v)
	return err
}

func (i *uint64Value) Get() interface{} {
	return uint64(*i)
}

func (i *uint64Value) Type() string {
	return "uint64"
}

func (i *uint64Value) String() string { return strconv.FormatUint(uint64(*i), 10) }

// GetUint64 return the uint64 value of a flag with the given name
func (f *FlagSet) GetUint64(name string) (uint64, error) {
	val, err := f.getFlagType(name, "uint64")
	if err != nil {
		return 0, err
	}
	return val.(uint64), nil
}

// MustGetUint64 is like GetUint64, but panics on error.
func (f *FlagSet) MustGetUint64(name string) uint64 {
	val, err := f.GetUint64(name)
	if err != nil {
		panic(err)
	}
	return val
}

// Uint64Var defines an uint64 flag with specified name, default value, and usage string.
// The argument p points to an uint64 variable in which to store the value of the flag.
func (f *FlagSet) Uint64Var(p *uint64, name string, value uint64, usage string, opts ...Opt) {
	f.Var(newUint64Value(value, p), name, usage, opts...)
}

// Uint64Var defines an uint64 flag with specified name, default value, and usage string.
// The argument p points to an uint64 variable in which to store the value of the flag.
func Uint64Var(p *uint64, name string, value uint64, usage string, opts ...Opt) {
	CommandLine.Uint64Var(p, name, value, usage, opts...)
}

// Uint64 defines an uint64 flag with specified name, default value, and usage string.
// The return value is the address of an uint64 variable that stores the value of the flag.
func (f *FlagSet) Uint64(name string, value uint64, usage string, opts ...Opt) *uint64 {
	var p uint64
	f.Uint64Var(&p, name, value, usage, opts...)
	return &p
}

// Uint64 defines an uint64 flag with specified name, default value, and usage string.
// The return value is the address of an uint64 variable that stores the value of the flag.
func Uint64(name string, value uint64, usage string, opts ...Opt) *uint64 {
	return CommandLine.Uint64(name, value, usage, opts...)
}
