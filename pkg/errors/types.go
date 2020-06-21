package errors

import (
	"errors"
	"fmt"
)

type wrongInput struct {
	error
}

func NewWrongInput(format string, args ...interface{}) error {
	return &wrongInput{fmt.Errorf(format, args...)}
}

func IsWrongInput(err error) bool {
	var target *wrongInput

	return errors.As(err, &target)
}

type notFound struct {
	error
}

func NewNotFound(format string, args ...interface{}) error {
	return &notFound{fmt.Errorf(format, args...)}
}

func WrapNotFound(err error, format string, args ...interface{}) error {
	args = append(args, err)

	return &notFound{fmt.Errorf(format+": %w", args...)}
}

func IsNotFound(err error) bool {
	var target *notFound

	return errors.As(err, &target)
}
