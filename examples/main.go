package main

import (
	"fmt"
	"os"

	"github.com/krishpranav/jsonparser"
)

func main() {
	f, _ := os.Open("input.json")
	decoder := jsonparser.NewDecoder(f, 1)
	for mv := range decoder.Stream() {
		fmt.Printf("%v\n ", mv.Value)
	}
}
