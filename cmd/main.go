package main

import (
	"os"
)

func execute() {
	rootCmd := newRootCmd()

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func main() {
	execute()
}
