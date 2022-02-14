// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"time"
)

// -- time.Duration Value
type durationValue time.Duration

func newDurationValue(val time.Duration, p *time.Duration) *durationValue {
	*p = val
	return (*durationValue)(p)
}

func (d *durationValue) Set(s string) error {
	v, err := time.ParseDuration(s)
	*d = durationValue(v)
	return err
}

func (d *durationValue) Get() interface{} {
	return time.Duration(*d)
}

func (d *durationValue) Type() string {
	return "duration"
}

func (d *durationValue) String() string { return (*time.Duration)(d).String() }

// GetDuration return the duration value of a flag with the given name
func (f *FlagSet) GetDuration(name string) (time.Duration, error) {
	val, err := f.getFlagType(name, "duration")
	if err != nil {
		return 0, err
	}
	return val.(time.Duration), nil
}

// MustGetDuration is like GetDuration, but panics on error.
func (f *FlagSet) MustGetDuration(name string) time.Duration {
	val, err := f.GetDuration(name)
	if err != nil {
		panic(err)
	}
	return val
}

// DurationVar defines a time.Duration flag with specified name, default value, and usage string.
// The argument p points to a time.Duration variable in which to store the value of the flag.
func (f *FlagSet) DurationVar(p *time.Duration, name string, value time.Duration, usage string, opts ...Opt) {
	f.Var(newDurationValue(value, p), name, usage, opts...)
}

// DurationVar defines a time.Duration flag with specified name, default value, and usage string.
// The argument p points to a time.Duration variable in which to store the value of the flag.
func DurationVar(p *time.Duration, name string, value time.Duration, usage string, opts ...Opt) {
	CommandLine.DurationVar(p, name, value, usage, opts...)
}

// Duration defines a time.Duration flag with specified name, default value, and usage string.
// The return value is the address of a time.Duration variable that stores the value of the flag.
func (f *FlagSet) Duration(name string, value time.Duration, usage string, opts ...Opt) *time.Duration {
	var p time.Duration
	f.DurationVar(&p, name, value, usage, opts...)
	return &p
}

// Duration defines a time.Duration flag with specified name, default value, and usage string.
// The return value is the address of a time.Duration variable that stores the value of the flag.
func Duration(name string, value time.Duration, usage string, opts ...Opt) *time.Duration {
	return CommandLine.Duration(name, value, usage, opts...)
}
