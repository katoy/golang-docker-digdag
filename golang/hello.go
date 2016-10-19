package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Printf("Hello world from go-lang.\n")
	fmt.Printf("こんにちは. こちらは GO 言語です。\n")
	fmt.Println(runtime.GOARCH, runtime.GOOS)
}
