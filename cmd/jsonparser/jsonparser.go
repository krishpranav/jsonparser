package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/krishpranav/jsonparser"
)

var (
	depthFlag   = flag.Int("d", 0, "")
	kvFlag      = flag.Bool("kv", false, "")
	verboseFlag = flag.Bool("v", false, "")
	helpFlag    = flag.Bool("h", false, "")
)

func printVal(mv *jsonparser.MetaValue) {
	b, err := json.Marshal(mv.Value)
	if err != nil {
		exitErr(err)
	}

	s := string(b)
	var label string

	switch mv.Value.(type) {
	case []interface{}:
		label = "array "
	case float64:
		label = "float "
	case jsonparser.KV:
		label = "kv "
	case string:
		label = "string "
	case map[string]interface{}:
		label = "object"
	}

	if *verboseFlag {
		end := mv.Offset + mv.Length
		fmt.Printf("%d\t%03d\t%03d\t%s| %s\n", mv.Depth, mv.Offset, end, label, s)
		return
	}
	fmt.Println(s)
}

func main() {

}

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
