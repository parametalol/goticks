package utils

import (
	"context"
	"errors"
	"fmt"
	"os"
)

func ExampleWithLog() {
	f := Log[string](os.Stdout, os.Stdout, "test", func(msg string) error {
		fmt.Println(msg)
		return errors.New("test error")
	})

	fmt.Println("Error:", f(context.Background(), "tick"))

	// Output:
	// Calling test
	// tick
	// Execution of test failed with error: test error
	// Error: test error
}

func ExampleAdapt() {
	f := func(s string) error {
		fmt.Println(s)
		return errors.New("error")
	}
	var adaptedF = Adapt[string](f) // func(context.Context, string) error
	err := adaptedF(context.Background(), "hello")
	fmt.Println(err)

	// Output:
	// hello
	// error
}
