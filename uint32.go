// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import "strconv"

// -- uint32 value
type uint32Value uint32

func newUint32Value(val uint32, p *uint32) *uint32Value {
	*p = val
	return (*uint32Value)(p)
}

func (i *uint32Value) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 32)
	*i = uint32Value(v)
	return err
}

func (i *uint32Value) Get() interface{} {
	return uint32(*i)
}

func (i *uint32Value) Type() string {
	return "uint32"
}

func (i *uint32Value) String() string { return strconv.FormatUint(uint64(*i), 10) }

// GetUint32 return the uint32 value of a flag with the given name
func (f *FlagSet) GetUint32(name string) (uint32, error) {
	val, err := f.getFlagType(name, "uint32")
	if err != nil {
		return 0, err
	}
	return val.(uint32), nil
}

// MustGetUint32 is like GetUint32, but panics on error.
func (f *FlagSet) MustGetUint32(name string) uint32 {
	val, err := f.GetUint32(name)
	if err != nil {
		panic(err)
	}
	return val
}

// Uint32Var defines an uint32 flag with specified name, default value, and usage string.
// The argument p points to an uint32 variable in which to store the value of the flag.
func (f *FlagSet) Uint32Var(p *uint32, name string, value uint32, usage string, opts ...Opt) {
	f.Var(newUint32Value(value, p), name, usage, opts...)
}

// Uint32Var defines an uint32 flag with specified name, default value, and usage string.
// The argument p points to an uint32 variable in which to store the value of the flag.
func Uint32Var(p *uint32, name string, value uint32, usage string, opts ...Opt) {
	CommandLine.Uint32Var(p, name, value, usage, opts...)
}

// Uint32 defines an uint32 flag with specified name, default value, and usage string.
// The return value is the address of an uint32 variable that stores the value of the flag.
func (f *FlagSet) Uint32(name string, value uint32, usage string, opts ...Opt) *uint32 {
	var p uint32
	f.Uint32Var(&p, name, value, usage, opts...)
	return &p
}

// Uint32 defines an uint32 flag with specified name, default value, and usage string.
// The return value is the address of an uint32 variable that stores the value of the flag.
func Uint32(name string, value uint32, usage string, opts ...Opt) *uint32 {
	return CommandLine.Uint32(name, value, usage, opts...)
}
