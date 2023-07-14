package main

import (
	"encoding/json"
	_ "fmt"
	"io"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof" // only to handle HandleFunc
)

type Book struct {
	Name string `json:"name"`
	Isbn string `json:"isbn"`
}

// it is a slice from book struct and initialize it
var bookList []Book = []Book{}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		b, _ := ioutil.ReadFile("resources/go.png")
		w.Header().Add("content-type", "image/png")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write(b)
		if err != nil {
			return
		}
	})

	mux.Handle("/books", http.HandlerFunc(handleBooks))
	err := http.ListenAndServe(`:8000`, mux)
	if err != nil {
		return
	}
}

func handleBooks(w http.ResponseWriter, r *http.Request) {
	// if username, password, ok := r.BasicAuth() {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// }

	// q := r.URL.Query()
	// fmt.Println(q.Get("name"))

	switch r.Method {
	case "GET":
		b, _ := json.Marshal(bookList)
		w.Write(b)
		return
	case "POST":
		newBook := Book{}
		b, _ := io.ReadAll(r.Body)
		err := json.Unmarshal(b, &newBook)
		if err != nil {
			return
		}

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
		b, err = json.Marshal(bookList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusAccepted)
		_, err = w.Write(b)
		if err != nil {
			return
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
