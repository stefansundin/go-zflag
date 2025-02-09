// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
)

func setUpDSFlagSet(dsp *[]time.Duration) *FlagSet {
	f := NewFlagSet("test", ContinueOnError)
	f.DurationSliceVar(dsp, "ds", []time.Duration{}, "Command separated list!")
	return f
}

func setUpDSFlagSetWithDefault(dsp *[]time.Duration) *FlagSet {
	f := NewFlagSet("test", ContinueOnError)
	f.DurationSliceVar(dsp, "ds", []time.Duration{0, 1}, "Command separated list!")
	return f
}

func TestDSValueImplementsGetter(t *testing.T) {
	var v Value = new(durationSliceValue)

	if _, ok := v.(Getter); !ok {
		t.Fatalf("%T should implement the Getter interface", v)
	}
}

func TestEmptyDS(t *testing.T) {
	var ds []time.Duration
	f := setUpDSFlagSet(&ds)
	err := f.Parse([]string{})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	getDS, err := f.GetDurationSlice("ds")
	if err != nil {
		t.Fatal("got an error from GetDurationSlice():", err)
	}

	if len(getDS) != 0 {
		t.Fatalf("got ds %v with len=%d but expected length=0", getDS, len(getDS))
	}

	getDS2, err := f.Get("ds")
	if err != nil {
		t.Fatal("got an error from Get():", err)
	}

	if !reflect.DeepEqual(getDS, getDS2) {
		t.Fatalf("expected %v with type %T but got %v with type %T ", getDS, getDS, getDS2, getDS2)
	}
}

func TestDS(t *testing.T) {
	var ds []time.Duration
	f := setUpDSFlagSet(&ds)

	vals := []string{"1ns", "2ms", "3m", "4h"}
	arg := fmt.Sprintf("--ds=%s", strings.Join(vals, ","))
	err := f.Parse([]string{arg})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range ds {
		d, err := time.ParseDuration(vals[i])
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		if d != v {
			t.Fatalf("expected ds[%d] to be %s but got: %d", i, vals[i], v)
		}
	}
	getDS, err := f.GetDurationSlice("ds")
	if err != nil {
		t.Fatalf("got error: %v", err)
	}
	for i, v := range getDS {
		d, err := time.ParseDuration(vals[i])
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		if d != v {
			t.Fatalf("expected ds[%d] to be %s but got: %d from GetDurationSlice", i, vals[i], v)
		}
	}

	getDS2, err := f.Get("ds")
	if err != nil {
		t.Fatal("got an error from Get():", err)
	}

	if !reflect.DeepEqual(getDS2, getDS) {
		t.Fatalf("got %v with type %T but expected %v with type %T", getDS2, getDS2, getDS, getDS)
	}
}

func TestDSDefault(t *testing.T) {
	var ds []time.Duration
	f := setUpDSFlagSetWithDefault(&ds)

	vals := []string{"0s", "1ns"}

	err := f.Parse([]string{})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range ds {
		d, err := time.ParseDuration(vals[i])
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		if d != v {
			t.Fatalf("expected ds[%d] to be %d but got: %d", i, d, v)
		}
	}

	getDS, err := f.GetDurationSlice("ds")
	if err != nil {
		t.Fatal("got an error from GetDurationSlice():", err)
	}
	for i, v := range getDS {
		d, err := time.ParseDuration(vals[i])
		if err != nil {
			t.Fatal("got an error from GetDurationSlice():", err)
		}
		if d != v {
			t.Fatalf("expected ds[%d] to be %d from GetDurationSlice but got: %d", i, d, v)
		}
	}
}

func TestDSWithDefault(t *testing.T) {
	var ds []time.Duration
	f := setUpDSFlagSetWithDefault(&ds)

	vals := []string{"1ns", "2ns"}
	arg := fmt.Sprintf("--ds=%s", strings.Join(vals, ","))
	err := f.Parse([]string{arg})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range ds {
		d, err := time.ParseDuration(vals[i])
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		if d != v {
			t.Fatalf("expected ds[%d] to be %d but got: %d", i, d, v)
		}
	}

	getDS, err := f.GetDurationSlice("ds")
	if err != nil {
		t.Fatal("got an error from GetDurationSlice():", err)
	}
	for i, v := range getDS {
		d, err := time.ParseDuration(vals[i])
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		if d != v {
			t.Fatalf("expected ds[%d] to be %d from GetDurationSlice but got: %d", i, d, v)
		}
	}
}

func TestDSAsSliceValue(t *testing.T) {
	var ds []time.Duration
	f := setUpDSFlagSet(&ds)

	in := []string{"1ns", "2ns"}
	argfmt := "--ds=%s"
	arg1 := fmt.Sprintf(argfmt, in[0])
	arg2 := fmt.Sprintf(argfmt, in[1])
	err := f.Parse([]string{arg1, arg2})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	f.VisitAll(func(f *Flag) {
		if val, ok := f.Value.(SliceValue); ok {
			_ = val.Replace([]string{"3ns"})
		}
	})
	if len(ds) != 1 || ds[0] != time.Duration(3) {
		t.Fatalf("Expected ss to be overwritten with '3ns', but got: %v", ds)
	}
}

func TestDSCalledTwice(t *testing.T) {
	var ds []time.Duration
	f := setUpDSFlagSet(&ds)

	in := []string{"1ns,2ns", "3ns"}
	expected := []time.Duration{1, 2, 3}
	argfmt := "--ds=%s"
	arg1 := fmt.Sprintf(argfmt, in[0])
	arg2 := fmt.Sprintf(argfmt, in[1])
	err := f.Parse([]string{arg1, arg2})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range ds {
		if expected[i] != v {
			t.Fatalf("expected ds[%d] to be %d but got: %d", i, expected[i], v)
		}
	}
}
