package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/gocolly/colly"
)

/*
func fetchForm(url string) error {

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
*/

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
		// fetchForm(sc.Text())
		c := colly.NewCollector()

		c.OnRequest(func(r *colly.Request) {
			fmt.Println("Visiting: ", r.URL)
		})

		c.OnError(func(_ *colly.Response, err error) {
			fmt.Println("Something went wrong: ", err)
		})

		c.OnResponse(func(r *colly.Response) {
			fmt.Println("Page visited: ", r.Request.URL)
		})
		c.OnHTML("div", func(e *colly.HTMLElement) {
			// printing all URLs associated with the a links in the page
			fmt.Println("title is: ", e.Text)
		})

		c.Visit("https://eu.floqast.app/login")
	}

	fmt.Println("fuck new tool!")
}
