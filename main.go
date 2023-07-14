package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)


type Book struct {   
	Name string `"json:name"`
	Isbn string `"json:isbn"`
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

	mux.Handle("/books", http.HandlerFunc(handleBooks))
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
        b, _ := io.ReadAll(r.Body)
        json.Unmarshal(b, &newBook)

        // Check if book with same ISBN already exists
        for _, b := range bookList {
            if b.Isbn == newBook.Isbn {
                w.WriteHeader(http.StatusBadRequest)
                return
            }
        }

        // Add new book to the bookList slice
        bookList = append(bookList, newBook)
        
        // Return updated bookList in response body
        b, err := json.Marshal(bookList)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            return
        }
        w.WriteHeader(http.StatusAccepted)
        w.Write(b)
    default:
        w.WriteHeader(http.StatusBadRequest)
        return
    }
}