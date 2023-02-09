package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

type Package struct {
	Name       string
	ImportedBy int
	GitHubStar int
}

func GetStringInBetweenTwoString(str string, startS string, endS string) (result string, found bool) {
	s := strings.Index(str, startS)
	if s == -1 {
		return result, false
	}
	newS := str[s+len(startS):]
	e := strings.Index(newS, endS)
	if e == -1 {
		return result, false
	}
	result = newS[:e]
	return result, true
}

func main() {
	GITHUB_TOKEN := os.Getenv("GITHUB_TOKEN")
	if GITHUB_TOKEN == "" {
		panic("GITHUB_TOKEN is empty")
	}
	file, err := os.Open("awesome_go.md")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var links []string

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		str, boolean := GetStringInBetweenTwoString(scanner.Text(), "(https://", ")")
		if boolean && strings.HasPrefix(str, "github.com") {
			links = append(links, str)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	totalNumOfLinks := len(links)
	println("Total:", totalNumOfLinks)

	visitLinksCollector := colly.NewCollector()

	var pkgs []Package

	visitLinksCollector.OnHTML("body > main > header > div > div.go-Main-headerDetails > span:nth-child(5)", func(e *colly.HTMLElement) {
		if strings.HasPrefix(e.Request.URL.String(), "https://pkg.go.dev") {
			ib := strings.ReplaceAll(e.ChildText("a"), "Imported by: ", "")
			importedBy, err := strconv.Atoi(strings.ReplaceAll(ib, ",", ""))
			if err != nil {
				log.Println(err)
			} else {

				str := strings.Replace(e.Request.URL.Path[1:], "github.com/", "", 1)
				client := &http.Client{}
				URL := "https://api.github.com/repos/" + str
				req, err := http.NewRequest("GET", URL, nil)

				req.SetBasicAuth("aidenkwong", GITHUB_TOKEN)
				resp, err := client.Do(req)
				if err != nil {
					panic(err)
				}
				bodyText, err := ioutil.ReadAll(resp.Body)

				var body map[string]interface{}
				json.Unmarshal(bodyText, &body)
				var githubStar float64
				var ok bool
				if x, found := body["watchers_count"]; found {
					if githubStar, ok = x.(float64); !ok {
						//do whatever you want to handle errors - this means this wasn't a string
					}
				} else {
					// log.Println(URL)
					// log.Println(string(bodyText))
					// log.Fatal("watchers_count not found")
				}
				pkgs = append(pkgs, Package{
					Name:       e.Request.URL.Path[1:],
					ImportedBy: importedBy,
					GitHubStar: int(githubStar),
				})
			}
		}

	})
	maxGoroutines := 5
	guard := make(chan struct{}, maxGoroutines)

	print("\033[s")
	progress := 0
	for i, link := range links {

		guard <- struct{}{}

		go func(i int, link string) {
			print("\033[u\033[K")
			curProgress := (i + 1) * 100 / totalNumOfLinks
			if curProgress > progress {
				progress = curProgress
			}
			print("Progress: ", progress, "% ", link)
			err := visitLinksCollector.Visit("https://pkg.go.dev/" + link)
			if err != nil {
				log.Println(err)
			}
			<-guard
		}(i, link)

	}
	sort.Slice(pkgs, func(i, j int) bool {
		return pkgs[i].ImportedBy > pkgs[j].ImportedBy
	})
	pkgsJSON, err := json.Marshal(pkgs)
	if err != nil {
		log.Fatal(err)
	}
	os.WriteFile("public/src/awesome_go.json", []byte(pkgsJSON), 0644)

}
