package main

import "net/http"

func routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/gratitude", gratitude)
	mux.HandleFunc("/gratitude/create", createGratitude)
	return mux
}
