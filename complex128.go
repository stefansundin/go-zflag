// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.15
// +build go1.15

package zflag

import "strconv"

// -- complex128 Value
type complex128Value complex128

func newComplex128Value(val complex128, p *complex128) *complex128Value {
	*p = val
	return (*complex128Value)(p)
}

func (f *complex128Value) Get() interface{} {
	return complex128(*f)
}

func (f *complex128Value) Set(s string) error {
	v, err := strconv.ParseComplex(s, 128)
	*f = complex128Value(v)
	return err
}

func (f *complex128Value) Type() string {
	return "complex128"
}

func (f *complex128Value) String() string { return strconv.FormatComplex(complex128(*f), 'g', -1, 128) }

// GetComplex128 return the complex128 value of a flag with the given name
func (f *FlagSet) GetComplex128(name string) (complex128, error) {
	val, err := f.getFlagType(name, "complex128")
	if err != nil {
		return 0, err
	}
	return val.(complex128), nil
}

// MustGetComplex128 is like GetComplex128, but panics on error.
func (f *FlagSet) MustGetComplex128(name string) complex128 {
	val, err := f.GetComplex128(name)
	if err != nil {
		panic(err)
	}
	return val
}

// Complex128Var defines a complex128 flag with specified name, default value, and usage string.
// The argument p points to a complex128 variable in which to store the value of the flag.
func (f *FlagSet) Complex128Var(p *complex128, name string, value complex128, usage string, opts ...Opt) {
	f.Var(newComplex128Value(value, p), name, usage, opts...)
}

// Complex128Var defines a complex128 flag with specified name, default value, and usage string.
// The argument p points to a complex128 variable in which to store the value of the flag.
func Complex128Var(p *complex128, name string, value complex128, usage string, opts ...Opt) {
	CommandLine.Complex128Var(p, name, value, usage, opts...)
}

// Complex128 defines a complex128 flag with specified name, default value, and usage string.
// The return value is the address of a complex128 variable that stores the value of the flag.
func (f *FlagSet) Complex128(name string, value complex128, usage string, opts ...Opt) *complex128 {
	var p complex128
	f.Complex128Var(&p, name, value, usage, opts...)
	return &p
}

// Complex128 defines a complex128 flag with specified name, default value, and usage string.
// The return value is the address of a complex128 variable that stores the value of the flag.
func Complex128(name string, value complex128, usage string, opts ...Opt) *complex128 {
	return CommandLine.Complex128(name, value, usage, opts...)
}
