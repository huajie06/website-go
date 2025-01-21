package test

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

type Note_Metadata struct {
	Fname         string `json:"Fname"`
	NoteTimestamp string `json:"NoteTimestamp"`
}

const meta_file = "./export_meta/big-json.json"

func WriteJson() {
	f1 := Note_Metadata{Fname: "emacs.el", NoteTimestamp: "2023-03-01"}
	f2 := Note_Metadata{Fname: "emacs.el", NoteTimestamp: "2023-01-01"}

	m := []Note_Metadata{f1, f2}

	b, err := json.MarshalIndent(m, "", "    ")

	if err != nil {
		log.Println(err)
	}

	os.WriteFile(meta_file, b, os.ModePerm)
}

func ReadJson() {
	b, err := os.ReadFile(meta_file)
	if err != nil {
		log.Println(err)
	}

	data := []Note_Metadata{}
	err = json.Unmarshal(b, &data)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(data)
}

func AppendJson() {
	b, err := os.ReadFile(meta_file)
	if err != nil {
		log.Println(err)
	}

	data := []Note_Metadata{}
	err = json.Unmarshal(b, &data)
	if err != nil {
		log.Println(err)
	}

	newItem := Note_Metadata{Fname: "init.vim", NoteTimestamp: "2023-01-01"}
	data = append(data, newItem)

	result, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		log.Println(err)
	}

	os.WriteFile(meta_file, result, os.ModePerm)
}

func AppendJson2() {
	// this is not working
	file, err := os.OpenFile(meta_file, os.O_RDWR, 0644)

	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		log.Println(err)
	}

	newItem := Note_Metadata{Fname: "xxx.vim", NoteTimestamp: "2019-01-01"}
	var current []Note_Metadata
	err = json.Unmarshal(b, &current)
	if err != nil {
		log.Println(err)
	}
	result := append(current, newItem)

	encoder := json.NewEncoder(file)
	err = encoder.Encode(result)
	if err != nil {
		log.Println(err)
	}
}

func Get_list_of_notes() []Note_Metadata {
	// this function show return a list of ID and IDs can be linked to the notes
	// example like `OWp2aMO6-fF`

	b, err := os.ReadFile(meta_file)
	if err != nil {
		log.Println(err)
	}

	var data []Note_Metadata
	err = json.Unmarshal(b, &data)
	if err != nil {
		log.Println(err)
	}

	// for _, v := range data {
	// 	fmt.Println(v)
	// }
	return data
}
