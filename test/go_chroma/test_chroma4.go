package test

import (
	"bytes"
	"html/template"
	"log"
	"os"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

func Run4() {
	fpath := "./main.go"
	content, err := os.ReadFile(fpath)
	if err != nil {
		log.Println(err)
	}

	l := lexers.Match(fpath)
	if l == nil {
		l = lexers.Fallback
	}

	// setup style and format
	style := styles.Get("github")
	formatter := html.New(html.WithClasses(true), html.WithLineNumbers(true))

	// tokenize each element
	iterator, err := l.Tokenise(nil, string(content))
	if err != nil {
		log.Println(err)
	}

	// css strings
	var buf_css bytes.Buffer
	err = formatter.WriteCSS(&buf_css, style)
	if err != nil {
		log.Println(err)
	}

	// html strings
	var buf_html bytes.Buffer
	err = formatter.Format(&buf_html, style, iterator)
	if err != nil {
		log.Println(err)
	}

	tmpl := template.Must(template.ParseFiles("../template/individual_note.html"))
	var html_out bytes.Buffer

	tmpl_data := struct {
		Css_str  template.CSS
		Html_str template.HTML
	}{
		template.CSS(buf_css.String()),
		template.HTML(buf_html.String()),
	}

	tmpl.Execute(&html_out, tmpl_data)

	os.WriteFile("./tmp/test_chroma_v5.html", html_out.Bytes(), os.ModePerm)
}
