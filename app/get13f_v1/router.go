package get13f_v1

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func FundList(w http.ResponseWriter, r *http.Request) {

	website_tmpl, err := template.ParseFiles("./template/header_footer.html")
	if err != nil {
		log.Println(err)
	}

	byte_file, err := os.ReadFile("./template/findata/list_of_fund.html")
	if err != nil {
		log.Println(err)
	}

	home_content_html := struct {
		Home_Content_HTML template.HTML
	}{
		Home_Content_HTML: template.HTML(string(byte_file)),
	}

	website_tmpl.Execute(w, home_content_html)
}

func FundName(w http.ResponseWriter, r *http.Request) {
	website_tmpl, err := template.ParseFiles("./template/header_footer.html")
	if err != nil {
		log.Println(err)
	}

	vars := mux.Vars(r)
	fund_name := strings.ToLower(strings.TrimSpace(vars["fund_name"]))

	fundList := make(map[string]string)
	fundList["himalaya"] = "0001709323"
	fundList["christopher-davis"] = "0001036325"
	fundList["valueact"] = "0001418814"
	fundList["scion"] = "0001649339"
	fundList["third-avenue"] = "0001099281"
	fundList["soros"] = "0001748240"
	fundList["bridge-water"] = "0001350694"
	fundList["abram"] = "0001358706"
	fundList["david-tepper"] = "0001656456"
	fundList["Seth-Klarman"] = "0001061768"
	fundList["point72"] = "0001603466"
	fundList["valueact"] = "0001418814"
	fundList["berkshire"] = "0001067983"
	fundList["brave-warrior"] = "0001553733"

	cik, found := fundList[fund_name]

	// check for cik existence as well
	wd, err := findGoProjectRoot()
	if err != nil {
		log.Println(err)
	}
	cik_json_path := fmt.Sprintf("%s/app/get13f/filings/%s/cik.json", wd, cik)

	if !found || !fileExists(cik_json_path) {
		home_content_html := struct {
			Home_Content_HTML template.HTML
		}{
			Home_Content_HTML: template.HTML("Fund not collected here!"),
		}
		website_tmpl.Execute(w, home_content_html)

		time.Sleep(1 * time.Second)
		return
	}

	cik_data := readCIK(cik)
	most_recent2_filing := filter13F(cik_data)

	// href = fund/accNum, link on report date. desc filing date
	// 	most_recent2_filing[0].AccessionNumber
	// most_recent2_filing[0].ReportDate
	// most_recent2_filing[0].FilingDate

	item_tmpl := template.Must(template.ParseFiles("./template/findata/list_of_items.html"))
	var html_string bytes.Buffer

	type fundTran struct {
		ItemHref string
		ItemLink string
		ItemDesc string
	}

	item1 := fundTran{
		ItemHref: fmt.Sprintf("%s/%s", cik, most_recent2_filing[0].AccessionNumber),
		ItemLink: "Filing",
		ItemDesc: fmt.Sprintf("Reported on: %s. For month end: %s.", most_recent2_filing[0].FilingDate, most_recent2_filing[0].ReportDate),
	}

	item2 := fundTran{
		ItemHref: fmt.Sprintf("%s/%s", cik, most_recent2_filing[1].AccessionNumber),
		ItemLink: "Filing",
		ItemDesc: fmt.Sprintf("Reported on: %s. For month end: %s.", most_recent2_filing[1].FilingDate, most_recent2_filing[1].ReportDate),
	}
	// fmt.Println(item1.ItemHref, item1.ItemLInk, item1.ItemDesc)

	data := []fundTran{item1, item2}

	item_tmpl.Execute(&html_string, data)

	// fmt.Println(html_string.String())

	home_content_html := struct {
		Home_Content_HTML template.HTML
	}{
		Home_Content_HTML: template.HTML(html_string.String()),
	}

	// fmt.Println(html_string.String())
	website_tmpl.Execute(w, home_content_html)

}

func FundHolding(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	fund_name := strings.TrimSpace(vars["fund_name"])

	website_tmpl, err := template.ParseFiles("./template/header_footer.html")
	if err != nil {
		log.Println(err)
	}

	var fund_holding []byte
	if fund_name == "Soros" {
		fund_holding, err = os.ReadFile("./test/tmp/html_table_v2.html")
		if err != nil {
			log.Println(err)
		}

	} else {
		fund_holding = []byte(`Not Found!`)
	}

	// fmt.Println(string(fund_holding))

	home_content_html := struct {
		Home_Content_HTML template.HTML
	}{
		Home_Content_HTML: template.HTML(string(fund_holding)),
	}

	website_tmpl.Execute(w, home_content_html)
}
