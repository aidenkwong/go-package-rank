package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aidenkwong/go-package-rank/pkg/model"
	"github.com/gocolly/colly"
	"gorm.io/gorm/clause"
)

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

	var stdPkgs []model.StandardPackage

	visitLinksCollector.OnHTML("body > main > header > div > div.go-Main-headerDetails > span:nth-child(5)", func(e *colly.HTMLElement) {

		noi := strings.ReplaceAll(e.ChildText("a"), "Imported by: ", "")
		numOfImports, err := strconv.Atoi(strings.ReplaceAll(noi, ",", ""))
		if err != nil {
			panic(err)
		}
		stdPkgs = append(stdPkgs, model.StandardPackage{
			Name:         e.Request.URL.Path,
			NumOfImports: numOfImports,
		})

	})

	for _, link := range stdPkgLinks {

		err := visitLinksCollector.Visit("https://pkg.go.dev" + link)
		if err != nil {
			fmt.Println(err)
		}

	}
	db := model.GetDB()
	db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&stdPkgs)

}
