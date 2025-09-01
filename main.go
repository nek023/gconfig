package main

import (
	"fmt"
	"os"

	"github.com/nek023/gconfig/cmd"
)

// main is the entry point for the gconfig CLI application
func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
