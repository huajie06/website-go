package get13f

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	bolt "go.etcd.io/bbolt"
)

func RenderFundList(w http.ResponseWriter, r *http.Request) {

	website_tmpl, err := template.ParseFiles("./template/header_footer.html")
	if err != nil {
		log.Println(err)
	}

	byte_file, err := os.ReadFile("./template/findata/list_of_fund.html")
	if err != nil {
		log.Println(err)
	}

	main_content_html := struct {
		Home_Content_HTML template.HTML
	}{
		Home_Content_HTML: template.HTML(string(byte_file)),
	}

	website_tmpl.Execute(w, main_content_html)
}

func FundQuarterlyList(w http.ResponseWriter, r *http.Request) {
	type filingList struct {
		ItemHref string
		ItemLink string
		ItemDesc string
	}

	data := []filingList{}
	vars := mux.Vars(r)
	fund_CIK := strings.ToLower(strings.TrimSpace(vars["fund_CIK"]))

	List_of_filings := GetListOfExistedACC(fund_CIK)

	for _, v := range List_of_filings {
		report_filing_dt := strings.Split(v[0], "|")
		item_01 := filingList{
			ItemHref: fmt.Sprintf("/fund/%s/q?accnum=%s&filingdt=%s", fund_CIK, v[1], report_filing_dt[0]),
			ItemLink: "Filing",
			ItemDesc: fmt.Sprintf("Reported on: %s. For month end: %s.", report_filing_dt[1], report_filing_dt[0]),
		}
		data = append(data, item_01)
	}

	fundname := FundList()[fund_CIK]
	fundname_inHTML := fmt.Sprintf("<p>%s</p>", fundname)

	item_tmpl := template.Must(template.ParseFiles("./template/findata/list_of_items.html"))
	var html_string bytes.Buffer
	item_tmpl.Execute(&html_string, data)
	home_content_html := struct {
		Home_Content_HTML template.HTML
	}{
		Home_Content_HTML: template.HTML(fundname_inHTML + html_string.String()),
	}

	website_tmpl, err := template.ParseFiles("./template/header_footer.html")
	if err != nil {
		log.Println(err)
	}
	website_tmpl.Execute(w, home_content_html)

}

func FundFiling_individual(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	fund_CIK := strings.ToLower(strings.TrimSpace(vars["fund_CIK"]))
	// fmt.Println("cik", fund_CIK)

	query := r.URL.Query()
	accnum := query.Get("accnum")
	filing_dt := query.Get("filingdt")

	if accnum == "" {
		http.Error(w, "Missing 'q' parameter", http.StatusBadRequest)
		return
	}

	hmlt_table_string := GetTable_fromDB(fmt.Sprintf("%s_html", fund_CIK), accnum)

	funcName := fmt.Sprintf("<h4>%s</h4><h5>Portfolio Date: %s</h5>", FundList()[fund_CIK], filing_dt)

	website_tmpl, err := template.ParseFiles("./template/header_footer.html")
	if err != nil {
		log.Println(err)
	}

	main_content_html := struct {
		Home_Content_HTML template.HTML
	}{
		Home_Content_HTML: template.HTML(funcName + hmlt_table_string),
	}

	website_tmpl.Execute(w, main_content_html)
}

func Fund_table(w http.ResponseWriter, r *http.Request) {
	website_tmpl, err := template.ParseFiles("./template/header_footer.html")
	if err != nil {
		log.Println(err)
	}

	vars := mux.Vars(r)
	fund_CIK := strings.ToLower(strings.TrimSpace(vars["fund_CIK"]))

	wd, err := findGoProjectRoot()
	if err != nil {
		log.Println(err)
	}

	// initiate a DB for a CIK
	db_dir := fmt.Sprintf("%s/db/filings/fin.db", wd)

	db, err := bolt.Open(db_dir, fs.FileMode(os.O_WRONLY), nil)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	var html_byte []byte

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(fund_CIK))
		c := b.Cursor()
		_, last := c.Last()

		html_byte = b.Get(last)

		return nil
	})

	funcName := fmt.Sprintf("<h4> %s </h4>", FundList()[fund_CIK])

	main_content_html := struct {
		Home_Content_HTML template.HTML
	}{
		Home_Content_HTML: template.HTML(funcName + string(html_byte)),
	}

	website_tmpl.Execute(w, main_content_html)
}
