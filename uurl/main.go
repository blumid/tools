package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

func openFile() (io.ReadWriteCloser, error) {
	fmt.Println("opening file:")

	pwd, err := filepath.Abs(".")
	if err != nil {
		return nil, err
	}

	file, err := os.Open(filepath.Join(pwd, "pattern"))
	if err != nil {
		return nil, errors.New("can't open pattern file")
	}

	return file, nil

}
func getKeysPath(URL string) ([]string, string) {
	u, _ := url.Parse(URL)
	path := u.Scheme + u.Host + "/" + u.Path

	keys := make([]string, 0, len(u.Query()))
	for k, _ := range u.Query() {
		keys = append(keys, k)
	}

	return keys, path

}

func checkPattern(file io.ReadWriteCloser, URL string) bool {

	keys, path := getKeysPath(URL)

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		patt_keys, patt_path := getKeysPath(sc.Text())
		if path == patt_path {
			if reflect.DeepEqual(keys, patt_keys) {
				return false
			}
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

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	sc := bufio.NewScanner(os.Stdin)

	for sc.Scan() {
		url := strings.TrimSpace(sc.Text())

		ck := checkPattern(file, url)
		if ck {
			fmt.Println(url)
		}

	}
	file.Close()

}
