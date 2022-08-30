package validator_test

import (
	"github.com/liviu-moraru/snippetbox/internal/validator"
	"testing"
)

func TestValidator_Permitted(t *testing.T) {
	permitted := validator.Permitted(10, 11, 20, 30)
	if permitted {
		t.Errorf("Not permitted int value, validated")
	}
	permitted = validator.Permitted("abc", "x", "y", "c")
	if permitted {
		t.Errorf("Not permitted string value, validated")
	}
}

func TestValidator_AddFieldError(t *testing.T) {
	// Add an error to a initial empty validator
	v := validator.Validator{}

	v.AddFieldError("content", "My content")

	if v.FieldErrors == nil || v.FieldErrors["content"] != "My content" {
		t.Error("Error not added to an empty validator")
	}

	// Same key, the error is not overwritten
	v = validator.Validator{}
	v.AddFieldError("content", "My first content")

	v.AddFieldError("content", "My second content")
	if v.FieldErrors == nil || v.FieldErrors["content"] != "My first content" {
		t.Error("Error adding an existing key")
	}
}
