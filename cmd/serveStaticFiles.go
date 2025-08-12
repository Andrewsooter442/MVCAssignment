package main

import "net/http"

func ServeStaticFiles(mux *http.ServeMux) {
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
}
