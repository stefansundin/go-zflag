// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"fmt"
	"strconv"
	"strings"
)

// -- uint64Slice Value
type uint64SliceValue struct {
	value   *[]uint64
	changed bool
}

func newUint64SliceValue(val []uint64, p *[]uint64) *uint64SliceValue {
	isv := new(uint64SliceValue)
	isv.value = p
	*isv.value = val
	return isv
}

func (s *uint64SliceValue) Get() interface{} {
	return *s.value
}

func (s *uint64SliceValue) Set(val string) error {
	ss := strings.Split(val, ",")
	out := make([]uint64, len(ss))
	for i, d := range ss {
		var err error
		out[i], err = strconv.ParseUint(d, 0, 64)
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

func (s *uint64SliceValue) Type() string {
	return "uint64Slice"
}

func (s *uint64SliceValue) String() string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = fmt.Sprintf("%d", d)
	}
	return "[" + strings.Join(out, ",") + "]"
}

func (s *uint64SliceValue) fromString(val string) (uint64, error) {
	return strconv.ParseUint(val, 0, 64)
}

func (s *uint64SliceValue) toString(val uint64) string {
	return fmt.Sprintf("%d", val)
}

func (s *uint64SliceValue) Append(val string) error {
	i, err := s.fromString(val)
	if err != nil {
		return err
	}
	*s.value = append(*s.value, i)
	return nil
}

func (s *uint64SliceValue) Replace(val []string) error {
	out := make([]uint64, len(val))
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

func (s *uint64SliceValue) GetSlice() []string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = s.toString(d)
	}
	return out
}

// GetUint64Slice return the []uint64 value of a flag with the given name
func (f *FlagSet) GetUint64Slice(name string) ([]uint64, error) {
	val, err := f.getFlagType(name, "uint64Slice")
	if err != nil {
		return []uint64{}, err
	}
	return val.([]uint64), nil
}

// MustGetUint64Slice is like GetUint64Slice, but panics on error.
func (f *FlagSet) MustGetUint64Slice(name string) []uint64 {
	val, err := f.GetUint64Slice(name)
	if err != nil {
		panic(err)
	}
	return val
}

// Uint64SliceVar defines a []uint64 flag with specified name, default value, and usage string.
// The argument p points to a []uint64 variable in which to store the value of the flag.
func (f *FlagSet) Uint64SliceVar(p *[]uint64, name string, value []uint64, usage string, opts ...Opt) {
	f.Var(newUint64SliceValue(value, p), name, usage, opts...)
}

// Uint64SliceVar defines a []uint64 flag with specified name, default value, and usage string.
// The argument p points to a []uint64 variable in which to store the value of the flag.
func Uint64SliceVar(p *[]uint64, name string, value []uint64, usage string, opts ...Opt) {
	CommandLine.Uint64SliceVar(p, name, value, usage, opts...)
}

// Uint64Slice defines a []uint64 flag with specified name, default value, and usage string.
// The return value is the address of a []uint64 variable that stores the value of the flag.
func (f *FlagSet) Uint64Slice(name string, value []uint64, usage string, opts ...Opt) *[]uint64 {
	var p []uint64
	f.Uint64SliceVar(&p, name, value, usage, opts...)
	return &p
}

// Uint64Slice defines a []uint64 flag with specified name, default value, and usage string.
// The return value is the address of a []uint64 variable that stores the value of the flag.
func Uint64Slice(name string, value []uint64, usage string, opts ...Opt) *[]uint64 {
	return CommandLine.Uint64Slice(name, value, usage, opts...)
}
