package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
)

func fetchForm(doc *goquery.Document) []endpoint {

	// search form tags:

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

	return result
}

func makeRequest(url string) error {

	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		return err
	}
	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println("body is: ", string(body))

	// doc, err := html.Parse(resp.Body)
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	data := fetchForm(doc)
	fmt.Println("data is: ", data)
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
		if err := makeRequest(sc.Text()); err != nil {
			fmt.Println("err is: ", err)
			continue
		}

	}

	fmt.Println("fuck new tool!")
}
