package main

import (
	"fmt"
	"io"
	"os"

	"github.com/barthollomew/check-this.nvim/analyzer/internal/cli"
)

func main() {
	stdin, _ := io.ReadAll(os.Stdin)
	code, err := cli.Run(os.Args[1:], stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	os.Exit(code)
}
