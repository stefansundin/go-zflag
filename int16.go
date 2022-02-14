// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import "strconv"

// -- int16 Value
type int16Value int16

func newInt16Value(val int16, p *int16) *int16Value {
	*p = val
	return (*int16Value)(p)
}

func (i *int16Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 16)
	*i = int16Value(v)
	return err
}

func (i *int16Value) Get() interface{} {
	return int16(*i)
}

func (i *int16Value) Type() string {
	return "int16"
}

func (i *int16Value) String() string { return strconv.FormatInt(int64(*i), 10) }

// GetInt16 returns the int16 value of a flag with the given name
func (f *FlagSet) GetInt16(name string) (int16, error) {
	val, err := f.getFlagType(name, "int16")
	if err != nil {
		return 0, err
	}
	return val.(int16), nil
}

// MustGetInt16 is like GetInt16, but panics on error.
func (f *FlagSet) MustGetInt16(name string) int16 {
	val, err := f.GetInt16(name)
	if err != nil {
		panic(err)
	}
	return val
}

// Int16Var defines an int16 flag with specified name, default value, and usage string.
// The argument p points to an int16 variable in which to store the value of the flag.
func (f *FlagSet) Int16Var(p *int16, name string, value int16, usage string, opts ...Opt) {
	f.Var(newInt16Value(value, p), name, usage, opts...)
}

// Int16Var defines an int16 flag with specified name, default value, and usage string.
// The argument p points to an int16 variable in which to store the value of the flag.
func Int16Var(p *int16, name string, value int16, usage string, opts ...Opt) {
	CommandLine.Int16Var(p, name, value, usage, opts...)
}

// Int16 defines an int16 flag with specified name, default value, and usage string.
// The return value is the address of an int16 variable that stores the value of the flag.
func (f *FlagSet) Int16(name string, value int16, usage string, opts ...Opt) *int16 {
	var p int16
	f.Int16Var(&p, name, value, usage, opts...)
	return &p
}

// Int16 defines an int16 flag with specified name, default value, and usage string.
// The return value is the address of an int16 variable that stores the value of the flag.
func Int16(name string, value int16, usage string, opts ...Opt) *int16 {
	return CommandLine.Int16(name, value, usage, opts...)
}
