package get13f

import (
	"encoding/xml"
	"time"
)

type CIKSub struct {
	Cik            string  `json:"cik,omitempty"`
	EntityType     string  `json:"entityType,omitempty"`
	Sic            string  `json:"sic,omitempty"`
	SicDescription string  `json:"sicDescription,omitempty"`
	Name           string  `json:"name,omitempty"`
	Filings        Filings `json:"filings,omitempty"`
}

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
	TitleOfClass         string
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
