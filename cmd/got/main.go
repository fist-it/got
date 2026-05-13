package main

import (
	"encoding/hex"
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
	case "hash-object":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "usage: got hash-object [-w] <file>")
			os.Exit(1)
		}

		write := false
		path := os.Args[2]
		if os.Args[2] == "-w" {
			if len(os.Args) < 4 {
				fmt.Fprintln(os.Stderr, "usage: got hash-object [-w] <file>")
				os.Exit(1)
			}
			write = true
			path = os.Args[3]
		}

		hash, err := porcelain.HashObject(path, write)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			os.Exit(1)
		}
		fmt.Println(hex.EncodeToString(hash[:]))

	case "cat-file":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "usage: got cat-file <hash>")
			os.Exit(1)
		}

		hashBytes, err := hex.DecodeString(os.Args[2])
		if err != nil || len(hashBytes) != 32 {
			fmt.Fprintln(os.Stderr, "error: invalid hash")
			os.Exit(1)
		}
		var hash [32]byte
		copy(hash[:], hashBytes)

		data, err := porcelain.CatFile(hash)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			os.Exit(1)
		}
		os.Stdout.Write(data)
	case "add":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "usage: got add <file|dir>")
			os.Exit(1)
		}
		if err := porcelain.Add(os.Args[2]); err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		os.Exit(1)

	}
}
