package test

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

func Run_new() {
	// code := `package main

	// func main() {
	// 	fmt.Println("Hello, world!")
	// }`

	code0, err := os.ReadFile("./go_chroma/test_chroma.go")
	if err != nil {
		log.Println(err)
	}
	code := string(code0)

	code2 := `what today`

	l := lexers.Analyse(code)
	// fmt.Println(l)

	fmt.Println(l.Config().Name)

	l2 := lexers.Analyse(code2) // return nil
	fmt.Println(l2)
	if l2 == nil {
		fmt.Println("no format found")
	}

	style := styles.Get("github")
	formatter := html.New(html.WithClasses(true), html.WithLineNumbers(true))

	var buf_css bytes.Buffer
	err = formatter.WriteCSS(&buf_css, style)
	if err != nil {
		log.Println(err)
	}
	// fmt.Println(buf.String())

	iterator, err := l.Tokenise(nil, code)
	if err != nil {
		log.Println(err)
	}
	var buf_html bytes.Buffer
	_ = formatter.Format(&buf_html, style, iterator)

	// fmt.Println(buf_html.String())

	//BELOW IS NOT WORKING!!!

	// tmpl := template.Must(template.ParseFiles("./go_chroma/one_note.html", "../template/header.html", "../template/footer.html"))
	tmpl := template.Must(template.ParseFiles("./go_chroma/one_note.html"))
	var html_out bytes.Buffer

	tmpl_data := struct {
		Css_str  string
		Html_str string
	}{
		buf_css.String(),
		buf_html.String(),
	}

	tmpl.Execute(&html_out, tmpl_data)

	// fmt.Println(html_out.String())
	os.WriteFile("./tmp/test_chroma_v2.html", html_out.Bytes(), os.ModePerm)

	// fmt.Println(buf_css.String())
	// fmt.Println("-----")
	// fmt.Println(buf_html.String())
}
