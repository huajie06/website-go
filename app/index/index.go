package index

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

// reference here: https://www.practical-go-lessons.com/chap-32-templates

func IndexPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./template/header_footer.html"))

	content_bytes, err := os.ReadFile("./template/home/index_content.html")
	if err != nil {
		log.Println(err)
	}

	index_data := Index_Content{
		Home_Content_HTML: template.HTML(string(content_bytes)),
	}
	// log.Println(index_data.Home_Content_HTML)

	tmpl.Execute(w, index_data)
}

func FaviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./image/favicon.ico")
}
