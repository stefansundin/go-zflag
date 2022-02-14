// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"fmt"
	"strconv"
	"strings"
)

// -- int64Slice Value
type int64SliceValue struct {
	value   *[]int64
	changed bool
}

func newInt64SliceValue(val []int64, p *[]int64) *int64SliceValue {
	isv := new(int64SliceValue)
	isv.value = p
	*isv.value = val
	return isv
}

func (s *int64SliceValue) Set(val string) error {
	ss := strings.Split(val, ",")
	out := make([]int64, len(ss))
	for i, d := range ss {
		var err error
		out[i], err = strconv.ParseInt(d, 0, 64)
		if err != nil {
			return err
		}

	}
	if !s.changed {
		*s.value = out
	} else {
		*s.value = append(*s.value, out...)
	}
	s.changed = true
	return nil
}

func (s *int64SliceValue) Get() interface{} {
	return *s.value
}

func (s *int64SliceValue) Type() string {
	return "int64Slice"
}

func (s *int64SliceValue) String() string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = fmt.Sprintf("%d", d)
	}
	return "[" + strings.Join(out, ",") + "]"
}

func (s *int64SliceValue) fromString(val string) (int64, error) {
	return strconv.ParseInt(val, 0, 64)
}

func (s *int64SliceValue) toString(val int64) string {
	return fmt.Sprintf("%d", val)
}

func (s *int64SliceValue) Append(val string) error {
	i, err := s.fromString(val)
	if err != nil {
		return err
	}
	*s.value = append(*s.value, i)
	return nil
}

func (s *int64SliceValue) Replace(val []string) error {
	out := make([]int64, len(val))
	for i, d := range val {
		var err error
		out[i], err = s.fromString(d)
		if err != nil {
			return err
		}
	}
	*s.value = out
	return nil
}

func (s *int64SliceValue) GetSlice() []string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = s.toString(d)
	}
	return out
}

// GetInt64Slice return the []int64 value of a flag with the given name
func (f *FlagSet) GetInt64Slice(name string) ([]int64, error) {
	val, err := f.getFlagType(name, "int64Slice")
	if err != nil {
		return []int64{}, err
	}
	return val.([]int64), nil
}

// MustGetInt64Slice is like GetInt64Slice, but panics on error.
func (f *FlagSet) MustGetInt64Slice(name string) []int64 {
	val, err := f.GetInt64Slice(name)
	if err != nil {
		panic(err)
	}
	return val
}

// Int64SliceVar defines a []int64 flag with specified name, default value, and usage string.
// The argument p points to a []int64 variable in which to store the value of the flag.
func (f *FlagSet) Int64SliceVar(p *[]int64, name string, value []int64, usage string, opts ...Opt) {
	f.Var(newInt64SliceValue(value, p), name, usage, opts...)
}

// Int64SliceVar defines a []int64 flag with specified name, default value, and usage string.
// The argument p points to a []int64 variable in which to store the value of the flag.
func Int64SliceVar(p *[]int64, name string, value []int64, usage string, opts ...Opt) {
	CommandLine.Int64SliceVar(p, name, value, usage, opts...)
}

// Int64Slice defines a []int64 flag with specified name, default value, and usage string.
// The return value is the address of a []int64 variable that stores the value of the flag.
func (f *FlagSet) Int64Slice(name string, value []int64, usage string, opts ...Opt) *[]int64 {
	var p []int64
	f.Int64SliceVar(&p, name, value, usage, opts...)
	return &p
}

// Int64Slice defines a []int64 flag with specified name, default value, and usage string.
// The return value is the address of a []int64 variable that stores the value of the flag.
func Int64Slice(name string, value []int64, usage string, opts ...Opt) *[]int64 {
	return CommandLine.Int64Slice(name, value, usage, opts...)
}
