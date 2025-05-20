package net

import (
	"net/http"
)

func Run() {
	// Use the http.NewServeMux() function to create an empty servemux.
	mux := http.NewServeMux()

	mux.HandleFunc("/hello", HelloWorld)
	http.ListenAndServe(":3000", mux)
}

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}
