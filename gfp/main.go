package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"golang.org/x/exp/slices"
)

/*
# install headless chrome:
sudo apt update
sudo apt install -y libappindicator1 fonts-liberation
sudo apt install -f
sudo apt --fix-broken install
wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb
sudo dpkg -i google-chrome-stable_current_amd64.deb
*/

var result = make(map[string]bool)

func makeRequest(url string) error {

	geziyor.NewGeziyor(&geziyor.Options{
		LogDisabled: true,
		StartRequestsFunc: func(g *geziyor.Geziyor) {
			g.GetRendered(url, g.Opt.ParseFunc)
		},
		ParseFunc: func(g *geziyor.Geziyor, r *client.Response) {
			if slices.Contains(r.Header["Content-Type"], "application/javascript") {
				fmt.Println("fuck, it's a js file!")
				fetchXHR(string(r.Body))
			} else {
				if doc := r.HTMLDoc; doc != nil {
					fetchForm(doc)
				}
			}

		},
	}).Start()

	// fmt.Println("result is: ", result)
	return nil
}

func fetchForm(doc *goquery.Document) {
	result := make([]endpoint, 0)
	doc.Find("form").Each(func(i int, form *goquery.Selection) {
		ep := endpoint{}
		ep.ref = "form"
		ep.method = form.AttrOr("method", "GET")
		ep.url = form.AttrOr("action", "/")
		ep.c_type = form.AttrOr("enctype", "x-www-form-urlencoded")
		ep.params = make([]string, 1)

		form.Find("input").Each(func(i int, input *goquery.Selection) {

			param := input.AttrOr("name", "")
			ep.params = append(ep.params, param)
			// ep.params[param] = input.AttrOr("value", "")
		})

		makeMap(ep)
		result = append(result, ep)
		fmt.Println("result is:", result)
	})

	doc.Find("script").Each(func(i int, scr *goquery.Selection) {
		if code := scr.Text(); code != "" {
			fetchXHR(scr.Text())
		}

	})

}

func fetchXHR(code string) error {
	re, _ := regexp.Compile(`var ([a-zA-Z1-9]+)\s?=\s?new XMLHttpRequest`)
	varName := re.FindStringSubmatch(code)
	if varName != nil {
		fmt.Println("var name is: ", varName[1])
		re2, _ := regexp.Compile(varName[1] + `\.open\((.*)\)`)
		// fmt.Println("regex is: ", re2)

		if url := re2.FindStringSubmatch(code); url != nil {
			re3, _ := regexp.Compile(varName[1] + `\.send\((.*)\)`)
			if body := re3.FindStringSubmatch(code); body != nil {
				fmt.Println("body is: ", body)
			}
			fmt.Println("url is: ", strings.Split(url[1], ",")[1])

		}
	}

	return nil
}

func makeMap(ep endpoint) {

	key := ep.method + "_" + ep.url + "_" + strings.Join(ep.params, ",") + "_" + ep.ref
	fmt.Println("key is:", key)
	if !result[key] {
		fmt.Println("fuck!! we found a new one.")
	}
	//
}

type endpoint struct {
	url string
	// params map[string]string
	params []string
	method string //get|post|put|delete
	c_type string //form|json|file
	ref    string //form|XHR|Jquery|Axios
}

func main() {

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		makeRequest(sc.Text())
	}
}
