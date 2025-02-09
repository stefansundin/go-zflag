// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

// -- stringSlice Value
type stringSliceValue struct {
	value   *[]string
	changed bool
}

func newStringSliceValue(val []string, p *[]string) *stringSliceValue {
	ssv := new(stringSliceValue)
	ssv.value = p
	*ssv.value = val
	return ssv
}

func readAsCSV(val string) ([]string, error) {
	if val == "" {
		return []string{}, nil
	}
	stringReader := strings.NewReader(val)
	csvReader := csv.NewReader(stringReader)
	return csvReader.Read()
}

func readCSVKeyValue(val string) (map[string]string, error) {
	strSlice, err := readAsCSV(val)
	if err != nil && err != io.EOF {
		return nil, err
	}

	out := make(map[string]string, len(strSlice))
	for _, pair := range strSlice {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("%s must be formatted as key=value", pair)
		}
		out[kv[0]] = kv[1]
	}

	return out, nil
}

func writeAsCSV(vals []string) (string, error) {
	b := &bytes.Buffer{}
	w := csv.NewWriter(b)
	err := w.Write(vals)
	if err != nil {
		return "", err
	}
	w.Flush()
	return strings.TrimSuffix(b.String(), "\n"), nil
}

func (s *stringSliceValue) Set(val string) error {
	v, err := readAsCSV(val)
	if err != nil {
		return err
	}
	if !s.changed {
		*s.value = v
	} else {
		*s.value = append(*s.value, v...)
	}
	s.changed = true
	return nil
}

func (s *stringSliceValue) Get() interface{} {
	return *s.value
}

func (s *stringSliceValue) Type() string {
	return "stringSlice"
}

func (s *stringSliceValue) String() string {
	str, _ := writeAsCSV(*s.value)
	return "[" + str + "]"
}

func (s *stringSliceValue) Append(val string) error {
	*s.value = append(*s.value, val)
	return nil
}

func (s *stringSliceValue) Replace(val []string) error {
	*s.value = val
	return nil
}

func (s *stringSliceValue) GetSlice() []string {
	return *s.value
}

// GetStringSlice return the []string value of a flag with the given name
func (f *FlagSet) GetStringSlice(name string) ([]string, error) {
	val, err := f.getFlagType(name, "stringSlice")
	if err != nil {
		return []string{}, err
	}
	return val.([]string), nil
}

// MustGetStringSlice is like GetStringSlice, but panics on error.
func (f *FlagSet) MustGetStringSlice(name string) []string {
	val, err := f.GetStringSlice(name)
	if err != nil {
		panic(err)
	}
	return val
}

// StringSliceVar defines a []string flag with specified name, default value, and usage string.
// The argument p points to a []string variable in which to store the value of the flag.
// Compared to StringArray flags, StringSlice flags take comma-separated value as arguments and split them accordingly.
// For example:
//   --ss="v1,v2" --ss="v3"
// will result in
//   []string{"v1", "v2", "v3"}
func (f *FlagSet) StringSliceVar(p *[]string, name string, value []string, usage string, opts ...Opt) {
	f.Var(newStringSliceValue(value, p), name, usage, opts...)
}

// StringSliceVar defines a []string flag with specified name, default value, and usage string.
// The argument p points to a []string variable in which to store the value of the flag.
// Compared to StringArray flags, StringSlice flags take comma-separated value as arguments and split them accordingly.
// For example:
//   --ss="v1,v2" --ss="v3"
// will result in
//   []string{"v1", "v2", "v3"}
func StringSliceVar(p *[]string, name string, value []string, usage string, opts ...Opt) {
	CommandLine.StringSliceVar(p, name, value, usage, opts...)
}

// StringSlice defines a []string flag with specified name, default value, and usage string.
// The return value is the address of a []string variable that stores the value of the flag.
// Compared to StringArray flags, StringSlice flags take comma-separated value as arguments and split them accordingly.
// For example:
//   --ss="v1,v2" --ss="v3"
// will result in
//   []string{"v1", "v2", "v3"}
func (f *FlagSet) StringSlice(name string, value []string, usage string, opts ...Opt) *[]string {
	var p []string
	f.StringSliceVar(&p, name, value, usage, opts...)
	return &p
}

// StringSlice defines a []string flag with specified name, default value, and usage string.
// The return value is the address of a []string variable that stores the value of the flag.
// Compared to StringArray flags, StringSlice flags take comma-separated value as arguments and split them accordingly.
// For example:
//   --ss="v1,v2" --ss="v3"
// will result in
//   []string{"v1", "v2", "v3"}
func StringSlice(name string, value []string, usage string, opts ...Opt) *[]string {
	return CommandLine.StringSlice(name, value, usage, opts...)
}
