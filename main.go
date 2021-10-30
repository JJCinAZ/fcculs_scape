package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"

	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/mgutz/ansi"
	"golang.org/x/net/html"
)

const ULSURI = "http://wireless2.fcc.gov/UlsApp/UlsSearch/licenseLocSum.jsp?licKey=%d"

func main() {
	var (
		licenseKey int
	)
	flag.IntVar(&licenseKey, "l", 0, "License key to retrieve")
	flag.Parse()
	if licenseKey == 0 {
		flag.Usage()
		os.Exit(1)
	}
	err := getLicense2(licenseKey)
	if err != nil {
		fmt.Println(err)
	}
}

func getLicense2(licenseKey int) error {
	url := fmt.Sprintf(ULSURI, licenseKey)
	c := colly.NewCollector()

	color1 := ansi.ColorCode("yellow")
	//color2 := ansi.ColorCode("blue+h")
	c.OnHTML(`table[summary="graphical layout table"]`, func(e *colly.HTMLElement) {
		fmt.Printf("%d\n", len(e.DOM.ChildrenFiltered("tbody>tr>td").Nodes))
		e.DOM.ChildrenFiltered("tbody").Each(func(i int, selection *goquery.Selection) {
			if selection.Nodes[0].Type == html.ElementNode {
				fmt.Printf("#%d-%s: %s%s%s\n", i, selection.Nodes[0].Data, color1, strings.TrimSpace(selection.Text()), ansi.DefaultFG)
			}
		})
	})
	/*
		if summary, found := e.DOM.Attr("summary"); found {
			switch summary {
			case "Location info table":
				e.ForEach("td", func(i int, element *colly.HTMLElement) {
					fmt.Printf("#%d: %s%s%s\n", i, color1, strings.TrimSpace(element.Text), ansi.DefaultFG)
				})
			case "graphical layout table":
				e.ForEach("td", func(i int, element *colly.HTMLElement) {
					fmt.Printf("#%d: %s%s%s\n", i, color2, strings.TrimSpace(element.Text), ansi.DefaultFG)
				})
			}
		}
	*/

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit(url)
	return nil
}

func getLicense(licenseKey int) error {
	url := fmt.Sprintf(ULSURI, licenseKey)
	fmt.Println(url)
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	if res == nil {
		return errors.New("response is nil")
	}
	defer res.Body.Close()
	if res.Request == nil {
		return errors.New("response.Request is nil")
	}
	fmt.Println("Reading body")
	doc, err := goquery.NewDocumentFromReader(res.Body)
	fmt.Println("find tables")
	doc.Find("table").Each(func(index int, item *goquery.Selection) {
		summary, _ := item.Attr("summary")
		fmt.Printf("Table #%d: %s\n", index, summary)
	})
	return nil
}
