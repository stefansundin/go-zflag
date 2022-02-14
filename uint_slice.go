// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"fmt"
	"strconv"
	"strings"
)

// -- uintSlice Value
type uintSliceValue struct {
	value   *[]uint
	changed bool
}

func newUintSliceValue(val []uint, p *[]uint) *uintSliceValue {
	uisv := new(uintSliceValue)
	uisv.value = p
	*uisv.value = val
	return uisv
}

func (s *uintSliceValue) Set(val string) error {
	ss := strings.Split(val, ",")
	out := make([]uint, len(ss))
	for i, d := range ss {
		u, err := strconv.ParseUint(d, 10, 0)
		if err != nil {
			return err
		}
		out[i] = uint(u)
	}
	if !s.changed {
		*s.value = out
	} else {
		*s.value = append(*s.value, out...)
	}
	s.changed = true
	return nil
}

func (s *uintSliceValue) Get() interface{} {
	return *s.value
}

func (s *uintSliceValue) Type() string {
	return "uintSlice"
}

func (s *uintSliceValue) String() string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = fmt.Sprintf("%d", d)
	}
	return "[" + strings.Join(out, ",") + "]"
}

func (s *uintSliceValue) fromString(val string) (uint, error) {
	t, err := strconv.ParseUint(val, 10, 0)
	if err != nil {
		return 0, err
	}
	return uint(t), nil
}

func (s *uintSliceValue) toString(val uint) string {
	return fmt.Sprintf("%d", val)
}

func (s *uintSliceValue) Append(val string) error {
	i, err := s.fromString(val)
	if err != nil {
		return err
	}
	*s.value = append(*s.value, i)
	return nil
}

func (s *uintSliceValue) Replace(val []string) error {
	out := make([]uint, len(val))
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

func (s *uintSliceValue) GetSlice() []string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = s.toString(d)
	}
	return out
}

// GetUintSlice returns the []uint value of a flag with the given name.
func (f *FlagSet) GetUintSlice(name string) ([]uint, error) {
	val, err := f.getFlagType(name, "uintSlice")
	if err != nil {
		return []uint{}, err
	}
	return val.([]uint), nil
}

// MustGetUintSlice is like GetUintSlice, but panics on error.
func (f *FlagSet) MustGetUintSlice(name string) []uint {
	val, err := f.GetUintSlice(name)
	if err != nil {
		panic(err)
	}
	return val
}

// UintSliceVar defines a []uint flag with specified name, default value, and usage string.
// The argument p points to a []uint variable in which to store the value of the flag.
func (f *FlagSet) UintSliceVar(p *[]uint, name string, value []uint, usage string, opts ...Opt) {
	f.Var(newUintSliceValue(value, p), name, usage, opts...)
}

// UintSliceVar defines a []uint flag with specified name, default value, and usage string.
// The argument p points to a []uint variable in which to store the value of the flag.
func UintSliceVar(p *[]uint, name string, value []uint, usage string, opts ...Opt) {
	CommandLine.UintSliceVar(p, name, value, usage, opts...)
}

// UintSlice defines a []uint flag with specified name, default value, and usage string.
// The return value is the address of a []uint variable that stores the value of the flag.
func (f *FlagSet) UintSlice(name string, value []uint, usage string, opts ...Opt) *[]uint {
	var p []uint
	f.UintSliceVar(&p, name, value, usage, opts...)
	return &p
}

// UintSlice defines a []uint flag with specified name, default value, and usage string.
// The return value is the address of a []uint variable that stores the value of the flag.
func UintSlice(name string, value []uint, usage string, opts ...Opt) *[]uint {
	return CommandLine.UintSlice(name, value, usage, opts...)
}
