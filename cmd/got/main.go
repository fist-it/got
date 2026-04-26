package main

import (
	"fmt"
	"os"

	"github.com/fist-it/got/internal/porcelain"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: got <command>")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "init":
		err := porcelain.Init(".")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			os.Exit(1)
		}
	}
}
