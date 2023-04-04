package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/browser"
	"gopkg.in/headzoo/surf.v1"
)

func fetchForm(url string) error {

	// search form tags:

	browser.OpenURL(url)

	doc := surf.NewBrowser()
	doc.AddRequestHeader("Accept", "text/html")
	doc.AddRequestHeader("Accept-Charset", "utf8")
	err := doc.Open(url)
	fmt.Println(string(doc.Body()))

	if err != nil {
		return err
	}

	result := make([]endpoint, 1)

	doc.Find("form").Each(func(i int, form *goquery.Selection) {
		fmt.Println("find a form")
		ep := endpoint{}
		ep.ref = "form"
		ep.method = form.AttrOr("method", "GET")
		ep.url = form.AttrOr("action", "/")
		ep.c_type = form.AttrOr("enctype", "application/x-www-form-urlencoded")

		form.Find("input").Each(func(i int, input *goquery.Selection) {

			param := input.AttrOr("name", "")
			ep.params[param] = input.AttrOr("value", "")
		})

		result = append(result, ep)

	})

	fmt.Println("result is: ", result)
	return nil
}

type endpoint struct {
	url    string
	params map[string]string
	method string //get|post|put|delete
	c_type string
	ref    string //form|XHR|Axios|Jquery
}

func main() {

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		fetchForm(sc.Text())

	}

	fmt.Println("fuck new tool!")
}
