package isbn

import (
	"errors"
	"fmt"
	"reflect"
)

var ErrTodo = errors.New("todo")

type invalidLengthError struct{ len int }

func (err invalidLengthError) Error() string {
	return fmt.Sprintf("invalid ISBN length %d", err.len)
}

type invalidTypeError struct{ value reflect.Type }

func (err invalidTypeError) Error() string {
	return fmt.Sprintf("failed to scan type %+v for value", err.value)
}
