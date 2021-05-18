package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/tendermint/tendermint/libs/bech32"
)

type Page struct {
	Body   []byte
	Result []byte
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Panic("PORT is not set")
	}
	addr := fmt.Sprintf(":%s", port)

	http.HandleFunc("/", handler)
	http.HandleFunc("/validate/", validateHandler)

	log.Printf("Listening %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	page := &Page{Body: []byte("Paste addresses")}

	tmpl, _ := template.ParseFiles("index.html")
	tmpl.Execute(w, page)
}

func validateHandler(w http.ResponseWriter, r *http.Request) {
	body := r.FormValue("body")
	page := &Page{Body: []byte(body)}

	var result strings.Builder

	for i, line := range strings.Split(body, "\n") {
		addr := strings.TrimSpace(line)
		if len(addr) == 0 {
			continue
		}

		if i > 0 {
			result.WriteString("\n")
		}

		hrp, _, err := bech32.DecodeAndConvert(addr)
		if err != nil || hrp != "panacea" {
			result.WriteString(fmt.Sprintf("%s,NO", addr))
		} else {
			result.WriteString(fmt.Sprintf("%s,YES", addr))
		}
	}

	page.Result = []byte(result.String())

	tmpl, _ := template.ParseFiles("index.html")
	tmpl.Execute(w, page)
}
