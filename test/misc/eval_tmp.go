package test

import (
	"fmt"
	"log"
	"os"
	"text/template"
)

type ListCards struct {
	LCard []Card
}

type Card struct {
	SectionTitle  string
	SectionText   string
	SectionButton string
}

func TestEvalTmpl() {

	// if useing cards, in the template using {{range .}}
	cards := []Card{{SectionTitle: "title1", SectionText: "long text", SectionButton: "submit 1"}, {SectionTitle: "title2", SectionText: "long text v2", SectionButton: "submit 2"}}
	fmt.Println(cards)

	a := Card{SectionTitle: "title1", SectionText: "long text", SectionButton: "submit 1"}
	b := Card{SectionTitle: "title2", SectionText: "long text v2", SectionButton: "submit 2"}

	NewCards := ListCards{LCard: []Card{a, b}}
	fmt.Println(NewCards)

	tmpl, err := template.ParseFiles("./test/template_test.html")
	if err != nil {
		log.Println(err)
	}
	err = tmpl.Execute(os.Stdout, NewCards)
	if err != nil {
		fmt.Println(err)
	}

}
