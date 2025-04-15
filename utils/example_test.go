package utils

import (
	"context"
	"errors"
	"fmt"
	"os"
)

func ExampleWithLog() {
	f := WithLog[string](os.Stdout, os.Stdout, "test", func(msg string) error {
		fmt.Println(msg)
		return errors.New("test error")
	})

	f(context.Background(), "tick")

	// Output:
	// Calling test
	// tick
	// Execution of test failed with error: test error
}

func ExampleAdapt() {
	f := func(s string) error {
		fmt.Println(s)
		return errors.New("error")
	}
	var adaptedF func(context.Context, string) error

	adaptedF = Adapt[string](f)
	err := adaptedF(context.Background(), "hello")
	fmt.Println(err)

	// Output:
	// hello
	// error
}
