package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
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

func fetchForm(url string) error {

	result := make([]endpoint, 0)

	geziyor.NewGeziyor(&geziyor.Options{
		LogDisabled: true,
		StartRequestsFunc: func(g *geziyor.Geziyor) {
			g.GetRendered(url, g.Opt.ParseFunc)
		},
		ParseFunc: func(g *geziyor.Geziyor, r *client.Response) {
			doc := r.HTMLDoc

			if doc != nil {
				doc.Find("form").Each(func(i int, form *goquery.Selection) {
					ep := endpoint{}
					ep.ref = "form"
					ep.method = form.AttrOr("method", "GET")
					ep.url = form.AttrOr("action", "/")
					ep.c_type = form.AttrOr("enctype", "x-www-form-urlencoded")
					ep.params = make(map[string]string)

					form.Find("input").Each(func(i int, input *goquery.Selection) {

						param := input.AttrOr("name", "")
						ep.params[param] = input.AttrOr("value", "")
					})

					result = append(result, ep)
					fmt.Println("result is:", result)
				})

			}
		},
	}).Start()

	// fmt.Println("result is: ", result)
	return nil
}

func fetchJquery(url string) error {

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
}
