// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

// -- string Value
type stringValue string

func newStringValue(val string, p *string) *stringValue {
	*p = val
	return (*stringValue)(p)
}

func (s *stringValue) Set(val string) error {
	*s = stringValue(val)
	return nil
}
func (s *stringValue) Get() interface{} {
	return string(*s)
}
func (s *stringValue) Type() string {
	return "string"
}

func (s *stringValue) String() string { return string(*s) }

// GetString return the string value of a flag with the given name
func (f *FlagSet) GetString(name string) (string, error) {
	val, err := f.getFlagType(name, "string")
	if err != nil {
		return "", err
	}
	return val.(string), nil
}

// MustGetString is like GetString, but panics on error.
func (f *FlagSet) MustGetString(name string) string {
	val, err := f.GetString(name)
	if err != nil {
		panic(err)
	}
	return val
}

// StringVar defines a string flag with specified name, default value, and usage string.
// The argument p points to a string variable in which to store the value of the flag.
func (f *FlagSet) StringVar(p *string, name string, value string, usage string, opts ...Opt) {
	f.Var(newStringValue(value, p), name, usage, opts...)
}

// StringVar defines a string flag with specified name, default value, and usage string.
// The argument p points to a string variable in which to store the value of the flag.
func StringVar(p *string, name string, value string, usage string, opts ...Opt) {
	CommandLine.StringVar(p, name, value, usage, opts...)
}

// String defines a string flag with specified name, default value, and usage string.
// The return value is the address of a string variable that stores the value of the flag.
func (f *FlagSet) String(name string, value string, usage string, opts ...Opt) *string {
	var p string
	f.StringVar(&p, name, value, usage, opts...)
	return &p
}

// String defines a string flag with specified name, default value, and usage string.
// The return value is the address of a string variable that stores the value of the flag.
func String(name string, value string, usage string, opts ...Opt) *string {
	return CommandLine.String(name, value, usage, opts...)
}
