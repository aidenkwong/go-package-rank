package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"

	"github.com/aidenkwong/go-package-rank/pkg/model"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	http.HandleFunc("/std", func(w http.ResponseWriter, r *http.Request) {
		stdPkgs := model.GetAllStandardPackages()
		sort.Slice(stdPkgs, func(p, q int) bool {
			return stdPkgs[p].ImportedBy > stdPkgs[q].ImportedBy
		})
		stdPkgsJSON, _ := json.Marshal(stdPkgs)
		fmt.Fprintf(w, "%v", string(stdPkgsJSON))
	})

	log.Fatal(http.ListenAndServe(":"+port, nil))
	log.Println("Listening on port ", port)
}
