package get13f

import (
	"fmt"
)

func Setup() {
	cik := "0001358706"
	fmt.Println(cik)
	// Delete_bucket("../db/filings/fin.db", fmt.Sprintf("%s_html", cik))

	// ReadFromBucket("../db/filings/fin.db", cik)
	// ReadFromBucket("../db/filings/fin.db", fmt.Sprintf("%s_html", cik))
	// fmt.Println("====================")

	// fmt.Println(string(byte_result))
	// GetTableKeys_fromDB(fmt.Sprintf("%s_html", cik))

	// ReadFromBucket("../db/filings/fin.db", "0001709323")
	// Fund_most_recent_table("0001061768")

	// Wrapper_v1()
	// Download_TXT_Filing_From_DB("0001418814")

	// ListEachBuckets("..//db/filings/fin.db")

	fundlist := FundList()

	for cik, nm := range fundlist {
		fmt.Printf("========runing %s==========", nm)
		ResetDB(cik)
		UpdateOneCIK(cik)
		fmt.Println("=========================")
	}
}
