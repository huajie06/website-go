package get13f

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	bolt "go.etcd.io/bbolt"
)

func checkOS() string {
	osName := runtime.GOOS

	// Get the architecture from the GOARCH environment variable
	arch := runtime.GOARCH

	// fmt.Printf("Operating System: %s\n", osName)
	// fmt.Printf("Architecture: %s\n", arch)

	// Check if it's running on a Raspberry Pi
	if isRaspberryPi(osName, arch) {
		// fmt.Println("Running on a Raspberry Pi")
		return "Pi"
	} else if osName == "windows" {
		// fmt.Println("Running on Windows")
		return "Windows"
	} else if osName == "linux" {
		// fmt.Println("Running on Linux (Ubuntu or other)")
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

func setupWD(cik string) {

	wd, err := findGoProjectRoot()
	if err != nil {
		log.Println(err)
	}
	// fmt.Println(wd)

	dir := fmt.Sprintf("%s/db/filings/%s", wd, cik)

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
	cik_json_path := fmt.Sprintf("%s/db/filings/%s/cik.json", wd, cik)

	os.WriteFile(cik_json_path, body, os.ModePerm)

}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func ListBuckets(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println(err)
		return
	}

	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, _ *bolt.Bucket) error {
			fmt.Println(string(name))
			return nil
		})
	})
	if err != nil {
		fmt.Println(err)
		return
	}
}

func ListEachBuckets(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println(err)
		return
	}

	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, _ *bolt.Bucket) error {
			fmt.Printf("Bucket: %s\n", name)
			fmt.Println("=====content below======")

			b := tx.Bucket(name)
			c := b.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				fmt.Printf("key=%s, value=%s\n", k, v)
			}
			fmt.Println("========================")

			// fmt.Println(string(name))
			return nil
		})
	})
	if err != nil {
		fmt.Println(err)
		return
	}
}

func ReadFromBucket(dbPath, bucketName string) {
	db, err := bolt.Open(dbPath, fs.FileMode(os.O_WRONLY), nil)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("key=%s, value=%s\n", k, v)
		}
		return nil
	})
}

func SECtxtURLBuild(cik, accNum string) string {
	accNum_nodash := strings.ReplaceAll(accNum, "-", "")
	txt_url := fmt.Sprintf("https://www.sec.gov/Archives/edgar/data/%s/%s/%s.txt", cik, accNum_nodash, accNum)
	return txt_url
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

func Delete_bucket(db_path, bucket string) {
	db, err := bolt.Open(db_path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte(bucket))
		return err
	})

	fmt.Println("completed")
}

func Delete_key(db_path, bucket, key string) {
	db, err := bolt.Open(db_path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		err := b.Delete([]byte(key))
		return err
	})

	fmt.Println("completed")
}

func Loop_to_2nd(CIK string) {
	wd, err := findGoProjectRoot()
	if err != nil {
		log.Println(err)
	}
	db_dir := fmt.Sprintf("%s/db/filings/fin.db", wd)
	db, err := bolt.Open(db_dir, fs.FileMode(os.O_RDWR), &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	// get existed accNum html
	existed_accNum := []string{}
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(fmt.Sprintf("%s_html", CIK)))
		c := b.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			existed_accNum = append(existed_accNum, string(k))
		}
		return nil
	})

	fmt.Println("existed html", existed_accNum)

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(CIK))
		c := b.Cursor()
		_, c_v := c.First()
		n_k, n_v := c.Next()

		for ; n_k != nil; n_k, n_v = c.Next() {
			fmt.Printf("older: %s. newer: %s\n", c_v, n_v)
			c_v = n_v
		}
		return nil
	})
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// func GetListOfExistedACC() {
// 	wd, err := findGoProjectRoot()
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	db_dir := fmt.Sprintf("%s/db/filings/fin.db", wd)
// 	db, err := bolt.Open(db_dir, fs.FileMode(os.O_RDWR), &bolt.Options{Timeout: 1 * time.Second})
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	defer db.Close()

// }

func GetListOfExistedACC(CIK string) [][]string {
	wd, err := findGoProjectRoot()
	if err != nil {
		log.Println(err)
	}
	db_dir := fmt.Sprintf("%s/db/filings/fin.db", wd)
	db, err := bolt.Open(db_dir, fs.FileMode(os.O_RDWR), &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	pairs_toshow := [][]string{}

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(CIK))
		c := b.Cursor()
		c.First() // skip first
		for n_k, n_v := c.Next(); n_k != nil; n_k, n_v = c.Next() {
			// fmt.Printf("key: %s. value: %s\n", n_k, n_v)
			pairs_toshow = append(pairs_toshow, []string{string(n_k), string(n_v)})
		}
		return nil
	})

	// fmt.Println(pairs_toshow)
	return pairs_toshow
}

func GetTable_fromDB(bucketName, key string) string {
	wd, err := findGoProjectRoot()
	if err != nil {
		log.Println(err)
	}
	db_dir := fmt.Sprintf("%s/db/filings/fin.db", wd)
	db, err := bolt.Open(db_dir, fs.FileMode(os.O_WRONLY), nil)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()
	// fmt.Printf("bucket name: %s, key: %s\n", bucketName, key)

	var html_byte []byte
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		html_byte = b.Get([]byte(key))
		return nil
	})

	return string(html_byte)
}

func DeleteAllKeys_inABucket(bucketName string) {
	wd, err := findGoProjectRoot()
	if err != nil {
		log.Println(err)
	}
	db_dir := fmt.Sprintf("%s/db/filings/fin.db", wd)
	db, err := bolt.Open(db_dir, fs.FileMode(os.O_WRONLY), nil)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		c := b.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			// fmt.Printf("Now deleting...key=%s, value=%s\n", k, v)
			fmt.Printf("Now deleting...key=%s\n", k)
			err := b.Delete(k)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func GetTableKeys_fromDB(bucketName string) {
	wd, err := findGoProjectRoot()
	if err != nil {
		log.Println(err)
	}
	db_dir := fmt.Sprintf("%s/db/filings/fin.db", wd)
	db, err := bolt.Open(db_dir, fs.FileMode(os.O_WRONLY), nil)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()
	fmt.Printf("bucket name: %s\n", bucketName)

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		c := b.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			fmt.Printf("key=%s\n", k)
		}
		return nil
	})

}

func ResetDB(cik string) {
	InitDB(cik)

	fmt.Printf("reset bucket now: %s\n", cik)
	DeleteAllKeys_inABucket(cik)

	fmt.Printf("reset bucket now: %s_html\n", cik)
	DeleteAllKeys_inABucket(fmt.Sprintf("%s_html", cik))
	fmt.Println("reset complete!")
	fmt.Println("===============")
}

func FundList() map[string]string {
	ListofFund := make(map[string]string)
	ListofFund["0001553733"] = "Brave warrior"
	ListofFund["0001709323"] = "Himalaya"
	ListofFund["0001036325"] = "Christopher Davis"
	ListofFund["0001418814"] = "Value Act"
	ListofFund["0001649339"] = "Scion"
	ListofFund["0001099281"] = "Third Avenue"
	ListofFund["0001748240"] = "Soros"
	ListofFund["0001350694"] = "Bridge Water"
	ListofFund["0001358706"] = "Abram"
	ListofFund["0001656456"] = "David Tepper"
	ListofFund["0001061768"] = "Seth Klarman"
	ListofFund["0001603466"] = "Point72"
	ListofFund["0001067983"] = "Berkshire"
	return ListofFund
}

// create a db and all the bucket
func InitDB(CIK string) {
	ListofFund := FundList()
	if _, found := ListofFund[CIK]; !found {
		fmt.Println("not in the list")
		return
	}

	wd, err := findGoProjectRoot()
	if err != nil {
		log.Println(err)
	}

	// initiate a DB for a CIK
	db_dir := fmt.Sprintf("%s/db/filings/fin.db", wd)
	db, err := bolt.Open(db_dir, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	// go read the cik
	local_filing_cik := fmt.Sprintf("%s/db/filings/%s/cik.json", wd, CIK)

	// if not exist, will create folder and download CIK.json from SEC
	if !fileExists(local_filing_cik) {
		setupWD(CIK)
		downloadCIKJson(CIK)
	}

	// once found, it will unmarshal the data and get all the filings
	cik_byte, err := os.ReadFile(local_filing_cik)
	if err != nil {
		log.Println(err)
	}

	var CIKSub_data CIKSub
	if err = json.Unmarshal(cik_byte, &CIKSub_data); err != nil {
		log.Println(err)
	}

	// read the filings and filtered to 13F only
	recent_filings := CIKSub_data.Filings.Recent
	filing_len := len(recent_filings.Form)

	var filtered_13f = []EachFiling{}

	for i := 0; i < filing_len; i++ {
		if recent_filings.Form[i] == "13F-HR" {
			filtered_13f = append(filtered_13f, EachFiling{
				AccessionNumber: recent_filings.AccessionNumber[i],
				FileNumber:      recent_filings.FileNumber[i],
				FilingDate:      recent_filings.FilingDate[i],
				ReportDate:      recent_filings.ReportDate[i],
				PrimaryDocument: recent_filings.PrimaryDocument[i]})
		}
	}

	// create a bucket for that CIK number, for each accNum, to link with the report-filing date
	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(CIK))
		if err != nil {
			return err
		}
		// err = b.Put([]byte("name"), []byte(ListofFund[CIK]))

		for _, v := range filtered_13f[:5] {
			rpt_file_dt := fmt.Sprintf("%s|%s", v.ReportDate, v.FilingDate)
			err = b.Put([]byte(rpt_file_dt), []byte(v.AccessionNumber))
			if err != nil {
				return err
			}
		}
		return nil
	})

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(fmt.Sprintf("%s_html", CIK)))
		if err != nil {
			return err
		}
		return nil
	})
}

func DownloadTXTFiling_v2(CIK string) {
	// read the CIK's last 2 filings and downlaod if not exist
	wd, err := findGoProjectRoot()
	if err != nil {
		log.Println(err)
	}
	db_dir := fmt.Sprintf("%s/db/filings/fin.db", wd)
	db, err := bolt.Open(db_dir, fs.FileMode(os.O_RDWR), &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Println(err)
	}
	defer db.Close()
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(CIK))
		c := b.Cursor()

		_, last_val := c.Last()
		_, prev_val := c.Prev()
		fmt.Println(string(last_val))
		fmt.Println(string(prev_val))
		return nil
	})

}

func DownloadTXTFiling(CIK string) {
	// read the CIK's last 5 filings and downlaod if not exist
	wd, err := findGoProjectRoot()
	if err != nil {
		log.Println(err)
	}
	db_dir := fmt.Sprintf("%s/db/filings/fin.db", wd)
	db, err := bolt.Open(db_dir, fs.FileMode(os.O_RDWR), &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Println(err)
	}
	defer db.Close()
	accNumList_toUpdate := []string{}

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(CIK))
		c := b.Cursor()

		// db is date -> accNum
		_, last_val := c.Last()
		accNumList_toUpdate = append(accNumList_toUpdate, string(last_val))
		for i := 0; i < 5; i++ {
			_, prev_val := c.Prev()
			accNumList_toUpdate = append(accNumList_toUpdate, string(prev_val))
		}
		return nil
	})

	fmt.Printf("List of AccNum to udpate: %s\n", accNumList_toUpdate)
	for _, v := range accNumList_toUpdate {
		filing_txt := fmt.Sprintf("%s/db/filings/%s/%s.txt", wd, CIK, v)
		if !fileExists(filing_txt) {
			txt_url_last := SECtxtURLBuild(CIK, v)
			downloadTXT(txt_url_last, filing_txt)
			time.Sleep(1 * time.Second)
		} else {
			fmt.Println("file already existed!")
		}
	}
}

func RunPython(dir, pre_acc, curr_acc string) []byte {
	wd, err := findGoProjectRoot()
	if err != nil {
		log.Println(err)
	}

	python_script := fmt.Sprintf("%s/app/get13f/parse_txt_to_html.py", wd)

	var python_executable string
	if checkOS() == "Pi" {
		python_executable = "python3"
	} else {
		python_executable = "python"
	}

	cmd := exec.Command(python_executable, python_script, "-d", dir, "-p", pre_acc, "-c", curr_acc)

	// fmt.Println(cmd.Args)
	stdout, err := cmd.Output()
	if err != nil {
		log.Println(err)
		log.Printf("there's some err with the cmd\n%s\n", cmd.Args)
	}
	return stdout
}

func CreateReportHtmltable(CIK string) {
	wd, err := findGoProjectRoot()
	if err != nil {
		log.Println(err)
	}
	db_dir := fmt.Sprintf("%s/db/filings/fin.db", wd)
	db, err := bolt.Open(db_dir, fs.FileMode(os.O_RDWR), &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	// get existed accNum html
	existed_accNum := []string{}
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(fmt.Sprintf("%s_html", CIK)))
		c := b.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			existed_accNum = append(existed_accNum, string(k))
		}
		return nil
	})

	// fmt.Println("existed html: ", existed_accNum)

	txt_filing_dir := fmt.Sprintf("%s/db/filings/%s", wd, CIK)

	update_pairs := [][]string{}
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(CIK))
		c := b.Cursor()
		_, c_v := c.First()
		n_k, n_v := c.Next()

		for ; n_k != nil; n_k, n_v = c.Next() {
			// fmt.Printf("older: %s. newer: %s\n", c_v, n_v)
			if !contains(existed_accNum, string(n_v)) {
				update_pairs = append(update_pairs, []string{string(c_v), string(n_v)})
			}
			c_v = n_v
		}
		return nil
	})

	for _, v := range update_pairs {
		c := v[0]
		n := v[1]
		// fmt.Printf("older: %s. newer: %s\n", c, n)
		fmt.Printf("Updating for %s...\n", c)
		html_byte := RunPython(txt_filing_dir, c, n)

		db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(fmt.Sprintf("%s_html", CIK)))
			err := b.Put([]byte(n), html_byte)
			return err
		})
	}
}

func UpdateOneCIK(cik string) {
	ListofFund := FundList()
	fmt.Printf("====Setting up DB for: %s. CIK key: %s====\n", ListofFund[cik], cik)
	InitDB(cik)
	DownloadTXTFiling(cik)
	time.Sleep(1 * time.Second)
	fmt.Printf("====SUCCESS set up DB for: %s. CIK key: %s===\n", ListofFund[cik], cik)
	fmt.Printf("====generating html table for company: %s...====\n", ListofFund[cik])
	CreateReportHtmltable(cik)
}

func Wrapper_v1() {
	ListofFund := FundList()
	for cik, company_name := range ListofFund {
		fmt.Printf("#############Updating company %s...#############", company_name)
		UpdateOneCIK(cik)
	}
}
