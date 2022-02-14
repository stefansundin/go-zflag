// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

// -- stringArray Value
type stringArrayValue struct {
	value   *[]string
	changed bool
}

func newStringArrayValue(val []string, p *[]string) *stringArrayValue {
	ssv := new(stringArrayValue)
	ssv.value = p
	*ssv.value = val
	return ssv
}

func (s *stringArrayValue) Get() interface{} {
	return *s.value
}

func (s *stringArrayValue) Set(val string) error {
	if val == "" {
		return nil
	}
	if !s.changed {
		*s.value = []string{val}
		s.changed = true
	} else {
		*s.value = append(*s.value, val)
	}
	return nil
}

func (s *stringArrayValue) Append(val string) error {
	*s.value = append(*s.value, val)
	return nil
}

func (s *stringArrayValue) Replace(val []string) error {
	out := make([]string, len(val))
	copy(out, val)
	*s.value = out
	return nil
}

func (s *stringArrayValue) GetSlice() []string {
	out := make([]string, len(*s.value))
	copy(out, *s.value)
	return out
}

func (s *stringArrayValue) Type() string {
	return "stringArray"
}

func (s *stringArrayValue) String() string {
	str, _ := writeAsCSV(*s.value)
	return "[" + str + "]"
}

// GetStringArray return the []string value of a flag with the given name
func (f *FlagSet) GetStringArray(name string) ([]string, error) {
	val, err := f.getFlagType(name, "stringArray")
	if err != nil {
		return []string{}, err
	}
	return val.([]string), nil
}

// MustGetStringArray is like GetStringArray, but panics on error.
func (f *FlagSet) MustGetStringArray(name string) []string {
	val, err := f.GetStringArray(name)
	if err != nil {
		panic(err)
	}
	return val
}

// StringArrayVar defines a []string flag with specified name, default value, and usage string.
// The argument p points to a []string variable in which to store the value of the flag.
// Compared to StringSlice flags, a StringArray will not be split on commas.
func (f *FlagSet) StringArrayVar(p *[]string, name string, value []string, usage string, opts ...Opt) {
	f.Var(newStringArrayValue(value, p), name, usage, opts...)
}

// StringArrayVar defines a []string flag with specified name, default value, and usage string.
// The argument p points to a []string variable in which to store the value of the flag.
// Compared to StringSlice flags, a StringArray will not be split on commas.
func StringArrayVar(p *[]string, name string, value []string, usage string, opts ...Opt) {
	CommandLine.StringArrayVar(p, name, value, usage, opts...)
}

// StringArray defines a []string flag with specified name, default value, and usage string.
// The return value is the address of a []string variable that stores the value of the flag.
// Compared to StringSlice flags, a StringArray will not be split on commas.
func (f *FlagSet) StringArray(name string, value []string, usage string, opts ...Opt) *[]string {
	var p []string
	f.StringArrayVar(&p, name, value, usage, opts...)
	return &p
}

// StringArray defines a []string flag with specified name, default value, and usage string.
// The return value is the address of a []string variable that stores the value of the flag.
// Compared to StringSlice flags, a StringArray will not be split on commas.
func StringArray(name string, value []string, usage string, opts ...Opt) *[]string {
	return CommandLine.StringArray(name, value, usage, opts...)
}
