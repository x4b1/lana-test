package errors_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xabi93/lana-test/pkg/errors"
)

func TestNewCreatesCorrectErrorType(t *testing.T) {
	cases := map[string]struct {
		newFunc   func(format string, args ...interface{}) error
		checkFunc func(err error) bool
	}{
		"wrongInput": {
			newFunc:   errors.NewWrongInput,
			checkFunc: errors.IsWrongInput,
		},
		"notFound": {
			newFunc:   errors.NewNotFound,
			checkFunc: errors.IsNotFound,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			require.True(t, c.checkFunc(c.newFunc("some error")))
		})
	}
}

func TestWrapsWithCorrectErrorType(t *testing.T) {
	cases := map[string]struct {
		newFunc   func(err error, format string, args ...interface{}) error
		checkFunc func(err error) bool
	}{
		"notFound": {
			newFunc:   errors.WrapNotFound,
			checkFunc: errors.IsNotFound,
		},
	}

	someError := fmt.Errorf("Some error")

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			require.True(t, c.checkFunc(c.newFunc(someError, "wrap error")))
		})
	}
}
