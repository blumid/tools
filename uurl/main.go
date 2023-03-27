package main

import (
	"bufio"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func openFile() (*os.File, error) {

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

func getKeysPath(URL string) string {
	u, _ := url.Parse(URL)
	path := u.Scheme + "://" + u.Host + u.Path

	keys := make([]string, 0, len(u.Query()))
	for k := range u.Query() {
		keys = append(keys, k)
	}

	if len(keys) != 0 {
		return path + "?" + strings.Join(keys, "&")
	} else {
		return path
	}

}

func main() {
	// file, err := openFile()

	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	lines := make(map[string]bool)

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		url := strings.TrimSpace(sc.Text())

		res := getKeysPath(url)

		if !lines[res] {
			lines[res] = true
			fmt.Println(url)
		}

	}

}
