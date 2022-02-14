// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"bytes"
	"io"
	"strconv"
)

// -- stringToInt64 Value
type stringToInt64Value struct {
	value   *map[string]int64
	changed bool
}

func newStringToInt64Value(val map[string]int64, p *map[string]int64) *stringToInt64Value {
	ssv := new(stringToInt64Value)
	ssv.value = p
	*ssv.value = val
	return ssv
}

// Format: a=1,b=2
func (s *stringToInt64Value) Set(val string) error {
	// read flag arguments with CSV parser
	mapStrInt, err := readCSVKeyValue(val)
	if err != nil && err != io.EOF {
		return err
	}

	out := make(map[string]int64, len(mapStrInt))
	for key, value := range mapStrInt {
		var err error
		out[key], err = strconv.ParseInt(value, 10, 64)
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

func (s *stringToInt64Value) Get() interface{} {
	return *s.value
}

func (s *stringToInt64Value) Type() string {
	return "stringToInt64"
}

func (s *stringToInt64Value) String() string {
	var buf bytes.Buffer
	i := 0
	for k, v := range *s.value {
		if i > 0 {
			buf.WriteRune(',')
		}
		buf.WriteString(k)
		buf.WriteRune('=')
		buf.WriteString(strconv.FormatInt(v, 10))
		i++
	}
	return "[" + buf.String() + "]"
}

// GetStringToInt64 return the map[string]int64 value of a flag with the given name
func (f *FlagSet) GetStringToInt64(name string) (map[string]int64, error) {
	val, err := f.getFlagType(name, "stringToInt64")
	if err != nil {
		return map[string]int64{}, err
	}
	return val.(map[string]int64), nil
}

// MustGetStringToInt64 is like GetStringToInt64, but panics on error.
func (f *FlagSet) MustGetStringToInt64(name string) map[string]int64 {
	val, err := f.GetStringToInt64(name)
	if err != nil {
		panic(err)
	}
	return val
}

// StringToInt64Var defines a map[string]int64 flag with specified name, default value, and usage string.
// The argument p points to a map[string]int64 variable in which to store the values of multiple flags.
// The values will be separated on comma. Items can be quoted, or escape commas to avoid splitting.
func (f *FlagSet) StringToInt64Var(p *map[string]int64, name string, value map[string]int64, usage string, opts ...Opt) {
	f.Var(newStringToInt64Value(value, p), name, usage, opts...)
}

// StringToInt64Var defines a map[string]int64 flag with specified name, default value, and usage string.
// The argument p points to a map[string]int64 variable in which to store the values of multiple flags.
// The values will be separated on comma. Items can be quoted, or escape commas to avoid splitting.
func StringToInt64Var(p *map[string]int64, name string, value map[string]int64, usage string, opts ...Opt) {
	CommandLine.StringToInt64Var(p, name, value, usage, opts...)
}

// StringToInt64 defines a map[string]int64 flag with specified name, default value, and usage string.
// The return value is the address of a map[string]int64 variable that stores the values of multiple flags.
// The values will be separated on comma. Items can be quoted, or escape commas to avoid splitting.
func (f *FlagSet) StringToInt64(name string, value map[string]int64, usage string, opts ...Opt) *map[string]int64 {
	var p map[string]int64
	f.StringToInt64Var(&p, name, value, usage, opts...)
	return &p
}

// StringToInt64 defines a map[string]int64 flag with specified name, default value, and usage string.
// The return value is the address of a map[string]int64 variable that stores the values of multiple flags.
// The values will be separated on comma. Items can be quoted, or escape commas to avoid splitting.
func StringToInt64(name string, value map[string]int64, usage string, opts ...Opt) *map[string]int64 {
	return CommandLine.StringToInt64(name, value, usage, opts...)
}
