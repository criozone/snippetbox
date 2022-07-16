package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	binPath, err := filepath.Abs(os.Args[0])
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Exec path: ", binPath)

}
