package test

import (
	"fmt"
	"log"
	"os"
)

func ReadFile_v1() {

	content, err := os.ReadFile("../db/notes_20220101.txt")
	if err != nil {
		log.Println(err)
	}

	fmt.Println(string(content))
}
