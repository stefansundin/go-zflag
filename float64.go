// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import "strconv"

// -- float64 Value
type float64Value float64

func newFloat64Value(val float64, p *float64) *float64Value {
	*p = val
	return (*float64Value)(p)
}

func (f *float64Value) Set(s string) error {
	v, err := strconv.ParseFloat(s, 64)
	*f = float64Value(v)
	return err
}

func (s *float64Value) Get() interface{} {
	return float64(*s)
}

func (f *float64Value) Type() string {
	return "float64"
}

func (f *float64Value) String() string { return strconv.FormatFloat(float64(*f), 'g', -1, 64) }

// GetFloat64 return the float64 value of a flag with the given name
func (f *FlagSet) GetFloat64(name string) (float64, error) {
	val, err := f.getFlagType(name, "float64")
	if err != nil {
		return 0, err
	}
	return val.(float64), nil
}

// MustGetFloat64 is like GetFloat64, but panics on error.
func (f *FlagSet) MustGetFloat64(name string) float64 {
	val, err := f.GetFloat64(name)
	if err != nil {
		panic(err)
	}
	return val
}

// Float64Var defines a float64 flag with specified name, default value, and usage string.
// The argument p points to a float64 variable in which to store the value of the flag.
func (f *FlagSet) Float64Var(p *float64, name string, value float64, usage string, opts ...Opt) {
	f.Var(newFloat64Value(value, p), name, usage, opts...)
}

// Float64Var defines a float64 flag with specified name, default value, and usage string.
// The argument p points to a float64 variable in which to store the value of the flag.
func Float64Var(p *float64, name string, value float64, usage string, opts ...Opt) {
	CommandLine.Float64Var(p, name, value, usage, opts...)
}

// Float64 defines a float64 flag with specified name, default value, and usage string.
// The return value is the address of a float64 variable that stores the value of the flag.
func (f *FlagSet) Float64(name string, value float64, usage string, opts ...Opt) *float64 {
	var p float64
	f.Float64Var(&p, name, value, usage, opts...)
	return &p
}

// Float64 defines a float64 flag with specified name, default value, and usage string.
// The return value is the address of a float64 variable that stores the value of the flag.
func Float64(name string, value float64, usage string, opts ...Opt) *float64 {
	return CommandLine.Float64(name, value, usage, opts...)
}
