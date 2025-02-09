// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"fmt"
	"io"
	"net"
	"strings"
)

// -- ipNetSlice Value
type ipNetSliceValue struct {
	value   *[]net.IPNet
	changed bool
}

func newIPNetSliceValue(val []net.IPNet, p *[]net.IPNet) *ipNetSliceValue {
	ipnsv := new(ipNetSliceValue)
	ipnsv.value = p
	*ipnsv.value = val
	return ipnsv
}

func (s *ipNetSliceValue) Get() interface{} {
	return *s.value
}

// Set converts, and assigns, the comma-separated IPNet argument string representation as the []net.IPNet value of this flag.
// If Set is called on a flag that already has a []net.IPNet assigned, the newly converted values will be appended.
func (s *ipNetSliceValue) Set(val string) error {

	// remove all quote characters
	rmQuote := strings.NewReplacer(`"`, "", `'`, "", "`", "")

	// read flag arguments with CSV parser
	ipNetStrSlice, err := readAsCSV(rmQuote.Replace(val))
	if err != nil && err != io.EOF {
		return err
	}

	// parse ip values into slice
	out := make([]net.IPNet, 0, len(ipNetStrSlice))
	for _, ipNetStr := range ipNetStrSlice {
		_, n, err := net.ParseCIDR(strings.TrimSpace(ipNetStr))
		if err != nil {
			return fmt.Errorf("invalid string being converted to CIDR: %s", ipNetStr)
		}
		out = append(out, *n)
	}

	if !s.changed {
		*s.value = out
	} else {
		*s.value = append(*s.value, out...)
	}

	s.changed = true

	return nil
}

// Type returns a string that uniquely represents this flag's type.
func (s *ipNetSliceValue) Type() string {
	return "ipNetSlice"
}

// String defines a "native" format for this net.IPNet slice flag value.
func (s *ipNetSliceValue) String() string {

	ipNetStrSlice := make([]string, len(*s.value))
	for i, n := range *s.value {
		ipNetStrSlice[i] = n.String()
	}

	out, _ := writeAsCSV(ipNetStrSlice)
	return "[" + out + "]"
}

// GetIPNetSlice returns the []net.IPNet value of a flag with the given name
func (f *FlagSet) GetIPNetSlice(name string) ([]net.IPNet, error) {
	val, err := f.getFlagType(name, "ipNetSlice")
	if err != nil {
		return []net.IPNet{}, err
	}
	return val.([]net.IPNet), nil
}

// MustGetIPNetSlice is like GetIPNetSlice, but panics on error.
func (f *FlagSet) MustGetIPNetSlice(name string) []net.IPNet {
	val, err := f.GetIPNetSlice(name)
	if err != nil {
		panic(err)
	}
	return val
}

// IPNetSliceVar defines a []net.IPNet flag with specified name, default value, and usage string.
// The argument p points to a []net.IPNet variable in which to store the value of the flag.
func (f *FlagSet) IPNetSliceVar(p *[]net.IPNet, name string, value []net.IPNet, usage string, opts ...Opt) {
	f.Var(newIPNetSliceValue(value, p), name, usage, opts...)
}

// IPNetSliceVar defines a []net.IPNet flag with specified name, default value, and usage string.
// The argument p points to a []net.IPNet variable in which to store the value of the flag.
func IPNetSliceVar(p *[]net.IPNet, name string, value []net.IPNet, usage string, opts ...Opt) {
	CommandLine.IPNetSliceVar(p, name, value, usage, opts...)
}

// IPNetSlice defines a []net.IPNet flag with specified name, default value, and usage string.
// The return value is the address of a []net.IPNet variable that stores the value of the flag.
func (f *FlagSet) IPNetSlice(name string, value []net.IPNet, usage string, opts ...Opt) *[]net.IPNet {
	var p []net.IPNet
	f.IPNetSliceVar(&p, name, value, usage, opts...)
	return &p
}

// IPNetSlice defines a []net.IPNet flag with specified name, default value, and usage string.
// The return value is the address of a []net.IPNet variable that stores the value of the flag.
func IPNetSlice(name string, value []net.IPNet, usage string, opts ...Opt) *[]net.IPNet {
	return CommandLine.IPNetSlice(name, value, usage, opts...)
}
