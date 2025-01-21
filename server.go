package main

import (
	"combo-site/app/get13f"
	"combo-site/app/index"
	notes "combo-site/app/notes"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/favicon.ico", index.FaviconHandler)

	// this helps on the html link to css but only index page
	// router.PathPrefix("/css/").Handler(http.FileServer(http.Dir("template")))

	// first route

	router.HandleFunc("/", index.IndexPage)

	// route for all taken notes
	router.HandleFunc("/notes", notes.ReturnAllNotes)

	// route for post a note
	router.HandleFunc("/take-note", notes.SubmitNote)

	router.HandleFunc("/note/{note_id}", notes.ReadOneNote)

	router.HandleFunc("/fund", get13f.RenderFundList)
	router.HandleFunc("/fund/{fund_CIK}", get13f.FundQuarterlyList)
	router.HandleFunc("/fund/{fund_CIK}/q", get13f.FundFiling_individual)

	fmt.Println("server starting on: 5000")
	http.ListenAndServe(":5000", router)
}

// func testFunc() {
// 	test.TestEvalTmpl()
// }
