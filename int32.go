// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import "strconv"

// -- int32 Value
type int32Value int32

func newInt32Value(val int32, p *int32) *int32Value {
	*p = val
	return (*int32Value)(p)
}

func (i *int32Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 32)
	*i = int32Value(v)
	return err
}

func (i *int32Value) Get() interface{} {
	return int32(*i)
}

func (i *int32Value) Type() string {
	return "int32"
}

func (i *int32Value) String() string { return strconv.FormatInt(int64(*i), 10) }

// GetInt32 return the int32 value of a flag with the given name
func (f *FlagSet) GetInt32(name string) (int32, error) {
	val, err := f.getFlagType(name, "int32")
	if err != nil {
		return 0, err
	}
	return val.(int32), nil
}

// MustGetInt32 is like GetInt32, but panics on error.
func (f *FlagSet) MustGetInt32(name string) int32 {
	val, err := f.GetInt32(name)
	if err != nil {
		panic(err)
	}
	return val
}

// Int32Var defines an int32 flag with specified name, default value, and usage string.
// The argument p points to an int32 variable in which to store the value of the flag.
func (f *FlagSet) Int32Var(p *int32, name string, value int32, usage string, opts ...Opt) {
	f.Var(newInt32Value(value, p), name, usage, opts...)
}

// Int32Var defines an int32 flag with specified name, default value, and usage string.
// The argument p points to an int32 variable in which to store the value of the flag.
func Int32Var(p *int32, name string, value int32, usage string, opts ...Opt) {
	CommandLine.Int32Var(p, name, value, usage, opts...)
}

// Int32 defines an int32 flag with specified name, default value, and usage string.
// The return value is the address of an int32 variable that stores the value of the flag.
func (f *FlagSet) Int32(name string, value int32, usage string, opts ...Opt) *int32 {
	var p int32
	f.Int32Var(&p, name, value, usage, opts...)
	return &p
}

// Int32 defines an int32 flag with specified name, default value, and usage string.
// The return value is the address of an int32 variable that stores the value of the flag.
func Int32(name string, value int32, usage string, opts ...Opt) *int32 {
	return CommandLine.Int32(name, value, usage, opts...)
}
