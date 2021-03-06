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

func exitErr(err error) {
	fmt.Fprintf(os.Stderr, "[\033[31merror\033[0m] %s", err)
	os.Exit(1)
}

func printVal(mv *jsonparser.MetaValue) {
	b, err := json.Marshal(mv.Value)
	if err != nil {
		exitErr(err)
	}

	s := string(b)
	var label string

	switch mv.Value.(type) {
	case []interface{}:
		label = "array  "
	case float64:
		label = "float  "
	case jsonparser.KV:
		label = "kv     "
	case string:
		label = "string "
	case map[string]interface{}:
		label = "object "
	}

	if *verboseFlag {
		end := mv.Offset + mv.Length
		fmt.Printf("%d\t%03d\t%03d\t%s| %s\n", mv.Depth, mv.Offset, end, label, s)
		return
	}
	fmt.Println(s)
}

func main() {
	flag.Parse()
	if *helpFlag {
		help()
		os.Exit(0)
	}

	if *verboseFlag {
		fmt.Println("depth\tstart\tend\ttype   | value")
	}

	decoder := jsonparser.NewDecoder(os.Stdin, *depthFlag)
	if *kvFlag {
		decoder = decoder.EmitKV()
	}
	for mv := range decoder.Stream() {
		printVal(mv)
	}
	if err := decoder.Err(); err != nil {
		exitErr(err)
	}
}

var helpMsg = `jsonparser - A json parsing tool built using golang

usage: jsonparser [options]

options:

  -d <n> emit values at depth n.
  -kv    output inner key value pairs as newly formed objects
  -v     output depth and offset details for each value
  -h     display this help dialog


example:
	jsonparser -d 1 < test.json
`

func help() {
	fmt.Println(helpMsg)
}
