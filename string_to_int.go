// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"bytes"
	"io"
	"strconv"
)

// -- stringToInt Value
type stringToIntValue struct {
	value   *map[string]int
	changed bool
}

func newStringToIntValue(val map[string]int, p *map[string]int) *stringToIntValue {
	ssv := new(stringToIntValue)
	ssv.value = p
	*ssv.value = val
	return ssv
}

// Format: a=1,b=2
func (s *stringToIntValue) Set(val string) error {
	// read flag arguments with CSV parser
	mapStrInt, err := readCSVKeyValue(val)
	if err != nil && err != io.EOF {
		return err
	}

	out := make(map[string]int, len(mapStrInt))
	for key, value := range mapStrInt {
		var err error
		out[key], err = strconv.Atoi(value)
		if err != nil {
			return err
		}
	}

	if !s.changed {
		*s.value = out
	} else {
		for k, v := range out {
			(*s.value)[k] = v
		}
	}

	s.changed = true
	return nil
}

func (s *stringToIntValue) Get() interface{} {
	return *s.value
}

func (s *stringToIntValue) Type() string {
	return "stringToInt"
}

func (s *stringToIntValue) String() string {
	var buf bytes.Buffer
	i := 0
	for k, v := range *s.value {
		if i > 0 {
			buf.WriteRune(',')
		}
		buf.WriteString(k)
		buf.WriteRune('=')
		buf.WriteString(strconv.Itoa(v))
		i++
	}
	return "[" + buf.String() + "]"
}

// GetStringToInt return the map[string]int value of a flag with the given name
func (f *FlagSet) GetStringToInt(name string) (map[string]int, error) {
	val, err := f.getFlagType(name, "stringToInt")
	if err != nil {
		return map[string]int{}, err
	}
	return val.(map[string]int), nil
}

// MustGetStringToInt is like GetStringToInt, but panics on error.
func (f *FlagSet) MustGetStringToInt(name string) map[string]int {
	val, err := f.GetStringToInt(name)
	if err != nil {
		panic(err)
	}
	return val
}

// StringToIntVar defines a map[string]int flag with specified name, default value, and usage string.
// The argument p points to a map[string]int variable in which to store the values of multiple flags.
// The values will be separated on comma. Items can be quoted, or escape commas to avoid splitting.
func (f *FlagSet) StringToIntVar(p *map[string]int, name string, value map[string]int, usage string, opts ...Opt) {
	f.Var(newStringToIntValue(value, p), name, usage, opts...)
}

// StringToIntVar defines a map[string]int flag with specified name, default value, and usage string.
// The argument p points to a map[string]int variable in which to store the values of multiple flags.
// The values will be separated on comma. Items can be quoted, or escape commas to avoid splitting.
func StringToIntVar(p *map[string]int, name string, value map[string]int, usage string, opts ...Opt) {
	CommandLine.StringToIntVar(p, name, value, usage, opts...)
}

// StringToInt defines a map[string]int flag with specified name, default value, and usage string.
// The return value is the address of a map[string]int variable that stores the values of multiple flags.
// The values will be separated on comma. Items can be quoted, or escape commas to avoid splitting.
func (f *FlagSet) StringToInt(name string, value map[string]int, usage string, opts ...Opt) *map[string]int {
	var p map[string]int
	f.StringToIntVar(&p, name, value, usage, opts...)
	return &p
}

// StringToInt defines a map[string]int flag with specified name, default value, and usage string.
// The return value is the address of a map[string]int variable that stores the values of multiple flags.
// The values will be separated on comma. Items can be quoted, or escape commas to avoid splitting.
func StringToInt(name string, value map[string]int, usage string, opts ...Opt) *map[string]int {
	return CommandLine.StringToInt(name, value, usage, opts...)
}
