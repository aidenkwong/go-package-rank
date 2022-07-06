package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/aidenkwong/go-package-rank/pkg/model"
)

func main() {

	http.HandleFunc("/std", func(w http.ResponseWriter, r *http.Request) {
		stdPkgs := model.GetAllStandardPackages()
		sort.Slice(stdPkgs, func(p, q int) bool {
			return stdPkgs[p].NumOfImports > stdPkgs[q].NumOfImports
		})
		stdPkgsJSON, _ := json.Marshal(stdPkgs)
		fmt.Fprintf(w, "%v", string(stdPkgsJSON))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))

}
