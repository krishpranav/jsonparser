# jsonparser
A simple json parser built using golang 

## Installation:
```
go get -u github.com/krishpranav/jsonparser
```

## Installing the cli tool:
- [macOS](https://github.com/krishpranav/jsonparser/releases/download/v1/macOS.zip)
- [windows](https://github.com/krishpranav/jsonparser/releases/download/v1/windows.zip)
- [linux](https://github.com/krishpranav/jsonparser/releases/download/v1/linux.zip)

## Cli tool:
```bash
$ jsonparser -d 1 < test.json
```

## Tutorial:
```golang
package main

import (
	"fmt"
	"os"

	"github.com/krishpranav/jsonparser"
)

func main() {
	f, _ := os.Open("test.json")
	decoder := jsonparser.NewDecoder(f, 1)
	for mv := range decoder.Stream() {
		fmt.Printf("%v\n ", mv.Value)
	}
}
```

## Author:
- [krishpranav](https://github.com/krishpranav)