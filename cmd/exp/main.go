package main

import (
	"errors"
	"fmt"
)

func main() {
	err := B()
	// TODO: Determind if the `err` variable is an `ErrNotFound`
	fmt.Printf("The err contains ErrNotFound: %v\n", errors.Is(err, ErrNotFound))
}

// It is common for packages like database/sql to return
// an error that is predefined like this one.
var ErrNotFound = errors.New("not found")

func A() error {
	return ErrNotFound
}

func B() error {
	err := A()
	return fmt.Errorf("b: %w", err)
}
