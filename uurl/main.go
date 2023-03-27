package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strings"
)

func makeKey(URL string) string {
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
	lines := make(map[string]bool)

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		url := strings.TrimSpace(sc.Text())

		res := makeKey(url)

		if !lines[res] {
			lines[res] = true
			fmt.Println(url)
		}

	}

}
