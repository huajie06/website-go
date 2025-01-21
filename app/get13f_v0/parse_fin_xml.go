package get13f_v0

import (
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
	"strconv"
	"strings"
	"time"
)

// 1. provide a CIK, which is a company ID
// 2. get the CIK.json, which has all recent filings, and filter out the most recent filing, need accession number
// 3. create a directory for CIK if not exist, and check for the accession-number.xml
// 4. if not exist, download and create csv, call python.  if exist call python
// 5. python wil send html output and then embed it into the page

func removeLeadingZeros(input string) string {
	// Attempt to parse the input string as an integer
	num, err := strconv.Atoi(input)
	if err != nil {
		// If parsing fails, return the original input
		return input
	}

	// Convert the integer back to a string, which will remove leading zeros
	return strconv.Itoa(num)
}

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

	dir := fmt.Sprintf("%s/app/get13f_v0/filings/%s", wd, removeLeadingZeros(cik))

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Println(err)
	}
}

// get the CIK.json
func getCompanyJson(cik string) {
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
	cik_json_path := fmt.Sprintf("%s/app/get13f_v0/filings/%s/cik.json", wd, removeLeadingZeros(cik))

	os.WriteFile(cik_json_path, body, os.ModePerm)

}

// read the CIK.json and parse into struct

func readCIK(cik string) CIKSub {

	wd, err := findGoProjectRoot()
	if err != nil {
		log.Println(err)
	}
	cik_json_path := fmt.Sprintf("%s/app/get13f_v0/filings/%s/cik.json", wd, removeLeadingZeros(cik))

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

func downloadXML(xml_url, output_path string) error {
	// if fileExists(output_path) {
	// 	log.Println("file existed")
	// 	return nil
	// }

	var client http.Client
	r, err := http.NewRequest("GET", xml_url, nil)
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

func xmlURLBuild(cik, accNum string) string {
	accnbr := strings.ReplaceAll(accNum, "-", "")
	cik_strip0 := removeLeadingZeros(cik)
	// xml_url := fmt.Sprintf("https://www.sec.gov/Archives/edgar/data/%s/%s/informationtable.xml", cik_strip0, accnbr)
	xml_url := fmt.Sprintf("https://www.sec.gov/Archives/edgar/data/%s/%s/infotable.xml", cik_strip0, accnbr)
	return xml_url
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

func AssetList_toCSV(list_of_asset []Asset, out_csv string) {
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

func ParseXML(xml_local string) []Asset {
	byte_content, err := os.ReadFile(xml_local)
	if err != nil {
		log.Println(err)
	}

	var infotables InformationTable
	var ListOfAsset []Asset

	xml.Unmarshal(byte_content, &infotables)
	holdings := infotables.InfoTable

	for _, v := range holdings {
		ListOfAsset = append(ListOfAsset, Asset{
			NameOfIssuer:         v.NameOfIssuer,
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

	// for id, v := range ListOfAsset {
	// 	fmt.Printf("%d: %s\n", id, v)
	// }
	return ListOfAsset
}

func Download_AccNum_toCSV(cik, accNum string) {
	cik_strip0 := removeLeadingZeros(cik)
	xml_url := xmlURLBuild(cik, accNum)

	wd, err := findGoProjectRoot()
	if err != nil {
		log.Println(err)
	}
	filing_local_xml := fmt.Sprintf("%s/app/get13f_v0/filings/%s/%s.xml", wd, cik_strip0, accNum)

	// if file exist already, then end
	if fileExists(filing_local_xml) {
		log.Printf("File existed:\n%s", filing_local_xml)
		return
	}

	// if not found, 1.download. 2.parse to csv
	downloadXML(xml_url, filing_local_xml)
	all_asset := ParseXML(filing_local_xml)

	filing_local_csv := fmt.Sprintf("%s/app/get13f_v0/filings/%s/%s.csv", wd, cik_strip0, accNum)
	AssetList_toCSV(all_asset, filing_local_csv)
}

func RunAll() {
	cik := "0001603466"
	setupWD(cik)
	getCompanyJson(cik)

	// 1. read local CIK.json
	filing_13f := readCIK(cik)

	// 2. filter to most recent 2 and remove un-needed fields
	most_recent2 := filter13F(filing_13f)
	// fmt.Println(most_recent2)

	curr := most_recent2[0].AccessionNumber
	prev := most_recent2[1].AccessionNumber
	fmt.Println(curr)
	fmt.Println(prev)

	// for each filing, get accNum
	Download_AccNum_toCSV(cik, curr)
	time.Sleep(1 * time.Second)
	Download_AccNum_toCSV(cik, prev)

	// xml_url := xmlURLBuild(cik, prev)
	// wd, err := findGoProjectRoot()
	// if err != nil {
	// 	log.Println(err)
	// }
	// filing_local_xml := fmt.Sprintf("%s/app/get13f/filings/%s/%s.xml", wd, removeLeadingZeros(cik), prev)

	// // 4. download xml
	// downloadXML(xml_url, filing_local_xml)

	// // 5. read local filing
	// all_asset := ParseXML(filing_local_xml)

	// // 6. save it to csv
	// filing_local_csv := fmt.Sprintf("%s/app/get13f/filings/%s/%s.csv", wd, removeLeadingZeros(cik), prev)
	// AssetList_toCSV(all_asset, filing_local_csv)

}
