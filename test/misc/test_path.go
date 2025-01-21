package test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func RunAll() {
	TestFunc01()
	TestFunc02()
	TestFunc03()
}

func TestFunc01() {
	// content, err := os.ReadFile(".\\readme.md")
	// content, err := os.ReadFile("./readme.md")
	content, err := os.ReadFile("./template/index.html")
	// content, err := os.ReadFile("C:\\Users\\huaji\\Docs\\repos\\combo_site\\app\\fin_data\\data_types.go")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("File contents: %s", content)
}

func TestFunc02() {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println("current path is:", path)
}

func TestFunc03() {
	p := filepath.FromSlash(".\\readme.md")
	fmt.Println("Path: " + p)

	content, err := os.ReadFile(p)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("File contents: %s", content)

}
