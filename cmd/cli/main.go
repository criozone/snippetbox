package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	list, err := os.ReadDir(".")
	if err != nil {
		log.Fatalln(err)
	}

	for _, dir := range list {
		fmt.Println(dir.Name())
	}
}
