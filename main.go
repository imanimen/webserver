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

	mux.Handle("/books", handleBooks())
	http.ListenAndServe(`:8000`, nil)
}



func handleBooks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		b, _ := json.Marshal(bookList)
		w.Write(b)
		return
	case "POST":
		newBook := Book{}
		b,_ := io.ReadAll(r.Body)
		json.Unmarshal(b, &newBook)
		for _, b := range bookList {
			if b.Isbn == newBook.Isbn {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
		bookList = append(bookList, newBook)
		w.WriteHeader(http.StatusAccepted)
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}