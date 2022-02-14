// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"strings"
	"time"
)

// -- durationSlice Value
type durationSliceValue struct {
	value   *[]time.Duration
	changed bool
}

func newDurationSliceValue(val []time.Duration, p *[]time.Duration) *durationSliceValue {
	dsv := new(durationSliceValue)
	dsv.value = p
	*dsv.value = val
	return dsv
}

func (s *durationSliceValue) Set(val string) error {
	ss := strings.Split(val, ",")
	out := make([]time.Duration, len(ss))
	for i, d := range ss {
		var err error
		out[i], err = time.ParseDuration(d)
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

func (s *durationSliceValue) Get() interface{} {
	return *s.value
}

func (s *durationSliceValue) Type() string {
	return "durationSlice"
}

func (s *durationSliceValue) String() string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = d.String()
	}
	return "[" + strings.Join(out, ",") + "]"
}

func (s *durationSliceValue) fromString(val string) (time.Duration, error) {
	return time.ParseDuration(val)
}

func (s *durationSliceValue) toString(val time.Duration) string {
	return val.String()
}

func (s *durationSliceValue) Append(val string) error {
	i, err := s.fromString(val)
	if err != nil {
		return err
	}
	*s.value = append(*s.value, i)
	return nil
}

func (s *durationSliceValue) Replace(val []string) error {
	out := make([]time.Duration, len(val))
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

func (s *durationSliceValue) GetSlice() []string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = s.toString(d)
	}
	return out
}

// GetDurationSlice returns the []time.Duration value of a flag with the given name
func (f *FlagSet) GetDurationSlice(name string) ([]time.Duration, error) {
	val, err := f.getFlagType(name, "durationSlice")
	if err != nil {
		return []time.Duration{}, err
	}
	return val.([]time.Duration), nil
}

// MustGetDurationSlice is like GetDurationSlice, but panics on error.
func (f *FlagSet) MustGetDurationSlice(name string) []time.Duration {
	val, err := f.GetDurationSlice(name)
	if err != nil {
		panic(err)
	}
	return val
}

// DurationSliceVar defines a durationSlice flag with specified name, default value, and usage string.
// The argument p points to a []time.Duration variable in which to store the value of the flag.
func (f *FlagSet) DurationSliceVar(p *[]time.Duration, name string, value []time.Duration, usage string, opts ...Opt) {
	f.Var(newDurationSliceValue(value, p), name, usage, opts...)
}

// DurationSliceVar defines a duration[] flag with specified name, default value, and usage string.
// The argument p points to a duration[] variable in which to store the value of the flag.
func DurationSliceVar(p *[]time.Duration, name string, value []time.Duration, usage string, opts ...Opt) {
	CommandLine.DurationSliceVar(p, name, value, usage, opts...)
}

// DurationSlice defines a []time.Duration flag with specified name, default value, and usage string.
// The return value is the address of a []time.Duration variable that stores the value of the flag.
func (f *FlagSet) DurationSlice(name string, value []time.Duration, usage string, opts ...Opt) *[]time.Duration {
	var p []time.Duration
	f.DurationSliceVar(&p, name, value, usage, opts...)
	return &p
}

// DurationSlice defines a []time.Duration flag with specified name, default value, and usage string.
// The return value is the address of a []time.Duration variable that stores the value of the flag.
func DurationSlice(name string, value []time.Duration, usage string, opts ...Opt) *[]time.Duration {
	return CommandLine.DurationSlice(name, value, usage, opts...)
}
