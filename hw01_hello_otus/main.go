package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	inputPhrase := "Hello, OTUS!"

	resultPhrase := stringutil.Reverse(inputPhrase)

	fmt.Println(resultPhrase)
}
