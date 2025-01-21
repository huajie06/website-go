package main

import (
	"fmt"
	"runtime"
	"strings"

	test2 "combo-site/app/get13f"
)

// test "combo-site/test/import_notes"

// test2 "combo-site/test/nest_template"

func main() {

	// test2.ReadFromBucket("../db/filings/fin.db", "0001709323")
	// test2.Fund_most_recent_table("0001061768")

	// test2.Wrapper_v1()
	// test2.Download_TXT_Filing_From_DB("0001418814")

	// test2.ListEachBuckets("..//db/filings/fin.db")

	cik := "0001358706"
	fmt.Println(cik)
	// test2.Delete_bucket("../db/filings/fin.db", fmt.Sprintf("%s_html", cik))

	// test2.ReadFromBucket("../db/filings/fin.db", cik)
	// test2.ReadFromBucket("../db/filings/fin.db", fmt.Sprintf("%s_html", cik))
	// fmt.Println("====================")

	// fmt.Println(string(byte_result))
	// test2.GetTableKeys_fromDB(fmt.Sprintf("%s_html", cik))

	fundlist := test2.FundList()

	fmt.Println(fundlist)

	for cik, nm := range fundlist {
		fmt.Printf("========runing %s==========", nm)
		test2.ResetDB(cik)
		test2.UpdateOneCIK(cik)
		fmt.Println("=========================")
	}
	// test2.ReadFromBucket("../db/filings/fin.db", fmt.Sprintf("%s_html", cik))
	// checkOS()
}

func checkOS() string {
	osName := runtime.GOOS

	// Get the architecture from the GOARCH environment variable
	arch := runtime.GOARCH

	fmt.Printf("Operating System: %s\n", osName)
	fmt.Printf("Architecture: %s\n", arch)

	// Check if it's running on a Raspberry Pi
	if isRaspberryPi(osName, arch) {
		fmt.Println("Running on a Raspberry Pi")
		return "Pi"
	} else if osName == "windows" {
		fmt.Println("Running on Windows")
		return "Windows"
	} else if osName == "linux" {
		fmt.Println("Running on Linux (Ubuntu or other)")
		return "Linux"
	} else {
		fmt.Println("Running on an unknown operating system")
		return "Unknown"
	}
}

func isRaspberryPi(osName, arch string) bool {
	// Check if it's running on a Linux-based system with ARM architecture
	return osName == "linux" && strings.HasPrefix(arch, "arm")
}
