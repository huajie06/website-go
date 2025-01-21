package get13f_v1

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

/***********************************************************
1. provide a CIK, which is a company ID
2. get the CIK.json, which has all recent filings, and filter out the most recent filing, need accession number
3. create a directory for CIK if not exist, and check for the accession-number.xml
4. if not exist, download and create csv, call python.  if exist call python
5. python wil send html output and then embed it into the page
***********************************************************/

func findGoProjectRoot() (string, error) {
	// Start from the current working directory
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Walk up the directory tree until we find a go.mod file or reach the root directory
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		// Get the parent directory
		parent := filepath.Dir(dir)

		// If we haven't reached the root directory, continue searching
		if parent == dir {
			break
		}

		dir = parent
	}

	// If no go.mod file was found, return an error
	return "", fmt.Errorf("project root not found")
}

// first step is to create the CIK folder (with leading zero removed)
func setupWD(cik string) {

	wd, err := findGoProjectRoot()
	if err != nil {
		log.Println(err)
	}
	// fmt.Println(wd)

	dir := fmt.Sprintf("%s/app/get13f/filings/%s", wd, cik)

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Println(err)
	}
}

// get the CIK.json
func downloadCIKJson(cik string) {
	url_s := fmt.Sprintf("https://data.sec.gov/submissions/CIK%s.json", cik)

	u, err := url.Parse(url_s)
	if err != nil {
		log.Println(err)
	}

	var client http.Client
	r, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		log.Println(err)
	}

	r.Header.Set("User-Agent", "Yorkshire LLC/1.0")

	resp, err := client.Do(r)
	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	wd, err := findGoProjectRoot()
	if err != nil {
		log.Println(err)
	}
	cik_json_path := fmt.Sprintf("%s/app/get13f/filings/%s/cik.json", wd, cik)

	os.WriteFile(cik_json_path, body, os.ModePerm)

}

// read the CIK.json and parse into struct

func readCIK(cik string) CIKSub {

	wd, err := findGoProjectRoot()
	if err != nil {
		log.Println(err)
	}
	cik_json_path := fmt.Sprintf("%s/app/get13f/filings/%s/cik.json", wd, cik)

	cik_byte, err := os.ReadFile(cik_json_path)
	if err != nil {
		log.Println(err)
	}

	var CIKSub_data CIKSub
	if err = json.Unmarshal(cik_byte, &CIKSub_data); err != nil {
		log.Println(err)
	}

	// fmt.Println(SEC13FJson_data)

	return CIKSub_data
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func downloadTXT(txt_url, output_path string) error {

	if fileExists(output_path) {
		log.Printf("File existed:\n%s", output_path)
		return nil
	}

	var client http.Client
	r, err := http.NewRequest("GET", txt_url, nil)
	if err != nil {
		log.Println(err)
	}
	r.Header.Set("User-Agent", "Yorkshire LLC/1.0")

	resp, err := client.Do(r)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	// Create the output file
	outFile, err := os.Create(output_path)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Copy the response body to the output file
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return err
	}

	return nil

}

// take CIK.json struct and filter to 13F, return the first 2 filings - (accesion number)
func filter13F(ciksubmissing CIKSub) []EachFiling {
	recent_filings := ciksubmissing.Filings.Recent
	filing_len := len(recent_filings.Form)

	var filtered13f = []EachFiling{}

	// will return only 13F-HRs
	for i := 0; i < filing_len; i++ {
		if recent_filings.Form[i] == "13F-HR" {
			filtered13f = append(filtered13f, EachFiling{
				AccessionNumber: recent_filings.AccessionNumber[i],
				FileNumber:      recent_filings.FileNumber[i],
				FilingDate:      recent_filings.FilingDate[i],
				ReportDate:      recent_filings.ReportDate[i],
				PrimaryDocument: recent_filings.PrimaryDocument[i]})
		}
	}
	return filtered13f[:2]
}

func SECtxtURLBuild(cik, accNum string) string {
	accNum_nodash := strings.ReplaceAll(accNum, "-", "")
	txt_url := fmt.Sprintf("https://www.sec.gov/Archives/edgar/data/%s/%s/%s.txt", cik, accNum_nodash, accNum)
	return txt_url
}

func LocalTXTFile(cik, accNum string) string {
	wd, err := findGoProjectRoot()
	if err != nil {
		log.Println(err)
	}
	local_txt := fmt.Sprintf("%s/app/get13f/filings/%s/%s.txt", wd, cik, accNum)

	return local_txt
}

func LocalCSVFile(cik, accNum string) string {
	wd, err := findGoProjectRoot()
	if err != nil {
		log.Println(err)
	}
	local_csv := fmt.Sprintf("%s/app/get13f/filings/%s/%s.csv", wd, cik, accNum)
	return local_csv
}

// getStructFieldNames returns the field names of a struct as strings.
func getStructFieldNames(s interface{}) []string {
	var fieldNames []string
	t := reflect.TypeOf(s)
	for i := 0; i < t.NumField(); i++ {
		fieldNames = append(fieldNames, t.Field(i).Name)
	}
	return fieldNames
}

func AssetListStruct_toCSV(list_of_asset []Asset, out_csv string) {
	csvFile, err := os.Create(out_csv)
	if err != nil {
		log.Println(err)
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()

	fieldNames := getStructFieldNames(Asset{})

	// Write the CSV header using field names
	if err := csvWriter.Write(fieldNames); err != nil {
		log.Println(err)
	}

	for _, row := range list_of_asset {
		string_row := []string{
			row.NameOfIssuer,
			row.TitleOfClass,
			row.Cusip,
			row.Value,
			row.SshPrnamt,
			row.PutCall,
			row.InvestmentDiscretion,
			row.OtherManager,
			row.Sole,
			row.Shared,
			row.None}
		err := csvWriter.Write(string_row)
		if err != nil {
			log.Println(err)
		}
	}
}

func ParseTXT_toXML(txt_file string) []byte {
	file, err := os.Open(txt_file)
	if err != nil {
		log.Println(err)
	}

	scanner := bufio.NewScanner(file)

	var xml_result = []byte{}
	informationTable := false

	for scanner.Scan() {
		line := scanner.Text()
		Trimedline := strings.TrimSpace(line)

		if !informationTable && strings.Contains(Trimedline, "informationTable") {
			informationTable = true
			xml_result = append(xml_result, Trimedline...)
			continue
		}

		if informationTable && strings.Contains(Trimedline, "informationTable") {
			informationTable = false
			xml_result = append(xml_result, Trimedline...)
			continue
		}

		if informationTable {
			xml_result = append(xml_result, Trimedline...)
		}
	}
	return xml_result
}

func ParseXML(xml_in_byte []byte) []Asset {
	var infotables InformationTable
	var ListOfAsset []Asset

	xml.Unmarshal(xml_in_byte, &infotables)
	holdings := infotables.InfoTable

	for _, v := range holdings {
		ListOfAsset = append(ListOfAsset, Asset{
			NameOfIssuer:         v.NameOfIssuer,
			TitleOfClass:         v.TitleOfClass,
			Cusip:                v.Cusip,
			Value:                v.Value,
			SshPrnamt:            v.ShrsOrPrnAmt.SshPrnamt,
			PutCall:              v.PutCall,
			InvestmentDiscretion: v.InvestmentDiscretion,
			OtherManager:         v.OtherManager,
			Sole:                 v.VotingAuthority.Sole,
			Shared:               v.VotingAuthority.Shared,
			None:                 v.VotingAuthority.None,
		})
	}

	return ListOfAsset
}

func TXT_toCSV(txt_file_path, out_csv_path string) {
	xml_results := ParseTXT_toXML(txt_file_path)
	all_asset := ParseXML(xml_results)
	AssetListStruct_toCSV(all_asset, out_csv_path)
}

func Download_AccNum_toCSV(cik, accNum string) {
	txt_url := SECtxtURLBuild(cik, accNum)
	filing_local_txt := LocalTXTFile(cik, accNum)
	filing_local_csv := LocalCSVFile(cik, accNum)

	// if file exist already, then end
	if fileExists(filing_local_txt) {
		log.Printf("File existed:\n%s", filing_local_txt)
		return
	}

	if fileExists(filing_local_csv) {
		log.Printf("File existed:\n%s", filing_local_csv)
		return
	}

	// if not found, 1.download. 2.parse to csv
	downloadTXT(txt_url, filing_local_txt)
	xml_results := ParseTXT_toXML(filing_local_txt)
	all_asset := ParseXML(xml_results)

	AssetListStruct_toCSV(all_asset, filing_local_csv)
}

func Test_v0() {
	cik := "0001067983"
	accNum := "0000950123-23-008074"
	filing_local_txt := LocalTXTFile(cik, accNum)

	xml_byte := ParseTXT_toXML(filing_local_txt)

	all_asset := ParseXML(xml_byte)

	fmt.Println(all_asset[0])
	printStructFieldsAndValues(all_asset[0])

	// var infotables InformationTable

	// xml.Unmarshal(xml_byte, &infotables)
	// holdings := infotables.InfoTable

	// for _, v := range holdings[:3] {

	// }
}

func printStructFieldsAndValues(s interface{}) {
	v := reflect.ValueOf(s)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		fmt.Printf("%s: %v\n", field.Name, value.Interface())
	}
}

func RunAll() {
	cik := "0001067983"
	setupWD(cik)
	downloadCIKJson(cik)

	// 1. read local CIK.json
	filing_13f := readCIK(cik)

	// 2. filter to most recent 2 and remove un-needed fields
	most_recent2 := filter13F(filing_13f)

	curr := most_recent2[0].AccessionNumber
	prev := most_recent2[1].AccessionNumber
	fmt.Println(curr)
	fmt.Println(prev)

	which_to_fetch := prev

	txt_url := SECtxtURLBuild(cik, which_to_fetch)
	filing_local_txt := LocalTXTFile(cik, which_to_fetch)
	filing_local_csv := LocalCSVFile(cik, which_to_fetch)

	// 4. download xml
	if err := downloadTXT(txt_url, filing_local_txt); err != nil {
		log.Println(err)
	}

	TXT_toCSV(filing_local_txt, filing_local_csv)
}
