// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"bytes"
	"errors"
	goflag "flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"unicode/utf8"
)

// ErrHelp is the error returned if the flag -help is invoked but no such flag is defined.
var ErrHelp = errors.New("zflag: help requested")

// ErrorHandling defines how to handle flag parsing errors.
type ErrorHandling int

const (
	// ContinueOnError will return an err from Parse() if an error is found
	ContinueOnError ErrorHandling = iota
	// ExitOnError will call os.Exit(2) if an error is found when parsing
	ExitOnError
	// PanicOnError will panic() if an error is found when parsing flags
	PanicOnError
)

// ParseErrorsAllowlist defines the parsing errors that can be ignored
type ParseErrorsAllowlist struct {
	// UnknownFlags will ignore unknown flags errors and continue parsing rest of the flags
	// See GetUnknownFlags to retrieve collected unknowns.
	UnknownFlags bool
}

// NormalizedName is a flag name that has been normalized according to rules
// for the FlagSet (e.g. making '-' and '_' equivalent).
type NormalizedName string

// A FlagSet represents a set of defined flags.
type FlagSet struct {
	// Usage is the function called when an error occurs while parsing flags.
	// The field is a function (not a method) that may be changed to point to
	// a custom error handler.
	Usage func()

	// SortFlags is used to indicate, if user wants to have sorted flags in
	// help/usage messages.
	SortFlags bool

	// ParseErrorsAllowlist is used to configure an allowlist of errors
	ParseErrorsAllowlist ParseErrorsAllowlist

	// DisableBuiltinHelp toggles the built-in convention of handling -h and --help
	DisableBuiltinHelp bool

	// FlagUsageFormatter allows for custom formatting of flag usage output.
	// Each individual item needs to be implemented. See FlagUsagesForGroupWrapped for info on what gets passed.
	FlagUsageFormatter FlagUsageFormatter

	name              string
	parsed            bool
	actual            map[NormalizedName]*Flag
	orderedActual     []*Flag
	sortedActual      []*Flag
	formal            map[NormalizedName]*Flag
	orderedFormal     []*Flag
	sortedFormal      []*Flag
	shorthands        map[rune]*Flag
	args              []string // arguments after flags
	argsLenAtDash     int      // len(args) when a '--' was located when parsing, or -1 if no --
	errorHandling     ErrorHandling
	output            io.Writer // nil means stderr; use Output() accessor
	interspersed      bool      // allow interspersed option/non-option args
	normalizeNameFunc func(f *FlagSet, name string) NormalizedName

	addedGoFlagSets []*goflag.FlagSet
	unknownFlags    []string
}

// A Flag represents the state of a flag.
type Flag struct {
	Name                string              // name as it appears on command line
	Shorthand           rune                // one-letter abbreviated flag
	ShorthandOnly       bool                // If the user set only the shorthand
	Usage               string              // help message
	UsageType           string              // flag type displayed in the help message
	DisableUnquoteUsage bool                // toggle unquoting and extraction of type from usage
	DisablePrintDefault bool                // toggle printing of the default value in usage message
	Value               Value               // value as set
	DefValue            string              // default value (as text); for usage message
	Changed             bool                // If the user set the value (or if left to default)
	NoOptDefVal         string              // default value (as text); if the flag is on the command line without any options
	Deprecated          string              // If this flag is deprecated, this string is the new or now thing to use
	Hidden              bool                // used by zulu.Command to allow flags to be hidden from help/usage text
	ShorthandDeprecated string              // If the shorthand of this flag is deprecated, this string is the new or now thing to use
	Group               string              // flag group
	Annotations         map[string][]string // Use it to annotate this specific flag for your application; used by zulu.Command bash completion code
}

// Value is the interface to the dynamic value stored in a flag.
// (The default value is represented as a string.)
type Value interface {
	String() string
	Set(string) error
}

type Getter interface {
	Value
	Get() interface{}
}

// Typed is an interface of Values that can communicate their type.
type Typed interface {
	Type() string
}

// SliceValue is a secondary interface to all flags which hold a list
// of values.  This allows full control over the value of list flags,
// and avoids complicated marshalling and unmarshalling to csv.
type SliceValue interface {
	// Append adds the specified value to the end of the flag value list.
	Append(string) error
	// Replace will fully overwrite any data currently in the flag value list.
	Replace([]string) error
	// GetSlice returns the flag value list as an array of strings.
	GetSlice() []string
}

// sortFlags returns the flags as a slice in lexicographical sorted order.
func sortFlags(flags map[NormalizedName]*Flag) []*Flag {
	list := make(sort.StringSlice, len(flags))
	i := 0
	for k := range flags {
		list[i] = string(k)
		i++
	}
	list.Sort()
	result := make([]*Flag, len(list))
	for i, name := range list {
		result[i] = flags[NormalizedName(name)]
	}
	return result
}

// SetNormalizeFunc allows you to add a function which can translate flag names.
// Flags added to the FlagSet will be translated and then when anything tries to
// look up the flag that will also be translated. So it would be possible to create
// a flag named "getURL" and have it translated to "geturl".  A user could then pass
// "--getUrl" which may also be translated to "geturl" and everything will work.
func (f *FlagSet) SetNormalizeFunc(n func(f *FlagSet, name string) NormalizedName) {
	f.normalizeNameFunc = n
	f.sortedFormal = f.sortedFormal[:0]
	for fname, flag := range f.formal {
		nname := f.normalizeFlagName(flag.Name)
		if fname == nname {
			continue
		}
		flag.Name = string(nname)
		delete(f.formal, fname)
		f.formal[nname] = flag
		if _, set := f.actual[fname]; set {
			delete(f.actual, fname)
			f.actual[nname] = flag
		}
	}
}

// GetNormalizeFunc returns the previously set NormalizeFunc of a function which
// does no translation, if not set previously.
func (f *FlagSet) GetNormalizeFunc() func(f *FlagSet, name string) NormalizedName {
	if f.normalizeNameFunc != nil {
		return f.normalizeNameFunc
	}
	return func(f *FlagSet, name string) NormalizedName { return NormalizedName(name) }
}

func (f *FlagSet) normalizeFlagName(name string) NormalizedName {
	n := f.GetNormalizeFunc()
	return n(f, name)
}

// Output returns the destination for usage and error messages. os.Stderr is returned if
// output was not set or was set to nil.
func (f *FlagSet) Output() io.Writer {
	if f.output == nil {
		return os.Stderr
	}
	return f.output
}

// Name returns the name of the flag set.
func (f *FlagSet) Name() string {
	return f.name
}

// SetOutput sets the destination for usage and error messages.
// If output is nil, os.Stderr is used.
func (f *FlagSet) SetOutput(output io.Writer) {
	f.output = output
}

// GetAllFlags return the flags in lexicographical order or
// in primordial order if f.SortFlags is false.
// It visits all flags, even those not set.
func (f *FlagSet) GetAllFlags() (flags []*Flag) {
	if f.SortFlags {
		if len(f.formal) != len(f.sortedFormal) {
			f.sortedFormal = sortFlags(f.formal)
		}
		flags = f.sortedFormal
	} else {
		flags = f.orderedFormal
	}
	return
}

// VisitAll visits the flags in lexicographical order or
// in primordial order if f.SortFlags is false, calling fn for each.
// It visits all flags, even those not set.
func (f *FlagSet) VisitAll(fn func(*Flag)) {
	if len(f.formal) == 0 {
		return
	}
	for _, flag := range f.GetAllFlags() {
		fn(flag)
	}
}

// HasFlags returns a bool to indicate if the FlagSet has any flags defined.
func (f *FlagSet) HasFlags() bool {
	return len(f.formal) > 0
}

// HasAvailableFlags returns a bool to indicate if the FlagSet has any flags
// that are not hidden.
func (f *FlagSet) HasAvailableFlags() bool {
	for _, flag := range f.formal {
		if !flag.Hidden {
			return true
		}
	}
	return false
}

// GetAllFlags return the flags in lexicographical order or
// in primordial order if f.SortFlags is false.
func GetAllFlags() []*Flag {
	return CommandLine.GetAllFlags()
}

// VisitAll visits the command-line flags in lexicographical order or
// in primordial order if f.SortFlags is false, calling fn for each.
// It visits all flags, even those not set.
func VisitAll(fn func(*Flag)) {
	CommandLine.VisitAll(fn)
}

// GetFlags return the flags in lexicographical order or
// in primordial order if f.SortFlags is false.
// It visits only those flags that have been set.
func (f *FlagSet) GetFlags() (flags []*Flag) {
	if f.SortFlags {
		if len(f.actual) != len(f.sortedActual) {
			f.sortedActual = sortFlags(f.actual)
		}
		flags = f.sortedActual
	} else {
		flags = f.orderedActual
	}
	return
}

// Visit visits the flags in lexicographical order or
// in primordial order if f.SortFlags is false, calling fn for each.
// It visits only those flags that have been set.
func (f *FlagSet) Visit(fn func(*Flag)) {
	if len(f.actual) == 0 {
		return
	}
	for _, flag := range f.GetFlags() {
		fn(flag)
	}
}

// GetFlags return the flags in lexicographical order or
// in primordial order if f.SortFlags is false.
func GetFlags() []*Flag {
	return CommandLine.GetFlags()
}

// Visit visits the command-line flags in lexicographical order or
// in primordial order if f.SortFlags is false, calling fn for each.
// It visits only those flags that have been set.
func Visit(fn func(*Flag)) {
	CommandLine.Visit(fn)
}

func (f *FlagSet) addUnknownFlag(s string) {
	f.unknownFlags = append(f.unknownFlags, s)
}

// GetUnknownFlags returns unknown flags in the order they were Parsed.
// This requires ParseErrorsWhitelist.UnknownFlags to be set so that parsing does
// not abort on the first unknown flag.
func (f *FlagSet) GetUnknownFlags() []string {
	return f.unknownFlags
}

// GetUnknownFlags returns unknown command-line flags in the order they were Parsed.
// This requires ParseErrorsWhitelist.UnknownFlags to be set so that parsing does
// not abort on the first unknown flag.
func GetUnknownFlags() []string {
	return CommandLine.GetUnknownFlags()
}

// Get returns the value of the named flag.
func (f *FlagSet) Get(name string) (interface{}, error) {
	return f.getFlagType(name, "")
}

// Get returns the value of the named flag.
func Get(name string) (interface{}, error) {
	return CommandLine.Get(name)
}

// Lookup returns the Flag structure of the named flag, returning nil if none exists.
func (f *FlagSet) Lookup(name string) *Flag {
	return f.lookup(f.normalizeFlagName(name))
}

// ShorthandLookup returns the Flag structure of the shorthand flag,
// returning nil if none exists.
func (f *FlagSet) ShorthandLookup(name rune) *Flag {
	if name == 0 {
		return nil
	}

	v, ok := f.shorthands[name]
	if !ok {
		return nil
	}
	return v
}

// ShorthandLookupStr is the same as ShorthandLookup, but you can look it up through a string.
// It panics if name contains more than one UTF-8 character.
func (f *FlagSet) ShorthandLookupStr(name string) *Flag {
	r, err := shorthandStrToRune(name)
	if err != nil {
		fmt.Fprintln(f.Output(), err)
		panic(err)
	}

	return f.ShorthandLookup(r)
}

func shorthandStrToRune(name string) (rune, error) {
	if utf8.RuneCountInString(name) > 1 {
		return 0, fmt.Errorf("cannot convert shorthand with more than one UTF-8 character: %q", name)
	}
	r, _ := utf8.DecodeRuneInString(name)
	if r == utf8.RuneError {
		return 0, nil
	}

	return r, nil
}

// lookup returns the Flag structure of the named flag, returning nil if none exists.
func (f *FlagSet) lookup(name NormalizedName) *Flag {
	return f.formal[name]
}

// func to return a given type for a given flag name
func (f *FlagSet) getFlagType(name string, ftype string) (interface{}, error) {
	flag := f.Lookup(name)
	if flag == nil {
		err := fmt.Errorf("flag accessed but not defined: %s", name)
		return nil, err
	}

	if ftype != "" {
		if v, ok := flag.Value.(Typed); ok && v.Type() != ftype {
			err := fmt.Errorf("trying to get %q value of flag of type %q", ftype, v.Type())
			return nil, err
		}
	}

	getter, ok := flag.Value.(Getter)
	if !ok {
		return nil, fmt.Errorf("flag %q does not implement the Getter interface", name)
	}

	return getter.Get(), nil
}

// ArgsLenAtDash will return the length of f.Args at the moment when a -- was
// found during arg parsing. This allows your program to know which args were
// before the -- and which came after.
func (f *FlagSet) ArgsLenAtDash() int {
	return f.argsLenAtDash
}

// Lookup returns the Flag structure of the named command-line flag,
// returning nil if none exists.
func Lookup(name string) *Flag {
	return CommandLine.Lookup(name)
}

// ShorthandLookup returns the Flag structure of the shorthand flag,
// returning nil if none exists.
func ShorthandLookup(name rune) *Flag {
	return CommandLine.ShorthandLookup(name)
}

// ShorthandLookupStr is the same as ShorthandLookup, but you can look it up through a string.
// It panics if name contains more than one UTF-8 character.
func ShorthandLookupStr(name string) *Flag {
	return CommandLine.ShorthandLookupStr(name)
}

// Set sets the value of the named flag.
func (f *FlagSet) Set(name, value string) error {
	normalName := f.normalizeFlagName(name)
	flag, ok := f.formal[normalName]
	if !ok {
		return NewUnknownFlagError(name)
	}

	err := flag.Value.Set(value)
	if err != nil {
		var flagName string
		if flag.Shorthand != 0 && flag.ShorthandDeprecated == "" {
			flagName = fmt.Sprintf("-%c", flag.Shorthand)
			if !flag.ShorthandOnly {
				flagName = fmt.Sprintf("%s, --%s", flagName, flag.Name)
			}
		} else {
			flagName = fmt.Sprintf("--%s", flag.Name)
		}
		return fmt.Errorf("invalid argument %q for %q flag: %v", value, flagName, err)
	}

	if !flag.Changed {
		if f.actual == nil {
			f.actual = make(map[NormalizedName]*Flag)
		}
		f.actual[normalName] = flag
		f.orderedActual = append(f.orderedActual, flag)

		flag.Changed = true
	}

	if flag.Deprecated != "" {
		fmt.Fprintf(f.Output(), "Flag --%s has been deprecated, %s\n", flag.Name, flag.Deprecated)
	}
	return nil
}

// SetAnnotation allows one to set arbitrary annotations on this flag.
// This is sometimes used by gowarden/zulu programs which want to generate additional
// bash completion information.
func (f *Flag) SetAnnotation(key string, values []string) error {
	if f.Annotations == nil {
		f.Annotations = map[string][]string{}
	}

	f.Annotations[key] = values

	return nil
}

// SetAnnotation allows one to set arbitrary annotations on a flag in the FlagSet.
// This is sometimes used by gowarden/zulu programs which want to generate additional
// bash completion information.
func (f *FlagSet) SetAnnotation(name, key string, values []string) error {
	normalName := f.normalizeFlagName(name)
	flag, ok := f.formal[normalName]
	if !ok {
		return NewUnknownFlagError(name)
	}
	return flag.SetAnnotation(key, values)
}

// Changed returns true if the flag was explicitly set during Parse() and false
// otherwise
func (f *FlagSet) Changed(name string) bool {
	flag := f.Lookup(name)
	// If a flag doesn't exist, it wasn't changed....
	if flag == nil {
		return false
	}
	return flag.Changed
}

// Set sets the value of the named command-line flag.
func Set(name, value string) error {
	return CommandLine.Set(name, value)
}

// PrintDefaults prints to standard error unless configured otherwise, the
// default values of all defined command-line flags in the set. See the
// documentation for the global function PrintDefaults for more information.
func (f *FlagSet) PrintDefaults() {
	usages := f.FlagUsages()
	fmt.Fprint(f.Output(), usages)
}

// defaultIsZeroValue returns true if the default value for this flag represents
// a zero value.
func (f *Flag) defaultIsZeroValue() bool {
	switch f.Value.(type) {
	case boolFlag:
		return f.DefValue == "false"
	case *durationValue:
		// Beginning in Go 1.7, duration zero values are "0s"
		return f.DefValue == "0" || f.DefValue == "0s"
	case *intValue, *int8Value, *int32Value, *int64Value, *uintValue, *uint8Value, *uint16Value, *uint32Value, *uint64Value, *countValue, *float32Value, *float64Value:
		return f.DefValue == "0"
	case *stringValue:
		return f.DefValue == ""
	case *ipValue, *ipMaskValue, *ipNetValue:
		return f.DefValue == "<nil>"
	case *intSliceValue, *stringSliceValue, *stringArrayValue:
		return f.DefValue == "[]"
	default:
		switch f.DefValue {
		case "false":
			return true
		case "<nil>":
			return true
		case "":
			return true
		case "0":
			return true
		}
		return false
	}
}

// UnquoteUsage extracts a back-quoted name from the usage
// string for a flag and returns it and the un-quoted usage.
// Given "a `name` to show" it returns ("name", "a name to show").
// If there are no back quotes, the name is an educated guess of the
// type of the flag's value, or the empty string if the flag is boolean.
func UnquoteUsage(flag *Flag) (name string, usage string) {
	name = flag.UsageType
	usage = flag.Usage

	// Look for a back-quoted name, but avoid the strings package.
	if !flag.DisableUnquoteUsage {
		for i := 0; i < len(usage); i++ {
			if usage[i] == '`' {
				for j := i + 1; j < len(usage); j++ {
					if usage[j] == '`' {
						extracted := usage[i+1 : j]
						if name == "" {
							name = extracted
						}
						usage = usage[:i] + extracted + usage[j+1:]
						return
					}
				}
				break // Only one back quote; use type name.
			}
		}
	}

	if name == "" {
		name = "value" // compatibility layer to be a drop-in replacement
		if v, ok := flag.Value.(Typed); ok {
			name = v.Type()
			switch name {
			case "bool":
				name = ""
			case "boolSlice":
				name = "bools"
			case "complex128":
				name = "complex"
			case "complex128Slice":
				name = "complexes"
			case "durationSlice":
				name = "durations"
			case "float32", "float64":
				name = "float"
			case "floatSlice", "float32Slice", "float64Slice":
				name = "floats"
			case "int8", "int16", "int32", "int64":
				name = "int"
			case "intSlice", "int8Slice", "int16Slice", "int32Slice", "int64Slice":
				name = "ints"
			case "stringSlice", "stringArray":
				name = "strings"
			case "uint8", "uint16", "uint32", "uint64":
				name = "uint"
			case "uintSlice", "uint8Slice", "uint16Slice", "uint32Slice", "uint64Slice":
				name = "uints"
			}
		}
	}

	return
}

// Splits the string `s` on whitespace into an initial substring up to
// `i` runes in length and the remainder. Will go `slop` over `i` if
// that encompasses the entire string (which allows the caller to
// avoid short orphan words on the final line).
func wrapN(i, slop int, s string) (string, string) {
	if i+slop > len(s) {
		return s, ""
	}

	w := strings.LastIndexAny(s[:i], " \t\n")
	if w <= 0 {
		return s, ""
	}
	nlPos := strings.LastIndex(s[:i], "\n")
	if nlPos > 0 && nlPos < w {
		return s[:nlPos], s[nlPos+1:]
	}
	return s[:w], s[w+1:]
}

// Wraps the string `s` to a maximum width `w` with leading indent
// `i`. The first line is not indented (this is assumed to be done by
// caller). Pass `w` == 0 to do no wrapping
func wrap(i, w int, s string) string {
	if w == 0 {
		return strings.Replace(s, "\n", "\n"+strings.Repeat(" ", i), -1)
	}

	// space between indent i and end of line width w into which
	// we should wrap the text.
	wrap := w - i

	var r, l string

	// Not enough space for sensible wrapping. Wrap as a block on
	// the next line instead.
	if wrap < 24 {
		i = 16
		wrap = w - i
		r += "\n" + strings.Repeat(" ", i)
	}
	// If still not enough space then don't even try to wrap.
	if wrap < 24 {
		return strings.Replace(s, "\n", r, -1)
	}

	// Try to avoid short orphan words on the final line, by
	// allowing wrapN to go a bit over if that would fit in the
	// remainder of the line.
	slop := 5
	wrap = wrap - slop

	// Handle first line, which is indented by the caller (or the
	// special case above)
	l, s = wrapN(wrap, slop, s)
	r = r + strings.Replace(l, "\n", "\n"+strings.Repeat(" ", i), -1)

	// Now wrap the rest
	for s != "" {
		var t string

		t, s = wrapN(wrap, slop, s)
		r = r + "\n" + strings.Repeat(" ", i) + strings.Replace(t, "\n", "\n"+strings.Repeat(" ", i), -1)
	}

	return r

}

func (f *FlagSet) flagUsageFormatter() FlagUsageFormatter {
	if f.FlagUsageFormatter == nil {
		return DefaultFlagUsageFormatter{}
	}

	return f.FlagUsageFormatter
}

// FlagUsagesWrapped returns a string containing the usage information
// for all flags in the FlagSet. Wrapped to `cols` columns (0 for no
// wrapping)
func (f *FlagSet) FlagUsagesWrapped(cols int) string {
	return f.FlagUsagesForGroupWrapped("", cols)
}

// FlagUsagesForGroupWrapped returns a string containing the usage information
// for all flags in the FlagSet for group. Wrapped to `cols` columns (0 for no
// wrapping)
func (f *FlagSet) FlagUsagesForGroupWrapped(group string, cols int) string {
	buf := new(bytes.Buffer)

	lines := make(map[string][]string)

	usageFormatter := f.flagUsageFormatter()

	maxlen := 0
	f.VisitAll(func(flag *Flag) {
		if flag.Hidden {
			return
		}

		line := usageFormatter.Name(flag)

		varname, usage := UnquoteUsage(flag)
		if varname != "" {
			line += " " + usageFormatter.UsageVarName(flag, varname)
		}
		if flag.NoOptDefVal != "" {
			line += usageFormatter.NoOptDefValue(flag)
		}

		// This special character will be replaced with spacing once the
		// correct alignment is calculated
		line += "\x00"
		if len(line) > maxlen {
			maxlen = len(line)
		}

		line += usageFormatter.Usage(flag, usage)
		if !flag.DisablePrintDefault && !flag.defaultIsZeroValue() {
			line += usageFormatter.DefaultValue(flag)
		}
		if len(flag.Deprecated) != 0 {
			line += usageFormatter.Deprecated(flag)
		}

		group := flag.Group
		if _, ok := lines[group]; !ok {
			lines[group] = make([]string, 0)
		}
		lines[group] = append(lines[group], line)
	})

	for _, line := range lines[group] {
		sidx := strings.Index(line, "\x00")
		spacing := strings.Repeat(" ", maxlen-sidx)
		// maxlen + 2 comes from + 1 for the \x00 and + 1 for the (deliberate) off-by-one in maxlen-sidx
		fmt.Fprintln(buf, line[:sidx], spacing, wrap(maxlen+2, cols, line[sidx+1:]))
	}

	return buf.String()
}

// FlagUsages returns a string containing the usage information for all flags in
// the FlagSet
func (f *FlagSet) FlagUsages() string {
	return f.FlagUsagesWrapped(0)
}

// FlagUsagesForGroup returns a string containing the usage information for all flags in
// the FlagSet for group
func (f *FlagSet) FlagUsagesForGroup(group string) string {
	return f.FlagUsagesForGroupWrapped(group, 0)
}

// Groups return an array of unique flag groups sorted in the same order
// as flags. Empty group (unassigned) is always placed at the beginning.
func (f *FlagSet) Groups() []string {
	groupsMap := make(map[string]bool)
	groups := make([]string, 0)
	hasUngrouped := false
	f.VisitAll(func(flag *Flag) {
		if flag.Group == "" {
			hasUngrouped = true
			return
		}
		if _, ok := groupsMap[flag.Group]; !ok {
			groupsMap[flag.Group] = true
			groups = append(groups, flag.Group)
		}
	})
	sort.Strings(groups)

	if hasUngrouped {
		groups = append([]string{""}, groups...)
	}

	return groups
}

// PrintDefaults prints, to standard error unless configured otherwise,
// a usage message showing the default settings of all defined
// command-line flags.
// For an integer valued flag x, the default output has the form
//	-x int
//		usage-message-for-x (default 7)
// The usage message will appear on a separate line for anything but
// a bool flag with a one-byte name. For bool flags, the type is
// omitted and if the flag name is one byte the usage message appears
// on the same line. The parenthetical default is omitted if the
// default is the zero value for the type. The listed type, here int,
// can be changed by placing a back-quoted name in the flag's usage
// string; the first such item in the message is taken to be a parameter
// name to show in the message and the back quotes are stripped from
// the message when displayed. For instance, given
//	flag.String("I", "", "search `directory` for include files")
// the output will be
//	-I directory
//		search directory for include files.
//
// To change the destination for flag messages, call CommandLine.SetOutput.
func PrintDefaults() {
	CommandLine.PrintDefaults()
}

// defaultUsage is the default function to print a usage message.
func (f *FlagSet) defaultUsage() {
	if f.name == "" {
		fmt.Fprintf(f.Output(), "Usage:\n")
	} else {
		fmt.Fprintf(f.Output(), "Usage of %s:\n", f.name)
	}
	f.PrintDefaults()
}

// NOTE: Usage is not just CommandLine.defaultUsage()
// because it serves (via godoc flag Usage) as the example
// for how to write your own usage function.

// Usage prints to standard error a usage message documenting all defined command-line flags.
// The function is a variable that may be changed to point to a custom function.
// By default it prints a simple header and calls PrintDefaults; for details about the
// format of the output and how to control it, see the documentation for PrintDefaults.
var Usage = func() {
	fmt.Fprintf(CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	PrintDefaults()
}

// NFlag returns the number of flags that have been set.
func (f *FlagSet) NFlag() int { return len(f.actual) }

// NFlag returns the number of command-line flags that have been set.
func NFlag() int { return len(CommandLine.actual) }

// Arg returns the i'th argument.  Arg(0) is the first remaining argument
// after flags have been processed.
func (f *FlagSet) Arg(i int) string {
	if i < 0 || i >= len(f.args) {
		return ""
	}
	return f.args[i]
}

// Arg returns the i'th command-line argument.  Arg(0) is the first remaining argument
// after flags have been processed.
func Arg(i int) string {
	return CommandLine.Arg(i)
}

// NArg is the number of arguments remaining after flags have been processed.
func (f *FlagSet) NArg() int { return len(f.args) }

// NArg is the number of arguments remaining after flags have been processed.
func NArg() int { return len(CommandLine.args) }

// Args returns the non-flag arguments.
func (f *FlagSet) Args() []string { return f.args }

// Args returns the non-flag command-line arguments.
func Args() []string { return CommandLine.args }

// Var defines a flag with the specified name and usage string. The type and
// value of the flag are represented by the first argument, of type Value, which
// typically holds a user-defined implementation of Value. For instance, the
// caller could create a flag that turns a comma-separated string into a slice
// of strings by giving the slice the methods of Value; in particular, Set would
// decompose the comma-separated string into the slice.
func (f *FlagSet) Var(value Value, name, usage string, opts ...Opt) *Flag {
	flag := &Flag{
		Name:     name,
		Usage:    usage,
		Value:    value,
		DefValue: value.String(),
	}

	err := applyFlagOptions(flag, opts...)
	if err != nil {
		panic(err)
	}

	f.AddFlag(flag)
	return flag
}

// AddFlag will add the flag to the FlagSet
func (f *FlagSet) AddFlag(flag *Flag) {
	normalizedFlagName := f.normalizeFlagName(flag.Name)

	_, alreadyThere := f.formal[normalizedFlagName]
	if alreadyThere {
		msg := fmt.Sprintf("%s flag redefined: %s", f.name, flag.Name)
		fmt.Fprintln(f.Output(), msg)
		panic(msg) // Happens only if flags are declared with identical names
	}
	if f.formal == nil {
		f.formal = make(map[NormalizedName]*Flag)
	}

	flag.Name = string(normalizedFlagName)
	f.formal[normalizedFlagName] = flag
	f.orderedFormal = append(f.orderedFormal, flag)

	if flag.Shorthand == 0 {
		return
	}
	if f.shorthands == nil {
		f.shorthands = make(map[rune]*Flag)
	}
	used, alreadyThere := f.shorthands[flag.Shorthand]
	if alreadyThere {
		msg := fmt.Sprintf("unable to redefine %q shorthand in %q flagset: it's already used for %q flag", flag.Shorthand, f.name, used.Name)
		fmt.Fprintln(f.Output(), msg)
		panic(msg)
	}
	f.shorthands[flag.Shorthand] = flag
}

// AddFlagSet adds one FlagSet to another. If a flag is already present in f
// the flag from newSet will be ignored.
func (f *FlagSet) AddFlagSet(newSet *FlagSet) {
	if newSet == nil {
		return
	}
	newSet.VisitAll(func(flag *Flag) {
		if f.Lookup(flag.Name) == nil {
			f.AddFlag(flag)
		}
	})
}

// Var defines a flag with the specified name and usage string. The type and
// value of the flag are represented by the first argument, of type Value, which
// typically holds a user-defined implementation of Value. For instance, the
// caller could create a flag that turns a comma-separated string into a slice
// of strings by giving the slice the methods of Value; in particular, Set would
// decompose the comma-separated string into the slice.
func Var(value Value, name, usage string, opts ...Opt) *Flag {
	return CommandLine.Var(value, name, usage, opts...)
}

// failf prints to standard error a formatted error and usage message and
// returns the error.
func (f *FlagSet) failf(format string, a ...interface{}) error {
	f.usage()
	err := fmt.Errorf(format, a...)
	fmt.Fprintln(f.Output())
	fmt.Fprintln(f.Output(), err)
	return err
}

// usage calls the Usage method for the flag set, or the usage function if
// the flag set is CommandLine.
func (f *FlagSet) usage() {
	if f == CommandLine {
		Usage()
	} else if f.Usage == nil {
		f.defaultUsage()
	} else {
		f.Usage()
	}
}

// --unknown (args will be empty)
// --unknown --next-flag ... (args will be --next-flag ...)
// --unknown arg ... (args will be arg ...)
func (f *FlagSet) stripUnknownFlagValue(args []string) []string {
	if len(args) == 0 {
		// --unknown
		return args
	}

	first := args[0]
	if len(first) > 0 && first[0] == '-' {
		// --unknown --next-flag ...
		return args
	}

	// --unknown arg ... (args will be arg ...)
	if len(args) > 1 {
		f.addUnknownFlag(args[0])
		return args[1:]
	}
	return nil
}

func (f *FlagSet) parseLongArg(s string, args []string, fn parseFunc) (outArgs []string, err error) {
	outArgs = args
	name := s[2:]
	if len(name) == 0 || name[0] == '-' || name[0] == '=' {
		err = f.failf("bad flag syntax: %s", s)
		return
	}

	split := strings.SplitN(name, "=", 2)
	name = split[0]
	flag, exists := f.formal[f.normalizeFlagName(name)]

	if !exists || (flag != nil && flag.ShorthandOnly) {
		switch {
		case !exists && name == "help" && !f.DisableBuiltinHelp:
			f.usage()
			err = ErrHelp
			return
		case f.ParseErrorsAllowlist.UnknownFlags || (flag != nil && flag.ShorthandOnly):
			// --unknown=unknownval arg ...
			// we do not want to lose arg in this case
			f.addUnknownFlag(s)
			if len(split) >= 2 {
				return
			}
			outArgs = f.stripUnknownFlagValue(outArgs)
			return
		default:
			err = f.failf(NewUnknownFlagError(name).Error())
			return
		}
	}

	var value string
	if len(split) == 2 {
		// '--flag=arg'
		value = split[1]
	} else if flag.NoOptDefVal != "" {
		// '--flag' (arg was optional)
		value = flag.NoOptDefVal
	} else if len(outArgs) > 0 {
		// '--flag arg'
		value = outArgs[0]
		outArgs = outArgs[1:]
	} else {
		// '--flag' (arg was required)
		err = f.failf("flag needs an argument: %s", s)
		return
	}

	err = fn(flag, value)
	if err != nil {
		err = f.failf(err.Error())
	}
	return
}

func (f *FlagSet) parseSingleShortArg(shorthands string, args []string, fn parseFunc) (outShorts string, outArgs []string, err error) {
	outArgs = args
	outShorts = shorthands[1:]
	char, _ := utf8.DecodeRuneInString(shorthands)

	flag, exists := f.shorthands[char]
	if !exists {
		switch {
		case char == 'h' && !f.DisableBuiltinHelp:
			f.usage()
			err = ErrHelp
			return
		case f.ParseErrorsAllowlist.UnknownFlags:
			if len(shorthands) > 2 {
				// '-f...'
				// we do not want to lose anything in this case
				f.addUnknownFlag("-" + shorthands)
				outShorts = ""
				return
			}
			f.addUnknownFlag("-" + string(char))
			if len(outShorts) == 0 {
				outArgs = f.stripUnknownFlagValue(outArgs)
			}
			return
		default:
			err = f.failf("unknown shorthand flag: %q in -%s", char, shorthands)
			return
		}
	}

	var value string
	if len(shorthands) > 2 && shorthands[1] == '=' {
		// '-f=arg'
		value = shorthands[2:]
		outShorts = ""
	} else if flag.NoOptDefVal != "" {
		// '-f' (arg was optional)
		value = flag.NoOptDefVal
	} else if len(shorthands) > 1 {
		// '-farg'
		value = shorthands[1:]
		outShorts = ""
	} else if len(args) > 0 {
		// '-f arg'
		value = args[0]
		outArgs = args[1:]
	} else {
		// '-f' (arg was required)
		err = f.failf("flag needs an argument: %q in -%s", char, shorthands)
		return
	}

	if flag.ShorthandDeprecated != "" {
		fmt.Fprintf(f.Output(), "Flag shorthand -%c has been deprecated, %s\n", flag.Shorthand, flag.ShorthandDeprecated)
	}

	err = fn(flag, value)
	if err != nil {
		err = f.failf(err.Error())
	}
	return
}

func (f *FlagSet) parseShortArg(s string, args []string, fn parseFunc) (outArgs []string, err error) {
	outArgs = args
	shorthands := s[1:]

	// "shorthands" can be a series of shorthand letters of flags (e.g. "-vvv").
	for utf8.RuneCountInString(shorthands) > 0 {
		shorthands, outArgs, err = f.parseSingleShortArg(shorthands, args, fn)
		if err != nil {
			return
		}
	}

	return
}

func (f *FlagSet) parseArgs(args []string, fn parseFunc) (err error) {
	for len(args) > 0 {
		s := args[0]
		args = args[1:]
		if len(s) == 0 || s[0] != '-' || len(s) == 1 {
			if !f.interspersed {
				f.args = append(f.args, s)
				f.args = append(f.args, args...)
				return nil
			}
			f.args = append(f.args, s)
			continue
		}

		if s[1] == '-' {
			if len(s) == 2 { // "--" terminates the flags
				f.argsLenAtDash = len(f.args)
				f.args = append(f.args, args...)
				break
			}
			args, err = f.parseLongArg(s, args, fn)
		} else {
			args, err = f.parseShortArg(s, args, fn)
		}
		if err != nil {
			return
		}
	}
	return
}

func (f *FlagSet) parseAll(arguments []string, fn parseFunc) error {
	if f.addedGoFlagSets != nil {
		for _, goFlagSet := range f.addedGoFlagSets {
			if err := goFlagSet.Parse(nil); err != nil {
				return err
			}
		}
	}
	f.parsed = true

	if len(arguments) == 0 {
		return nil
	}

	f.args = make([]string, 0, len(arguments))

	err := f.parseArgs(arguments, fn)
	if err != nil {
		switch f.errorHandling {
		case ContinueOnError:
			return err
		case ExitOnError:
			if err == ErrHelp {
				os.Exit(0)
			}
			os.Exit(2)
		case PanicOnError:
			panic(err)
		}
	}
	return nil
}

// Parse parses flag definitions from the argument list, which should not
// include the command name.  Must be called after all flags in the FlagSet
// are defined and before flags are accessed by the program.
// The return value will be ErrHelp if -help was set but not defined.
func (f *FlagSet) Parse(arguments []string) error {
	set := func(flag *Flag, value string) error {
		return f.Set(flag.Name, value)
	}
	return f.parseAll(arguments, set)
}

type parseFunc func(flag *Flag, value string) error

// ParseAll parses flag definitions from the argument list, which should not
// include the command name. The arguments for fn are flag and value. Must be
// called after all flags in the FlagSet are defined and before flags are
// accessed by the program. The return value will be ErrHelp if -help was set
// but not defined.
func (f *FlagSet) ParseAll(arguments []string, fn func(flag *Flag, value string) error) error {
	return f.parseAll(arguments, fn)
}

// Parsed reports whether f.Parse has been called.
func (f *FlagSet) Parsed() bool {
	return f.parsed
}

// Parse parses the command-line flags from os.Args[1:].  Must be called
// after all flags are defined and before flags are accessed by the program.
func Parse() {
	// Ignore errors; CommandLine is set for ExitOnError.
	CommandLine.Parse(os.Args[1:])
}

// ParseAll parses the command-line flags from os.Args[1:] and called fn for each.
// The arguments for fn are flag and value. Must be called after all flags are
// defined and before flags are accessed by the program.
func ParseAll(fn func(flag *Flag, value string) error) {
	// Ignore errors; CommandLine is set for ExitOnError.
	CommandLine.ParseAll(os.Args[1:], fn)
}

// SetInterspersed sets whether to support interspersed option/non-option arguments.
func SetInterspersed(interspersed bool) {
	CommandLine.SetInterspersed(interspersed)
}

// Parsed returns true if the command-line flags have been parsed.
func Parsed() bool {
	return CommandLine.Parsed()
}

// CommandLine is the default set of command-line flags, parsed from os.Args.
var CommandLine = NewFlagSet(os.Args[0], ExitOnError)

// NewFlagSet returns a new, empty flag set with the specified name,
// error handling property and SortFlags set to true.
func NewFlagSet(name string, errorHandling ErrorHandling) *FlagSet {
	f := &FlagSet{
		name:          name,
		errorHandling: errorHandling,
		argsLenAtDash: -1,
		interspersed:  true,
		SortFlags:     true,
	}
	return f
}

// SetInterspersed sets whether to support interspersed option/non-option arguments.
func (f *FlagSet) SetInterspersed(interspersed bool) {
	f.interspersed = interspersed
}

// Init sets the name and error handling property for a flag set.
// By default, the zero FlagSet uses an empty name and the
// ContinueOnError error handling policy.
func (f *FlagSet) Init(name string, errorHandling ErrorHandling) {
	f.name = name
	f.errorHandling = errorHandling
	f.argsLenAtDash = -1
}
