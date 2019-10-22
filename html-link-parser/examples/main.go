package main

import (
	"fmt"
	link "gophercises/html-link-parser"
	"os"
)

func main() {
	filename := "ex2.html"
	r, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Failed to open file %s\n", filename)
		os.Exit(1)
	}
	defer r.Close()

	links, err := link.Parse(r)

	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", links)
}
