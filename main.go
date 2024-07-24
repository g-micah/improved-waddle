package main

import (
	"fmt"

	"rsc.io/quote/v4"
)

func main() {
	test := "some kind of string" //idkasdf asdf
	test2 := "another strhinga"   // idk 2
	fmt.Println("Hello there :) " + test + test2)
	fmt.Println(quote.Go())
}
