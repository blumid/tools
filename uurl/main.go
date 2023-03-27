package main

import (
	"bufio"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

func openFile() (*os.File, error) {
	// fmt.Println("opening file:")

	pwd, err := filepath.Abs(".")
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(filepath.Join(pwd, "pattern"), os.O_RDWR|os.O_CREATE, 0660)
	// file, err := os.Open(filepath.Join(pwd, "pattern"))
	if err != nil {
		return nil, errors.New("can't open pattern file")
	}

	return file, nil

}
func getKeysPath(URL string) ([]string, string) {
	u, _ := url.Parse(URL)
	path := u.Scheme + "://" + u.Host + u.Path

	// fmt.Println("path is:", path)

	keys := make([]string, 0, len(u.Query()))
	for k := range u.Query() {
		keys = append(keys, k)
	}
	// fmt.Println("keys are:", keys)

	return keys, path

}

func addTo(file *os.File, URL string) bool {

	keys, path := getKeysPath(URL)

	sc := bufio.NewScanner(file)
	sc.Split(bufio.ScanLines)

	fmt.Println("sc scan is:", sc.Scan())
	for sc.Scan() {
		fmt.Println("in side of addTo, for scan")
		patt_keys, patt_path := getKeysPath(sc.Text())
		if path == patt_path {
			fmt.Println("first condition: ", path)
			if reflect.DeepEqual(keys, patt_keys) {
				fmt.Println("duplicate happend.", URL)
				return false
			}
		} else {
			continue
		}
	}
	return true
}

type pattern struct {
	hname string
	query []string
}

func main() {

	file, err := openFile()
	// dw := bufio.NewWriter(file)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		url := strings.TrimSpace(sc.Text())

		if addTo(file, url) {
			// dw.WriteString(url + "\n")
			if _, err := file.WriteString(url + "\n"); err != nil {
				fmt.Println("err is: ", err)
			}
		}

	}
	// defer dw.Flush()
	defer file.Close()

}
