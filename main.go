package main

import (
	"io/ioutil"
	"net/http"
)

func main() {

	mux := http.DefaultServeMux
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		b, _ := ioutil.ReadFile("resources/go.png")
		w.Header().Add("content-type", "image/png")
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	})
	http.ListenAndServe(`:8000`, nil)
	
}