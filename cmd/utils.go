package cmd

import (
	"fmt"
	"os"
)

func fatalln(s string) {
	fmt.Fprintln(os.Stderr, s)
	os.Exit(1)
}

func fatalf(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}
