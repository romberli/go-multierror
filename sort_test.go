package multierror

import (
	"errors"
	"reflect"
	"sort"
	"testing"
)

func TestSortSingle(t *testing.T) {
	errFoo := errors.New("foo")

	expected := []error{
		errFoo,
	}

	err := &Error{
		Errs: []error{
			errFoo,
		},
	}

	sort.Sort(err)
	if !reflect.DeepEqual(err.Errs, expected) {
		t.Fatalf("bad: %#v", err)
	}
}

func TestSortMultiple(t *testing.T) {
	errBar := errors.New("bar")
	errBaz := errors.New("baz")
	errFoo := errors.New("foo")

	expected := []error{
		errBar,
		errBaz,
		errFoo,
	}

	err := &Error{
		Errs: []error{
			errFoo,
			errBar,
			errBaz,
		},
	}

	sort.Sort(err)
	if !reflect.DeepEqual(err.Errs, expected) {
		t.Fatalf("bad: %#v", err)
	}
}
