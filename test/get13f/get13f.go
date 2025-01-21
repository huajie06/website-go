package test

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
	"reflect"
	"strconv"
	"strings"
	"time"
)

type SEC13FJson struct {
	Cik            string `json:"cik,omitempty"`
	EntityType     string `json:"entityType,omitempty"`
	Sic            string `json:"sic,omitempty"`
	SicDescription string `json:"sicDescription,omitempty"`
	// InsiderTransactionForOwnerExists  int       `json:"insiderTransactionForOwnerExists,omitempty"`
	// InsiderTransactionForIssuerExists int       `json:"insiderTransactionForIssuerExists,omitempty"`
	Name string `json:"name,omitempty"`
	// Tickers                           []any     `json:"tickers,omitempty"`
	// Exchanges                         []any     `json:"exchanges,omitempty"`
	// Ein                               string    `json:"ein,omitempty"`
	// Description                       string    `json:"description,omitempty"`
	// Website                           string    `json:"website,omitempty"`
	// InvestorWebsite                   string    `json:"investorWebsite,omitempty"`
	// Category                          string    `json:"category,omitempty"`
	// FiscalYearEnd                     string    `json:"fiscalYearEnd,omitempty"`
	// StateOfIncorporation              string    `json:"stateOfIncorporation,omitempty"`
	// StateOfIncorporationDescription   string    `json:"stateOfIncorporationDescription,omitempty"`
	// Addresses                         Addresses `json:"addresses,omitempty"`
	// Phone                             string    `json:"phone,omitempty"`
	// Flags                             string    `json:"flags,omitempty"`
	// FormerNames                       []any     `json:"formerNames,omitempty"`
	Filings Filings `json:"filings,omitempty"`
}

//	type Mailing struct {
//		Street1                   string `json:"street1,omitempty"`
//		Street2                   string `json:"street2,omitempty"`
//		City                      string `json:"city,omitempty"`
//		StateOrCountry            string `json:"stateOrCountry,omitempty"`
//		ZipCode                   string `json:"zipCode,omitempty"`
//		StateOrCountryDescription string `json:"stateOrCountryDescription,omitempty"`
//	}
//
//	type Business struct {
//		Street1                   string `json:"street1,omitempty"`
//		Street2                   string `json:"street2,omitempty"`
//		City                      string `json:"city,omitempty"`
//		StateOrCountry            string `json:"stateOrCountry,omitempty"`
//		ZipCode                   string `json:"zipCode,omitempty"`
//		StateOrCountryDescription string `json:"stateOrCountryDescription,omitempty"`
//	}
//
//	type Addresses struct {
//		Mailing  Mailing  `json:"mailing,omitempty"`
//		Business Business `json:"business,omitempty"`
//	}

type Recent struct {
	AccessionNumber       []string    `json:"accessionNumber,omitempty"`
	FilingDate            []string    `json:"filingDate,omitempty"`
	ReportDate            []string    `json:"reportDate,omitempty"`
	AcceptanceDateTime    []time.Time `json:"acceptanceDateTime,omitempty"`
	Act                   []string    `json:"act,omitempty"`
	Form                  []string    `json:"form,omitempty"`
	FileNumber            []string    `json:"fileNumber,omitempty"`
	FilmNumber            []string    `json:"filmNumber,omitempty"`
	Items                 []string    `json:"items,omitempty"`
	Size                  []int       `json:"size,omitempty"`
	IsXBRL                []int       `json:"isXBRL,omitempty"`
	IsInlineXBRL          []int       `json:"isInlineXBRL,omitempty"`
	PrimaryDocument       []string    `json:"primaryDocument,omitempty"`
	PrimaryDocDescription []string    `json:"primaryDocDescription,omitempty"`
}
type Filings struct {
	Recent Recent `json:"recent,omitempty"`
	Files  []any  `json:"files,omitempty"`
}

type EachFiling struct {
	AccessionNumber string
	FileNumber      string
	FilingDate      string
	ReportDate      string
	PrimaryDocument string
	// Form            string
}

type InformationTable1 struct {
	XMLName   xml.Name `xml:"informationTable"`
	Text      string   `xml:",chardata"`
	Xsd       string   `xml:"xsd,attr"`
	Xsi       string   `xml:"xsi,attr"`
	Xmlns     string   `xml:"xmlns,attr"`
	InfoTable []struct {
		Text         string `xml:",chardata"`
		NameOfIssuer string `xml:"nameOfIssuer"`
		TitleOfClass string `xml:"titleOfClass"`
		Cusip        string `xml:"cusip"`
		Value        string `xml:"value"`
		ShrsOrPrnAmt struct {
			Text          string `xml:",chardata"`
			SshPrnamt     string `xml:"sshPrnamt"`
			SshPrnamtType string `xml:"sshPrnamtType"`
		} `xml:"shrsOrPrnAmt"`
		InvestmentDiscretion string `xml:"investmentDiscretion"`
		OtherManager         string `xml:"otherManager"`
		VotingAuthority      struct {
			Text   string `xml:",chardata"`
			Sole   string `xml:"Sole"`
			Shared string `xml:"Shared"`
			None   string `xml:"None"`
		} `xml:"votingAuthority"`
		PutCall string `xml:"putCall"`
	} `xml:"infoTable"`
}

type OneAsset struct {
	NameOfIssuer         string
	Cusip                string
	Value                string
	SshPrnamt            string
	PutCall              string
	InvestmentDiscretion string
	OtherManager         string
	Sole                 string
	Shared               string
	None                 string
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

func ParsedXML_to_CSV(list_of_asset []OneAsset) {
	csvFile, err := os.Create("./get13f/all_asset.csv")
	if err != nil {
		log.Println(err)
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()

	fieldNames := getStructFieldNames(OneAsset{})

	// Write the CSV header using field names
	if err := csvWriter.Write(fieldNames); err != nil {
		panic(err)
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

func ParseXML() {
	byte_content, err := os.ReadFile("./get13f/informationtable.xml")
	if err != nil {
		log.Println(err)
	}

	var infotables InformationTable
	var ListOfAsset []OneAsset

	xml.Unmarshal(byte_content, &infotables)
	holdings := infotables.InfoTable

	for _, v := range holdings {
		ListOfAsset = append(ListOfAsset, OneAsset{
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

	for id, v := range ListOfAsset {
		fmt.Printf("%d: %s\n", id, v)
	}
	ParsedXML_to_CSV(ListOfAsset)
}

func Parse13FJson() {
	fname := "./get13f/cik.json"

	byte_content, err := os.ReadFile(fname)

	if err != nil {
		log.Println(err)
	}

	// fmt.Println(string(byte_content))

	var secresult SEC13FJson

	err = json.Unmarshal(byte_content, &secresult)
	if err != nil {
		log.Println(err)
	}
	//

	recent_filings := secresult.Filings.Recent
	// fmt.Println(recent_filings)
	// var filtered13f Recent
	filing_len := len(recent_filings.Form)

	var valid13F = []EachFiling{}

	for i := 0; i < filing_len; i++ {
		// fmt.Println(recent_filings.AccessionNumber[i])
		if recent_filings.Form[i] == "13F-HR" {
			valid13F = append(valid13F, EachFiling{
				AccessionNumber: recent_filings.AccessionNumber[i],
				FileNumber:      recent_filings.FileNumber[i],
				FilingDate:      recent_filings.FilingDate[i],
				ReportDate:      recent_filings.ReportDate[i],
				PrimaryDocument: recent_filings.PrimaryDocument[i]})
		}
	}
	fmt.Println(valid13F[0])

	// filing_url := "https://www.sec.gov/Archives/edgar/data/{cik}/{accessionNumber_noDash}/{accessionNumber}.txt"

	// fmt.Println(filing_url)

	cik := "0001649339"
	accnbr := strings.ReplaceAll(valid13F[0].AccessionNumber, "-", "")
	xml_url := fmt.Sprintf("https://www.sec.gov/Archives/edgar/data/%s/%s/informationtable.xml", removeLeadingZeros(cik), accnbr)
	fmt.Println(xml_url)

	// fmt.Println(len(recent_filings.AccessionNumber))
	// fmt.Println(len(recent_filings.FileNumber))
	// fmt.Println(len(recent_filings.FilingDate))
	// fmt.Println(len(recent_filings.ReportDate))
	// fmt.Println(len(recent_filings.Form))
	outpath := fmt.Sprintf("./get13f/%s-%s.xml", removeLeadingZeros(cik), accnbr)
	fmt.Println(outpath)

	err = downloadXML(xml_url, outpath)
	if err != nil {
		log.Println(err)
	}
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func downloadXML(xml_url, output_path string) error {
	if fileExists(output_path) {
		log.Println("file existed")
		return nil
	}
	response, err := http.Get(xml_url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Create the output file
	outFile, err := os.Create(output_path)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Copy the response body to the output file
	_, err = io.Copy(outFile, response.Body)
	if err != nil {
		return err
	}

	return nil

}

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

func Get13F() {
	cik := "0001649339"
	url_s := fmt.Sprintf("https://data.sec.gov/submissions/CIK%s.json", cik)

	// fmt.Println(url_s)

	u, err := url.Parse(url_s)
	if err != nil {
		log.Println(err)
	}

	// fmt.Println(u.String())

	var client http.Client
	r, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		log.Println(err)
	}

	r.Header.Set("User-Agent", "Yorkshire LLC/1.0")
	// r.Header.Set("Content-Type", "application/json")

	// fmt.Println(r.Header)

	resp, err := client.Do(r)
	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(string(body))

	var SEC13FJson_data SEC13FJson
	err = json.Unmarshal(body, &SEC13FJson_data)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(SEC13FJson_data)
}
