package test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func Test_folder() {
	dir := filepath.Dir("./test/get13f/x1/x2")
	fmt.Println(dir)

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Println(err)
	}
}
