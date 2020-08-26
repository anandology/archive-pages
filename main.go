package main

import "fmt"

func main() {
	fmt.Printf("Hello, world!\n")
	Serve()
	item, err := GetItem("aaronsw-archive")

	if err != nil {
		panic("Failed to fetch metadata")
	}

	fmt.Println(item)
}
