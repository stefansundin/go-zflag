// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"fmt"
)

type ApplyOptFunc func(c *Flag) error

func (f ApplyOptFunc) apply(c *Flag) error {
	return f(c)
}

func applyFlagOptions(c *Flag, options ...Opt) error {
	for _, o := range options {
		if err := o.apply(c); err != nil {
			return err
		}
	}
	return nil
}

type Opt interface {
	apply(*Flag) error
}

type optShorthandImpl struct{ shorthand rune }

func (o optShorthandImpl) apply(c *Flag) error { c.Shorthand = o.shorthand; return nil }

// OptShorthand one-letter abbreviated flag
func OptShorthand(o rune) Opt { return optShorthandImpl{shorthand: o} }

// OptShorthandStr one-letter abbreviated flag
func OptShorthandStr(shorthand string) Opt {
	r, err := shorthandStrToRune(shorthand)
	if err != nil {
		panic(err)
	}

	return optShorthandImpl{shorthand: r}
}

type optShorthandOnlyImpl struct{}

func (o optShorthandOnlyImpl) apply(c *Flag) error { c.ShorthandOnly = true; return nil }

// OptShorthandOnly If the user set only the shorthand
func OptShorthandOnly() Opt { return optShorthandOnlyImpl{} }

type optUsageTypeImpl struct{ usageType string }

func (o optUsageTypeImpl) apply(c *Flag) error { c.UsageType = o.usageType; return nil }

// OptUsageType flag type displayed in the help message
func OptUsageType(usageType string) Opt { return optUsageTypeImpl{usageType: usageType} }

type optDisableUnquoteUsageImpl struct{}

func (o optDisableUnquoteUsageImpl) apply(c *Flag) error { c.DisableUnquoteUsage = true; return nil }

// OptDisableUnquoteUsage disable unquoting and extraction of type from usage
func OptDisableUnquoteUsage() Opt { return optDisableUnquoteUsageImpl{} }

type optDisablePrintDefaultImpl struct{}

func (o optDisablePrintDefaultImpl) apply(c *Flag) error { c.DisablePrintDefault = true; return nil }

// OptDisablePrintDefault toggle printing of the default value in usage message
func OptDisablePrintDefault(o bool) Opt { return optDisablePrintDefaultImpl{} }

type optDefValueImpl struct{ defValue string }

func (o optDefValueImpl) apply(c *Flag) error { c.DefValue = o.defValue; return nil }

// OptDefValue default value (as text); for usage message
func OptDefValue(defValue string) Opt { return optDefValueImpl{defValue: defValue} }

type optNoOptDefValImpl struct{ noOptDefVal string }

func (o optNoOptDefValImpl) apply(c *Flag) error { c.NoOptDefVal = o.noOptDefVal; return nil }

// OptNoOptDefVal default value (as text); if the flag is on the command line without any options
func OptNoOptDefVal(noOptDefVal string) Opt { return optNoOptDefValImpl{noOptDefVal: noOptDefVal} }

type optDeprecatedImpl struct{ msg string }

func (o optDeprecatedImpl) apply(c *Flag) error {
	if o.msg == "" {
		return fmt.Errorf("deprecated message for flag %q must be set", c.Name)
	}

	c.Deprecated = o.msg
	return OptHidden().apply(c)
}

// OptDeprecated indicated that a flag is deprecated in your program. It will
// continue to function but will not show up in help or usage messages. Using
// this flag will also print the given usageMessage.
func OptDeprecated(msg string) Opt { return optDeprecatedImpl{msg: msg} }

type optHiddenImpl struct{}

func (o optHiddenImpl) apply(c *Flag) error { c.Hidden = true; return nil }

// OptHidden used by zulu.Command to allow flags to be hidden from help/usage text
func OptHidden() Opt { return optHiddenImpl{} }

type optShorthandDeprecatedImpl struct{ msg string }

func (o optShorthandDeprecatedImpl) apply(c *Flag) error {
	if o.msg == "" {
		return fmt.Errorf("shorthand deprecated message for flag %q must be set", c.Name)
	}

	c.ShorthandDeprecated = o.msg
	return nil
}

// OptShorthandDeprecated If the shorthand of this flag is deprecated, this string is the new or now thing to use
func OptShorthandDeprecated(msg string) Opt { return optShorthandDeprecatedImpl{msg: msg} }

type optGroupImpl struct{ group string }

func (o optGroupImpl) apply(c *Flag) error { c.Group = o.group; return nil }

// OptGroup flag group
func OptGroup(group string) Opt { return optGroupImpl{group: group} }

type optAnnotationImpl struct {
	key   string
	value []string
}

func (o optAnnotationImpl) apply(c *Flag) error {
	return c.SetAnnotation(o.key, o.value)
}

// OptAnnotation Use it to annotate this specific flag for your application
func OptAnnotation(key string, value []string) Opt {
	return optAnnotationImpl{
		key:   key,
		value: value,
	}
}
