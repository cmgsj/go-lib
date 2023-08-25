package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	for _, kv := range os.Environ() {
		k, v, ok := strings.Cut(kv, "=")
		if ok {
			fmt.Printf("%s=%q\n", k, v)
		}
	}
}
