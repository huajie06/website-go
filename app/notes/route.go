package notes

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/gorilla/mux"
)

func ReadOneNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	note_id := strings.TrimSpace(vars["note_id"])

	fpath := fmt.Sprintf("./db/stored_notes/%s", note_id)

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

	tmpl := template.Must(template.ParseFiles("./template/notes/individual_note.html"))
	// var html_out bytes.Buffer

	tmpl_data := struct {
		Css_str  template.CSS
		Html_str template.HTML
	}{
		template.CSS(buf_css.String()),
		template.HTML(buf_html.String()),
	}

	tmpl.Execute(w, tmpl_data)

	// os.WriteFile("./test/tmp/test_chroma_v3.html", html_out.Bytes(), os.ModePerm)
}

func SubmitNote(w http.ResponseWriter, r *http.Request) {
	website_tmpl := template.Must(template.ParseFiles("./template/header_footer.html"))
	submit_note_tmpl := template.Must(template.ParseFiles("./template/notes/submit_notes.html"))
	var html_string bytes.Buffer

	// the if/else determine what to send to the `take-note` page
	if r.Method != http.MethodPost {
		// if not POST => html_string => FORM
		submit_note_tmpl.Execute(&html_string, nil)

	} else if r.FormValue("Fname") == "" || r.FormValue("Content") == "" {
		// if missing => html_string => FORM
		submit_note_tmpl.Execute(&html_string, nil)

	} else {
		submit_note_tmpl.Execute(&html_string, struct{ Success bool }{true})
		details := Note{
			Fname:         r.FormValue("Fname"),
			Content:       r.FormValue("Content"),
			NoteTimestamp: time.Now(),
		}

		// save the file to server
		SaveFilename := fmt.Sprintf("./db/stored_notes/%s_%s", details.NoteTimestamp.Format("20060102_150405"), strings.ReplaceAll(details.Fname, " ", ""))
		os.WriteFile(SaveFilename, []byte(details.Content), os.ModePerm)

		// modify meta data to include
		newItem := Note_Metadata{Fname: strings.ReplaceAll(details.Fname, " ", ""), NoteTimestamp: details.NoteTimestamp.Format("20060102_150405")}
		append_meta(newItem)
	}

	// send the note html to website
	home_content_html := struct {
		Home_Content_HTML template.HTML
	}{
		Home_Content_HTML: template.HTML(html_string.String()),
	}

	website_tmpl.Execute(w, home_content_html)

}

func ReturnAllNotes(w http.ResponseWriter, request *http.Request) {
	// note template
	note_lst_tmpl, err := template.ParseFiles("./template/notes/list_of_notes.html")
	if err != nil {
		log.Println(err)
	}
	// get a list of notes
	data := get_list_of_notes()
	var html_string bytes.Buffer
	// create note html string
	note_lst_tmpl.Execute(&html_string, data)

	// website template
	website_tmpl, err := template.ParseFiles("./template/header_footer.html")
	if err != nil {
		log.Println(err)
	}
	// create data struct
	home_content_html := struct {
		Home_Content_HTML template.HTML
	}{
		Home_Content_HTML: template.HTML(html_string.String()),
	}

	website_tmpl.Execute(w, home_content_html)

}
