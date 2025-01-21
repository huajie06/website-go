package test

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type InformationTable struct {
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

type Asset struct {
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

func printStructFieldsAndValues(s interface{}) {
	v := reflect.ValueOf(s)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		fmt.Printf("%s: %v\n", field.Name, value.Interface())
	}
}

func Parse_txt(txt_file string) {

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

	// fmt.Println(string(xml_result))

	var infotables InformationTable
	var ListOfAsset []Asset

	if err = xml.Unmarshal(xml_result, &infotables); err != nil {
		log.Println(err)
	}

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

	// fmt.Println(xml_struct)

	// var ListOfAsset []Asset
	// for _, v := range xml_struct.InfoTables {
	// 	ListOfAsset = append(ListOfAsset, Asset{
	// 		NameOfIssuer:         v.NameOfIssuer,
	// 		Cusip:                v.Cusip,
	// 		Value:                v.Value,
	// 		SshPrnamt:            v.ShrsOrPrnAmt.SshPrnamt,
	// 		PutCall:              v.PutCall,
	// 		InvestmentDiscretion: v.InvestmentDiscretion,
	// 		OtherManager:         v.OtherManager,
	// 		Sole:                 v.VotingAuthority.Sole,
	// 		Shared:               v.VotingAuthority.Shared,
	// 		None:                 v.VotingAuthority.None,
	// 	})
	// }

	for _, v := range ListOfAsset[:3] {
		printStructFieldsAndValues(v)
		fmt.Println("------")
	}

}

func Loop_txt() {
	rootDir := "../app/get13f/filings" // Replace with the path to your 'filings' directory

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}

		// Check if the file is a .txt file (case-insensitive)
		if strings.ToLower(filepath.Ext(path)) == ".txt" {
			fmt.Println("=================")
			fmt.Println(path)
			Parse_txt(path)
			fmt.Println("=================")
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
	}
}

func Parse_txt_v1() {

	// test_file := "C:/Users/huaji/Docs/repos/website-go/test/tmp/test_txt.txt"
	test_file := "C:/Users/huaji/Docs/repos/website-go/app/get13f/filings/1418814/0001418812-23-000026.txt"

	file, err := os.Open(test_file)
	if err != nil {
		log.Println(err)
	}

	scanner := bufio.NewScanner(file)

	var xml_result = []byte{}

	xml_counter := 0
	infortable_counter := 0

	for scanner.Scan() {
		line := scanner.Text()
		Trimedline := strings.TrimSpace(line)
		Trimed_len := len(Trimedline)

		if Trimed_len >= 5 && Trimedline[:5] == "<XML>" {
			xml_counter++
		}
		if Trimed_len >= 6 && Trimedline[:6] == "</XML>" {
			xml_counter--
		}

		if Trimed_len >= 17 && strings.Contains(Trimedline, "informationTable ") {
			infortable_counter++
		}

		if Trimed_len >= 18 && strings.Contains(Trimedline, "informationTable>") {
			infortable_counter--
		}

		if xml_counter == 1 && infortable_counter == 1 {
			xml_result = append(xml_result, line...)
		}

		// fmt.Printf("value: %s. length: %d. xml: %d. info: %d\n", Trimedline, Trimed_len, xml_counter, infortable_counter)
	}
	fmt.Println(string(xml_result))
}
