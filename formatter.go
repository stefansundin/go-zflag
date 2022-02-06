// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pflag

import "fmt"

type FlagUsageFormatter interface {
	Name(*Flag) string
	Usage(*Flag, string) string
	UsageVarName(*Flag, string) string
	DefaultValue(*Flag) string
	NoOptDefValue(*Flag) string
	Deprecated(*Flag) string
}

type DefaultFlagUsageFormatter struct{}

var _ FlagUsageFormatter = (*DefaultFlagUsageFormatter)(nil)

func (d DefaultFlagUsageFormatter) Name(flag *Flag) string {
	name := "  "
	if flag.Shorthand != 0 && flag.ShorthandDeprecated == "" {
		name += fmt.Sprintf("-%c", flag.Shorthand)
		if !flag.ShorthandOnly {
			name += ", "
		}
	} else {
		name += "    "
	}
	name += fmt.Sprintf("--%s", flag.Name)

	return name
}

func (d DefaultFlagUsageFormatter) Usage(flag *Flag, s string) string {
	return s
}

func (d DefaultFlagUsageFormatter) UsageVarName(flag *Flag, s string) string {
	return s
}

func (d DefaultFlagUsageFormatter) DefaultValue(flag *Flag) string {
	if v, ok := flag.Value.(Typed); ok && v.Type() == "string" {
		return fmt.Sprintf(" (default %q)", flag.DefValue)
	}

	return fmt.Sprintf(" (default %s)", flag.DefValue)
}

func (d DefaultFlagUsageFormatter) NoOptDefValue(flag *Flag) string {
	if v, ok := flag.Value.(Typed); ok {
		switch v.Type() {
		case "string":
			return fmt.Sprintf("[=\"%s\"]", flag.NoOptDefVal)
		case "bool":
			if flag.NoOptDefVal != "true" {
				return fmt.Sprintf("[=%s]", flag.NoOptDefVal)
			}
		case "count":
			if flag.NoOptDefVal != "+1" {
				return fmt.Sprintf("[=%s]", flag.NoOptDefVal)
			}
		default:
			return fmt.Sprintf("[=%s]", flag.NoOptDefVal)
		}
		return ""
	}

	return fmt.Sprintf("[=%s]", flag.NoOptDefVal)
}

func (d DefaultFlagUsageFormatter) Deprecated(flag *Flag) string {
	return fmt.Sprintf(" (DEPRECATED: %s)", flag.Deprecated)
}
