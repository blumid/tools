package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"sort"
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
				fetchXHR(string(r.Body), url)
			} else {
				if doc := r.HTMLDoc; doc != nil {
					fetchForm(doc, url)
				}
			}

		},
	}).Start()

	// fmt.Println("result is: ", result)
	return nil
}

func fetchForm(doc *goquery.Document, url string) {
	result := make([]endpoint, 0)
	doc.Find("form").Each(func(i int, form *goquery.Selection) {
		ep := endpoint{}
		ep.ref = "html"
		ep.method = form.AttrOr("method", "GET")
		ep.url = url + form.AttrOr("action", "/")
		ep.c_type = form.AttrOr("enctype", "form")
		if strings.Contains(ep.c_type, "form") {
			ep.c_type = "form"
		} else {
			ep.c_type = "file"
		}

		ep.params = make([]string, 1)

		form.Find("input").Each(func(i int, input *goquery.Selection) {

			param := input.AttrOr("name", "")
			ep.params = append(ep.params, param)
			// ep.params[param] = input.AttrOr("value", "")
		})

		makeMap(ep)
		result = append(result, ep)

	})

	doc.Find("script").Each(func(i int, scr *goquery.Selection) {
		if code := scr.Text(); code != "" {
			fetchXHR(scr.Text(), url)
		}

	})

}

func fetchXHR(code string, url string) error {
	ep := endpoint{}
	ep.url = url
	re, _ := regexp.Compile(`var ([a-zA-Z1-9]+)\s?=\s?new XMLHttpRequest`)
	varName := re.FindStringSubmatch(code)
	if varName != nil {
		re2, _ := regexp.Compile(varName[1] + `\.open\([\"\'](\S+)[\"\'],[\"\'](\S+)[\"\'].*\)`)

		header, _ := regexp.Compile(varName[1] + `\.setRequestHeader\(\)`)
		fmt.Println("the header is: ", header)

		if open := re2.FindStringSubmatch(code); open != nil {
			fmt.Println("the method is: ", open[1])
			fmt.Println("the url is: ", open[2])
			ep.method = open[1]
			ep.url = open[2]

			if !slices.Contains([]string{"GET", "HEAD"}, strings.ToUpper(ep.method)) {
				re3, _ := regexp.Compile(varName[1] + `\.send\((.*)\)`)
				if body := re3.FindStringSubmatch(code); body != nil {
					fmt.Println("body type is: ", reflect.TypeOf(body[1]))
					ep.params = append(ep.params, body[1])
				}
			} else {
				ep.params = nil
			}

		}
	}

	return nil
}

func makeMap(ep endpoint) {

	// make paramters string sorted.
	sort.Strings(ep.params)

	/*key structure is:

	method$url$c-type$ref$SortedJoinedParams(&)

	*/

	// make key string
	key := ep.method + "$" + ep.url + "$" + ep.c_type + "$" + strings.Join(ep.params, "=&") + "$" + ep.ref
	if !result[key] {
		result[key] = true
	}

}

type endpoint struct {
	url string
	// params map[string]string ,i comment it and had repalced with the below line
	params []string
	method string //get|post|put|delete
	c_type string //query|form|file|json
	ref    string //html|XHR|Jquery|Axios
}

func main() {

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		makeRequest(sc.Text())
	}
}
