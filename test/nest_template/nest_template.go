package test

import (
	"log"
	"os"
	"text/template"
)

type Card struct {
	SectionTitle  string
	SectionText   string
	SectionButton string
}

type Renderdata struct {
	Cards []Card
}

func Nest_template() {

	// the template needs to expand to others come first
	tmpl, err := template.ParseFiles("./nest_template/content.html", "./nest_template/header.html", "./nest_template/footer.html")

	if err != nil {
		log.Println(err)
	}

	card1 := Card{SectionTitle: "title1", SectionText: "long text", SectionButton: "submit 1"}
	card2 := Card{SectionTitle: "title2", SectionText: "long text v2", SectionButton: "submit 2"}

	data := Renderdata{Cards: []Card{card1, card2}}
	err = tmpl.Execute(os.Stdout, data)

	if err != nil {
		log.Println(err)
	}
}
