// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"fmt"
	"net"
	"strings"
)

// -- net.IP value
type ipValue net.IP

func newIPValue(val net.IP, p *net.IP) *ipValue {
	*p = val
	return (*ipValue)(p)
}

func (i *ipValue) String() string { return net.IP(*i).String() }
func (i *ipValue) Set(s string) error {
	if s == "" {
		return nil
	}
	ip := net.ParseIP(strings.TrimSpace(s))
	if ip == nil {
		return fmt.Errorf("failed to parse IP: %q", s)
	}
	*i = ipValue(ip)
	return nil
}

func (i *ipValue) Get() interface{} {
	return net.IP(*i)
}

func (i *ipValue) Type() string {
	return "ip"
}

// GetIP return the net.IP value of a flag with the given name
func (f *FlagSet) GetIP(name string) (net.IP, error) {
	val, err := f.getFlagType(name, "ip")
	if err != nil {
		return nil, err
	}
	return val.(net.IP), nil
}

// MustGetIP is like GetIP, but panics on error.
func (f *FlagSet) MustGetIP(name string) net.IP {
	val, err := f.GetIP(name)
	if err != nil {
		panic(err)
	}
	return val
}

// IPVar defines a net.IP flag with specified name, default value, and usage string.
// The argument p points to a net.IP variable in which to store the value of the flag.
func (f *FlagSet) IPVar(p *net.IP, name string, value net.IP, usage string, opts ...Opt) {
	f.Var(newIPValue(value, p), name, usage, opts...)
}

// IPVar defines a net.IP flag with specified name, default value, and usage string.
// The argument p points to a net.IP variable in which to store the value of the flag.
func IPVar(p *net.IP, name string, value net.IP, usage string, opts ...Opt) {
	CommandLine.IPVar(p, name, value, usage, opts...)
}

// IP defines a net.IP flag with specified name, default value, and usage string.
// The return value is the address of a net.IP variable that stores the value of the flag.
func (f *FlagSet) IP(name string, value net.IP, usage string, opts ...Opt) *net.IP {
	var p net.IP
	f.IPVar(&p, name, value, usage, opts...)
	return &p
}

// IP defines a net.IP flag with specified name, default value, and usage string.
// The return value is the address of a net.IP variable that stores the value of the flag.
func IP(name string, value net.IP, usage string, opts ...Opt) *net.IP {
	return CommandLine.IP(name, value, usage, opts...)
}
