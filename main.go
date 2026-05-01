package main

import (
	"github.com/shawn0915/kwcli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}