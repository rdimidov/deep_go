package main

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type MultiError struct {
	Errors []error
}

func (e *MultiError) Error() string {
	if len(e.Errors) == 0 {
		return ""
	}

	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("%d errors occured:\n", len(e.Errors)))

	for _, err := range e.Errors {
		b.WriteString(fmt.Sprintf("\t* %s", err.Error()))
	}

	b.WriteString("\n")
	return b.String()
}

func Append(err error, errs ...error) *MultiError {
	var (
		errList []error
		me      *MultiError
	)

	if err != nil && errors.As(err, &me) {
		errList = append(errList, me.Errors...)
	} else if err != nil {
		errList = append(errList, err)
	}

	for _, e := range errs {
		if e == nil {
			continue
		}

		var m *MultiError
		if errors.As(e, &m) {
			errList = append(errList, m.Errors...)
		} else {
			errList = append(errList, e)
		}
	}

	if len(errList) == 0 {
		return nil
	}

	return &MultiError{Errors: errList}
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))

	expectedMessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)
}
