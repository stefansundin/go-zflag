# zflag

***This is a fork of [zulucmd/zflag](https://github.com/zulucmd/zflag), which is a fork of [cornfeedhobo/pflag](https://github.com/cornfeedhobo/pflag), which in turn is a fork of [spf13/pflag](https://github.com/spf13/pflag). All forks were created due to poor maintenance. This is madness.***

[![GoDoc](https://godoc.org/github.com/stefansundin/go-zflag?status.svg)](https://godoc.org/github.com/stefansundin/go-zflag)
[![Go Report Card](https://goreportcard.com/badge/github.com/stefansundin/go-zflag)](https://goreportcard.com/report/github.com/stefansundin/go-zflag)
[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/stefansundin/go-zflag?sort=semver)](https://github.com/stefansundin/go-zflag/releases)
[![Build Status](https://github.com/stefansundin/go-zflag/actions/workflows/validate.yml/badge.svg)](https://github.com/stefansundin/go-zflag/actions/workflows/validate.yml)

* [Installation](#installation)
* [Supported Syntax](#supported-syntax)
* [Documentation](#documentation)
  * [Create a flag set](#create-a-flag-set)
  * [Set a custom default for flags passed without values](#set-a-custom-default-for-flags-passed-without-values)
  * [Mutating or "Normalizing" Flag names](#mutating-or-normalizing-flag-names)
  * [Deprecating a flag or its shorthand](#deprecating-a-flag-or-its-shorthand)
  * [Hidden flags](#hidden-flags)
  * [Disable sorting of flags](#disable-sorting-of-flags)
  * [Supporting Go flags when using zflag](#supporting-go-flags-when-using-zflag)
  * [Shorthand flags](#shorthand-flags)
  * [Shorthand-only flags](#shorthand-only-flags)
  * [Unknown flags](#unknown-flags)
  * [Custom flag types in usage](#custom-flag-types-in-usage)
  * [Disable printing a flag's default value](#disable-printing-a-flags-default-value)
  * [Disable built-in help flags](#disable-built-in-help-flags)

## Installation

zflag is available using the standard `go get` command.

Install by running:

```bash
go get github.com/stefansundin/go-zflag
```

You can import the package with the `flag` name using:

```go
import (
	flag "github.com/stefansundin/go-zflag"
)
```

## Supported Syntax

```plain
--flag    // boolean flags, or flags with no option default values
--flag x  // only on flags without a default value
--flag=x
```

Unlike the flag package, a single dash before an option means something
different than a double dash. Single dashes signify a series of shorthand
letters for flags. All but the last shorthand letter must be boolean flags
or a flag with a default value

```plain
// boolean or flags where the 'no option default value' is set
-f
-f=true
-abc
but
-b true is INVALID

// non-boolean and flags without a 'no option default value'
-n 1234
-n=1234
-n1234

// mixed
-abcs "hello"
-absd="hello"
-abcs1234
```

Slice flags can be specified multiple times, or specified with an equal sign and csv.

```plain
--sliceVal one --sliceVal=two
--sliceVal=one,two
```

Integer flags accept 1234, 0664, 0x1234 and may be negative.
Boolean flags (in their long form) accept 1, 0, t, f, true, false,
TRUE, FALSE, True, False.
Duration flags accept any input valid for time.ParseDuration.

Flag parsing stops after the terminator `--`. Unlike the flag package,
flags can be interspersed with arguments anywhere on the command line
before this terminator.

## Documentation

You can see the full reference documentation of the zflag package
[at godoc.org](http://godoc.org/github.com/stefansundin/go-zflag), querying with
[`go doc`](https://golang.org/cmd/doc/), or through go's standard documentation
system by running `godoc -http=:6060` and browsing to
[http://localhost:6060/pkg/github.com/stefansundin/go-zflag](http://localhost:6060/pkg/github.com/stefansundin/go-zflag)
after installation.

### Create a flag set

The default set of command-line flags is controlled by top-level functions. The FlagSet type allows one to define independent sets of flags, such as to implement subcommands in a command-line interface. The methods of FlagSet are analogous to the top-level functions for the command-line flag set.

```go
flagSet := flag.NewFlagSet("myapp", flag.ExitOnError)
```

### Set a custom default for flags passed without values

If a flag has a NoOptDefVal and the flag is set on the command line
without an option, the flag will be set to the NoOptDefVal.

**Example**:

```go
var ip = flag.Int("flagname", 1234, "help message", flag.OptShorthand('f'), flag.OptNoOptDefVal("4321"))
```

**Results**:

| Parsed Arguments | Resulting Value |
| -------------    | -------------   |
| --flagname=1357  | ip=1357         |
| --flagname       | ip=4321         |
| [nothing]        | ip=1234         |

### Mutating or "Normalizing" Flag names

It is possible to set a custom flag name 'normalization function.' It allows
flag names to be mutated both when created in the code and when used on the
command line to some 'normalized' form. The 'normalized' form is used for
comparison. Two examples of using the custom normalization func follow.

**Example #1**: You want -, _, and . in flags to compare the same, i.e. `--my-flag` == `--my_flag` == `--my.flag`

```go
func wordSepNormalizeFunc(f *flag.FlagSet, name string) flag.NormalizedName {
	from := []string{"-", "_"}
	to := "."
	for _, sep := range from {
		name = strings.Replace(name, sep, to, -1)
	}
	return flag.NormalizedName(name)
}

flagSet.SetNormalizeFunc(wordSepNormalizeFunc)
```

**Example #2**: You want to alias two flags, i.e. `--old-flag-name` == `--new-flag-name`

```go
func aliasNormalizeFunc(f *flag.FlagSet, name string) flag.NormalizedName {
	switch name {
	case "old-flag-name":
		name = "new-flag-name"
		break
	}
	return flag.NormalizedName(name)
}

flagSet.SetNormalizeFunc(aliasNormalizeFunc)
```

### Deprecating a flag or its shorthand

It is possible to deprecate a flag, or just its shorthand. Deprecating a
flag/shorthand hides it from help text and prints a usage message when the
deprecated flag/shorthand is used.

**Example #1**: You want to deprecate a flag named "badflag" as well as
inform the users what flag they should use instead.

```go
// deprecate a flag by specifying its name and a usage message
flag.Bool("badflag", false, "this does something", flag.OptDeprecated("please use --good-flag instead"))
```

This hides "badflag" from help text, and prints
`Flag --badflag has been deprecated, please use --good-flag instead`
when "badflag" is used.

**Example #2**: You want to keep a flag name "noshorthandflag" but deprecate
it's shortname "n".

```go
// deprecate a flag shorthand by specifying its flag name and a usage message
flag.Bool("noshorthandflag", false, "this does something", flag.OptShorthand('n'), flag.OptShorthandDeprecated("please use --noshorthandflag only"))
```

This hides the shortname "n" from help text, and prints
`Flag shorthand -n has been deprecated, please use --noshorthandflag only`
when the shorthand `n` is used.

Note that usage message is essential here, and it should not be empty. If it is empty, it will panic at runtime.

### Hidden flags

It is possible to mark a flag as hidden, meaning it will still function as
normal, however will not show up in usage/help text.

**Example**: You have a flag named `secretFlag` that you need for internal use
only and don't want it showing up in help text, or for its usage text to be available.

```go
// hide a flag by specifying its name
flag.Bool("secretFlag", false, "this does something", flag.OptHidden())
```

### Disable sorting of flags

It is possible to disable sorting of flags for help and usage message.

**Example**:

```go
flag.Bool("verbose", false, "verbose output", flag.OptShorthand('v'))
flag.String("coolflag", "yeaah", "it's really cool flag")
flag.Int("usefulflag", 777, "sometimes it's very useful")
flag.SortFlags = false
flag.PrintDefaults()
```

**Output**:

```plain
  -v, --verbose           verbose output
      --coolflag string   it's really cool flag (default "yeaah")
      --usefulflag int    sometimes it's very useful (default 777)
```

### Supporting Go flags when using zflag

In order to support flags defined using Go's `flag` package, they must be added
to the `zflag` flagset. This is usually necessary to support flags defined by
third-party dependencies (e.g. `golang/glog`).

**Example**: You want to add the Go flags to the `CommandLine` flagset

```go
import (
	goflag "flag"
	flag "github.com/stefansundin/go-zflag"
)

var ip *int = flag.Int("flagname", 1234, "help message for flagname")

func main() {
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	flag.Parse()
}
```

### Shorthand flags

A flag supporting both long and short formats can be created with the new
`OptShorthand` argument:

```go
flag.Bool("toggle", false, "toggle help message", flag.OptShorthand('t'))
```

### Shorthand-only flags

A shorthand-only flag can be created with `OptShorthandOnly`:

```go
flag.String("value", "", "value help message", flag.OptShorthandOnly('l'))
```

This flag can be looked up using its long name, but will only be parsed when
the short form is passed.

### Unknown flags

Normally zflag will error when an unknown flag is passed, but it's also possible
to disable that using `FlagSet.ParseErrorsAllowlist.UnknownFlags`:

```go
flagSet.ParseErrorsAllowlist.UnknownFlags = true
flag.Parse()
```

These can then be obtained as a slice of strings using `FlagSet.GetUnknownFlags()`.

### Custom flag types in usage

There are two methods to set a custom type to be printed in the usage.

First, it's possible to set explicitly with `UsageType`:

```go
flag.String("character", "", "character name", flag.OptUsageType("enum"))
```

Output:

```plain
  --character enum   character name (default "")
```

Alternatively, it's possible to include backticks around a single word in the
usage string, which will be extracted and printed with the usage:

```go
flag.String("character", "", "`character` name")
```

Output:

```plain
  --character character   character name (default "")
```

_Note: This unquoting behavior can be disabled with `Flag.DisableUnquoteUsage`, or `flag.OptDisableUnquoteUsage`_.

### Customizing flag usages

You can customize the flag usages by overriding the `FlagSet.FlagUsageFormatter` field
with a struct that implements [`FlagUsageFormatter`](./formatter.go).
You can re-implement the whole interface or embed the `DefaultFlagUsageFormatter` to only
re-implement some some functions.

For example:

```go
type myFlagFormatter struct {
  flag.DefaultFlagUsageFormatter
}

func (f myFlagFormatter) Name(flag *flag.Flag) string {
  return "--not-hello"
}


flagSet := flag.NewFlagSet("myapp", flag.ExitOnError)
flagSet.String("hello", "", "myusage")

flagSet.FlagUsageFormatter = myFlagFormatter{}
flagSet.PrintDefaults()
```

Which will print:

```
--not-hello string   myusage
```

### Disable printing a flag's default value

The printing of a flag's default value can be suppressed with `Flag.DisablePrintDefault`.

**Example**:

```go
flag.Int("in", -1, "help message", flag.OptDisablePrintDefault(true))
```

**Output**:

```plain
  --in int   help message
```

### Disable built-in help flags

Normally zflag will handle `--help` and `-h` when the flags aren't explicitly defined.

If for some reason there is a need to capture the error returned in this condition, it
is possible to disable this built-in handling.

```go
flagSet.DisableBuiltinHelp = true
```
