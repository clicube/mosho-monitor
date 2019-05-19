package main

import (
	"fmt"
	"os"

	"mosho-monitor/internal"
)

func main() {
	err := internal.Start()
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(-1)
	}
	os.Exit(0)
}
