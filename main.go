package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)


type Book struct {   
	Name string `"json:isbn`
	Isbn string `"json:name`
}

// it is a slice from book struct and initialize it
var bookList []Book = []Book{}

func main() {

	mux := http.DefaultServeMux
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		b, _ := ioutil.ReadFile("resources/go.png")
		w.Header().Add("content-type", "image/png")
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	})

	mux.Handle("/books")
	http.ListenAndServe(`:8000`, nil)
}



func handleBooks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		b, _ := json.Marshal(bookList)
		w.Write(b)
		return
	case "POST":
		io.ReadAll(r.Body)

	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}