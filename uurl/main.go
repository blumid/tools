package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strings"
)

func makeKey(URL string) (string, error) {
	u, err := url.Parse(URL)

	if err != nil {
		fmt.Println("we some erros: ", err)
		return "", err
	}
	path := u.Scheme + "://" + u.Host + u.Path

	keys := make([]string, len(u.Query()))
	for k := range u.Query() {
		keys = append(keys, k)
	}

	if len(keys) != 0 {
		return path + "?" + strings.Join(keys, "&"), nil
	} else {
		return path, nil
	}

}

func main() {
	lines := make(map[string]bool)

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		url := strings.TrimSpace(sc.Text())

		res, err := makeKey(url)

		if err != nil {
			continue
		}

		if !lines[res] {
			lines[res] = true
			fmt.Println(url)
		}

	}

}
