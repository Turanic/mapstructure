package mapstructure

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// Error implements the error interface and can represents multiple
// errors that occur in the course of a single decode.
type Error struct {
	Errors []string
}

func (e *Error) Error() string {
	points := make([]string, len(e.Errors))
	for i, err := range e.Errors {
		points[i] = fmt.Sprintf("* %s", err)
	}

	sort.Strings(points)
	return fmt.Sprintf(
		"%d error(s) decoding:\n\n%s",
		len(e.Errors), strings.Join(points, "\n"))
}

type TypeError struct {
	FieldName string
	FieldType reflect.Type
	FromType  reflect.Type
	Value     interface{}
}

func (e *TypeError) Error() string {
	return fmt.Sprintf("'%s' expected type '%s', got unconvertible type '%s', value: '%v'",
  e.FieldName, e.FieldType, e.FromType, e.Value)
}

type DecodeNumberError struct {
	FieldName string
  NumberKind reflect.Kind
	Err       error
}

func (e *DecodeNumberError) Error() string {
  return fmt.Sprintf("cannot parse '%s' as %s: %s", e.FieldName, e.NumberKind, e.Err.Error())
}

// WrappedErrors implements the errwrap.Wrapper interface to make this
// return value more useful with the errwrap and go-multierror libraries.
func (e *Error) WrappedErrors() []error {
	if e == nil {
		return nil
	}

	result := make([]error, len(e.Errors))
	for i, e := range e.Errors {
		result[i] = errors.New(e)
	}

	return result
}

func appendErrors(errors []string, err error) []string {
	switch e := err.(type) {
	case *Error:
		return append(errors, e.Errors...)
	default:
		return append(errors, e.Error())
	}
}
