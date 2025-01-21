package notes

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

const meta_file = "./db/meta.json"

// go to a db or a loc to find all the notes
// need to assign each notes to a unique ID
// or need a metadata about the all the notes

func get_list_of_notes() []NewNote_Meta {
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

	var dt time.Time
	var data2 []NewNote_Meta
	for _, v := range data {
		dt, err = time.Parse("20060102_150405", v.NoteTimestamp)
		if err != nil {
			log.Println(err)
		}
		data2 = append(data2, NewNote_Meta{Note_Metadata: v, DatetimeStr: dt.Format("2006-01-02 15:04")})
	}

	// sort data2 descending order
	sort.Slice(data2, func(i, j int) bool { return data2[i].DatetimeStr > data2[j].DatetimeStr })

	return data2
}

func append_meta(newItem Note_Metadata) {

	b, err := os.ReadFile(meta_file)
	if err != nil {
		log.Println(err)
	}

	data := []Note_Metadata{}
	err = json.Unmarshal(b, &data)
	if err != nil {
		log.Println(err)
	}

	// newItem := Note_Metadata{Fname: "init.vim", NoteTimestamp: "2023-01-01"}
	data = append(data, newItem)

	result, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		log.Println(err)
	}

	os.WriteFile(meta_file, result, os.ModePerm)
}

func ReadOneNote_v0(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	note_id := strings.TrimSpace(vars["note_id"])

	fpath := fmt.Sprintf("./db/stored_notes/%s", note_id)

	content, err := os.ReadFile(fpath)
	if err != nil {
		log.Println(err)
	}

	tmpl, err := template.ParseFiles("./template/individual_note.html", "./template/header.html", "./template/footer.html")
	if err != nil {
		log.Println(err)
	}

	tmpl.Execute(w, struct{ Note string }{string(content)})
}
