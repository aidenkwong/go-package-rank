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

		ib := strings.ReplaceAll(e.ChildText("a"), "Imported by: ", "")
		importedBy, err := strconv.Atoi(strings.ReplaceAll(ib, ",", ""))
		if err != nil {
			panic(err)
		}
		stdPkgs = append(stdPkgs, model.StandardPackage{
			Name:       e.Request.URL.Path,
			ImportedBy: importedBy,
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
		Columns:   []clause.Column{{Name: "name"}},
		DoUpdates: clause.AssignmentColumns([]string{"imported_by"}),
	}).Create(&stdPkgs)

}
