// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"fmt"
	"io"
	"net"
	"strings"
)

// -- ipSlice Value
type ipSliceValue struct {
	value   *[]net.IP
	changed bool
}

func newIPSliceValue(val []net.IP, p *[]net.IP) *ipSliceValue {
	ipsv := new(ipSliceValue)
	ipsv.value = p
	*ipsv.value = val
	return ipsv
}

// Set converts, and assigns, the comma-separated IP argument string representation as the []net.IP value of this flag.
// If Set is called on a flag that already has a []net.IP assigned, the newly converted values will be appended.
func (s *ipSliceValue) Set(val string) error {

	// remove all quote characters
	rmQuote := strings.NewReplacer(`"`, "", `'`, "", "`", "")

	// read flag arguments with CSV parser
	ipStrSlice, err := readAsCSV(rmQuote.Replace(val))
	if err != nil && err != io.EOF {
		return err
	}

	// parse ip values into slice
	out := make([]net.IP, 0, len(ipStrSlice))
	for _, ipStr := range ipStrSlice {
		ip := net.ParseIP(strings.TrimSpace(ipStr))
		if ip == nil {
			return fmt.Errorf("invalid string being converted to IP address: %s", ipStr)
		}
		out = append(out, ip)
	}

	if !s.changed {
		*s.value = out
	} else {
		*s.value = append(*s.value, out...)
	}

	s.changed = true

	return nil
}

func (s *ipSliceValue) Get() interface{} {
	return *s.value
}

// Type returns a string that uniquely represents this flag's type.
func (s *ipSliceValue) Type() string {
	return "ipSlice"
}

// String defines a "native" format for this net.IP slice flag value.
func (s *ipSliceValue) String() string {

	ipStrSlice := make([]string, len(*s.value))
	for i, ip := range *s.value {
		ipStrSlice[i] = ip.String()
	}

	out, _ := writeAsCSV(ipStrSlice)

	return "[" + out + "]"
}

func (s *ipSliceValue) fromString(val string) (net.IP, error) {
	return net.ParseIP(strings.TrimSpace(val)), nil
}

func (s *ipSliceValue) toString(val net.IP) string {
	return val.String()
}

func (s *ipSliceValue) Append(val string) error {
	i, err := s.fromString(val)
	if err != nil {
		return err
	}
	*s.value = append(*s.value, i)
	return nil
}

func (s *ipSliceValue) Replace(val []string) error {
	out := make([]net.IP, len(val))
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

func (s *ipSliceValue) GetSlice() []string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = s.toString(d)
	}
	return out
}

// GetIPSlice returns the []net.IP value of a flag with the given name
func (f *FlagSet) GetIPSlice(name string) ([]net.IP, error) {
	val, err := f.getFlagType(name, "ipSlice")
	if err != nil {
		return []net.IP{}, err
	}
	return val.([]net.IP), nil
}

// MustGetIPSlice is like GetIPSlice, but panics on error.
func (f *FlagSet) MustGetIPSlice(name string) []net.IP {
	val, err := f.GetIPSlice(name)
	if err != nil {
		panic(err)
	}
	return val
}

// IPSliceVar defines a []net.IP flag with specified name, default value, and usage string.
// The argument p points to a []net.IP variable in which to store the value of the flag.
func (f *FlagSet) IPSliceVar(p *[]net.IP, name string, value []net.IP, usage string, opts ...Opt) {
	f.Var(newIPSliceValue(value, p), name, usage, opts...)
}

// IPSliceVar defines a []net.IP flag with specified name, default value, and usage string.
// The argument p points to a []net.IP variable in which to store the value of the flag.
func IPSliceVar(p *[]net.IP, name string, value []net.IP, usage string, opts ...Opt) {
	CommandLine.IPSliceVar(p, name, value, usage, opts...)
}

// IPSlice defines a []net.IP flag with specified name, default value, and usage string.
// The return value is the address of a []net.IP variable that stores the value of the flag.
func (f *FlagSet) IPSlice(name string, value []net.IP, usage string, opts ...Opt) *[]net.IP {
	var p []net.IP
	f.IPSliceVar(&p, name, value, usage, opts...)
	return &p
}

// IPSlice defines a []net.IP flag with specified name, default value, and usage string.
// The return value is the address of a []net.IP variable that stores the value of the flag.
func IPSlice(name string, value []net.IP, usage string, opts ...Opt) *[]net.IP {
	return CommandLine.IPSlice(name, value, usage, opts...)
}
