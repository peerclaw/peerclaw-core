package errors

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	e := New(CodeNotFound, "agent not found")
	if e.Code != CodeNotFound {
		t.Errorf("Code = %q, want %q", e.Code, CodeNotFound)
	}
	if e.Message != "agent not found" {
		t.Errorf("Message = %q", e.Message)
	}
	if e.Unwrap() != nil {
		t.Error("expected nil cause")
	}
}

func TestWrap(t *testing.T) {
	cause := fmt.Errorf("connection refused")
	e := Wrap(CodeInternal, "database error", cause)
	if e.Unwrap() != cause {
		t.Error("Unwrap should return original cause")
	}
	if got := e.Error(); got != "internal_error: database error: connection refused" {
		t.Errorf("Error() = %q", got)
	}
}

func TestIs(t *testing.T) {
	e := New(CodeValidation, "bad input")
	if !Is(e, CodeValidation) {
		t.Error("Is should match CodeValidation")
	}
	if Is(e, CodeNotFound) {
		t.Error("Is should not match CodeNotFound")
	}

	// Wrapped error.
	wrapped := fmt.Errorf("wrap: %w", e)
	if !Is(wrapped, CodeValidation) {
		t.Error("Is should match through wrapping")
	}

	// Non-Error type.
	if Is(fmt.Errorf("plain"), CodeInternal) {
		t.Error("Is should return false for non-Error types")
	}
}
