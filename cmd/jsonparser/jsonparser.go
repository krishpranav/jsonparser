package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	depthFlag   = flag.Int("d", 0, "")
	kvFlag      = flag.Bool("kv", false, "")
	verboseFlag = flag.Bool("v", false, "")
	helpFlag    = flag.Bool("h", false, "")
)

func exitErr(err error) {
	fmt.Fprintf(os.Stderr, "[\033[31merror\033[0m] %s", err)
	os.Exit(1)
}

var Msg = `jsonparse - a golang json parsing tool

using: jsonpars [options]

options:

`

func help() {
	fmt.Println(Msg)
}
