package multierror

import (
	"errors"
	"testing"
)

func TestPrefix_Error(t *testing.T) {
	original := &Error{
		Errs: []error{errors.New("foo")},
	}

	result := Prefix(original, "bar")
	if result.(*Error).Errs[0].Error() != "bar foo" {
		t.Fatalf("bad: %s", result)
	}
}

func TestPrefix_NilError(t *testing.T) {
	var err error
	result := Prefix(err, "bar")
	if result != nil {
		t.Fatalf("bad: %#v", result)
	}
}

func TestPrefix_NonError(t *testing.T) {
	original := errors.New("foo")
	result := Prefix(original, "bar")
	if result == nil {
		t.Fatal("error result was nil")
	}
	if result.Error() != "bar foo" {
		t.Fatalf("bad: %s", result)
	}
}
