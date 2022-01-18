package multierror

import (
	"errors"
	"testing"
)

func TestAppend_Error(t *testing.T) {
	original := &Error{
		Errs: []error{errors.New("foo")},
	}

	result := Append(original, errors.New("bar"))
	if len(result.Errs) != 2 {
		t.Fatalf("wrong len: %d", len(result.Errs))
	}

	original = &Error{}
	result = Append(original, errors.New("bar"))
	if len(result.Errs) != 1 {
		t.Fatalf("wrong len: %d", len(result.Errs))
	}

	// Test when a typed nil is passed
	var e *Error
	result = Append(e, errors.New("baz"))
	if len(result.Errs) != 1 {
		t.Fatalf("wrong len: %d", len(result.Errs))
	}

	// Test flattening
	original = &Error{
		Errs: []error{errors.New("foo")},
	}

	result = Append(original, Append(nil, errors.New("foo"), errors.New("bar")))
	if len(result.Errs) != 3 {
		t.Fatalf("wrong len: %d", len(result.Errs))
	}
}

func TestAppend_NilError(t *testing.T) {
	var err error
	result := Append(err, errors.New("bar"))
	if len(result.Errs) != 1 {
		t.Fatalf("wrong len: %d", len(result.Errs))
	}
}

func TestAppend_NilErrorArg(t *testing.T) {
	var err error
	var nilErr *Error
	result := Append(err, nilErr)
	if len(result.Errs) != 0 {
		t.Fatalf("wrong len: %d", len(result.Errs))
	}
}

func TestAppend_NilErrorIfaceArg(t *testing.T) {
	var err error
	var nilErr error
	result := Append(err, nilErr)
	if len(result.Errs) != 0 {
		t.Fatalf("wrong len: %d", len(result.Errs))
	}
}

func TestAppend_NonError(t *testing.T) {
	original := errors.New("foo")
	result := Append(original, errors.New("bar"))
	if len(result.Errs) != 2 {
		t.Fatalf("wrong len: %d", len(result.Errs))
	}
}

func TestAppend_NonError_Error(t *testing.T) {
	original := errors.New("foo")
	result := Append(original, Append(nil, errors.New("bar")))
	if len(result.Errs) != 2 {
		t.Fatalf("wrong len: %d", len(result.Errs))
	}
}
