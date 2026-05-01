package main

import (
	"github.com/KWDB/kwcli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}