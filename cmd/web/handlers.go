package main

import "net/http"

func home(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "Welcome to Gratitude Jar",
	}
	render(w, r, "home.tmpl", data)
}
