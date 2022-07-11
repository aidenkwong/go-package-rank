package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

type StandardPackage struct {
	Name       string
	ImportedBy int
}

func main() {

	getLinksCollector := colly.NewCollector()

	var stdPkgLinks []string

	getLinksCollector.OnHTML("table", func(e *colly.HTMLElement) {

		e.ForEach("tr > td:nth-child(1) > div > span > a", func(_ int, e *colly.HTMLElement) {
			stdPkgLinks = append(stdPkgLinks, e.Attr("href"))
		})

		e.ForEach("tr > td:nth-child(1) > div > div > a", func(_ int, e *colly.HTMLElement) {
			stdPkgLinks = append(stdPkgLinks, e.Attr("href"))
		})
	})

	getLinksCollector.Visit("https://pkg.go.dev/std")

	fmt.Println(strconv.Itoa(len(stdPkgLinks)) + " packages found in std")

	visitLinksCollector := colly.NewCollector()

	var stdPkgs []StandardPackage

	visitLinksCollector.OnHTML("body > main > header > div > div.go-Main-headerDetails > span:nth-child(5)", func(e *colly.HTMLElement) {

		ib := strings.ReplaceAll(e.ChildText("a"), "Imported by: ", "")
		importedBy, err := strconv.Atoi(strings.ReplaceAll(ib, ",", ""))
		if err != nil {
			panic(err)
		}
		stdPkgs = append(stdPkgs, StandardPackage{
			Name:       e.Request.URL.Path,
			ImportedBy: importedBy,
		})

	})

	maxGoroutines := 10
	guard := make(chan struct{}, maxGoroutines)
	for i, link := range stdPkgLinks {
		guard <- struct{}{}
		go func(i int, link string) {
			println(i, " of ", len(stdPkgLinks), " ", link)
			err := visitLinksCollector.Visit("https://pkg.go.dev" + link)
			if err != nil {
				fmt.Println(err)
			}
			<-guard
		}(i, link)

	}
	sort.Slice(stdPkgs, func(i, j int) bool {
		return stdPkgs[i].ImportedBy > stdPkgs[j].ImportedBy
	})
	stdPkgsJSON, err := json.Marshal(stdPkgs)
	if err != nil {
		println(err)
	}
	os.WriteFile("std_pkg.json", []byte(stdPkgsJSON), 0644)

}
