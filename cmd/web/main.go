package main

import (
	"flag"
	"log"
	"net/http"
)

type config struct {
	addr      string
	staticDir string
}

var cfg config

func main() {

	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")

	//addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()
	// Use the http.NewServeMux() function to initialize a new servemux, then
	// register the home function as the handler for the "/" URL pattern.
	mux := http.NewServeMux()

	// Create a file server which serves files out of the "./ui/static" directory.
	// Note that the path given to the http.Dir function is relative to the project
	// directory root.
	fileServer := http.FileServer(neuteredFileSystem{http.Dir(cfg.staticDir)})

	// Use the mux.Handle() function to register the file server as the handler for
	// all URL paths that start with "/static/". For matching paths, we strip the
	// "/static" prefix before the request reaches the file server.
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)
	mux.Handle("/handler/", &handlerImpl{})
	log.Printf("Starting server on %s\n", cfg.addr)
	err := http.ListenAndServe(cfg.addr, mux)
	log.Fatal(err)

}
