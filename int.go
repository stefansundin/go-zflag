// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import "strconv"

// -- int Value
type intValue int

func newIntValue(val int, p *int) *intValue {
	*p = val
	return (*intValue)(p)
}

func (i *intValue) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	*i = intValue(v)
	return err
}

func (i *intValue) Get() interface{} {
	return int(*i)
}

func (i *intValue) Type() string {
	return "int"
}

func (i *intValue) String() string { return strconv.Itoa(int(*i)) }

// GetInt return the int value of a flag with the given name
func (f *FlagSet) GetInt(name string) (int, error) {
	val, err := f.getFlagType(name, "int")
	if err != nil {
		return 0, err
	}
	return val.(int), nil
}

// MustGetInt is like GetInt, but panics on error.
func (f *FlagSet) MustGetInt(name string) int {
	val, err := f.GetInt(name)
	if err != nil {
		panic(err)
	}
	return val
}

// IntVar defines an int flag with specified name, default value, and usage string.
// The argument p points to an int variable in which to store the value of the flag.
func (f *FlagSet) IntVar(p *int, name string, value int, usage string, opts ...Opt) {
	f.Var(newIntValue(value, p), name, usage, opts...)
}

// IntVar defines an int flag with specified name, default value, and usage string.
// The argument p points to an int variable in which to store the value of the flag.
func IntVar(p *int, name string, value int, usage string, opts ...Opt) {
	CommandLine.IntVar(p, name, value, usage, opts...)
}

// Int defines an int flag with specified name, default value, and usage string.
// The return value is the address of an int variable that stores the value of the flag.
func (f *FlagSet) Int(name string, value int, usage string, opts ...Opt) *int {
	var p int
	f.IntVar(&p, name, value, usage, opts...)
	return &p
}

// Int defines an int flag with specified name, default value, and usage string.
// The return value is the address of an int variable that stores the value of the flag.
func Int(name string, value int, usage string, opts ...Opt) *int {
	return CommandLine.Int(name, value, usage, opts...)
}
