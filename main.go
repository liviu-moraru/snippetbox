package main

import (
	"log"
	"net/http"
)

// Define a home handler function which writes a byte slice containing
// "Hello from Snippetbox" as the response body.
func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		log.Println("Request path:", r.URL.Path)
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Hello from Snippetbox"))
}

// Add a snippetView handler function.
func snippetView(w http.ResponseWriter, r *http.Request) {
	testHeaderMap(w)
	w.Write([]byte("Display a specific snippet..."))
}

func testHeaderMap(w http.ResponseWriter) {
	//w.Header().Set("Cache-Control", "public, max-age=31536000")
	w.Header().Add("Cache-Control", "public")
	w.Header().Add("Cache-control", "max-age-31536000")
	//w.Header().Del("Cache-Control")
	w.Header()["Date"] = nil
	log.Println("Cache-Control:", w.Header().Get("Cache-Control"))
	log.Println("Cache-Controls:", w.Header().Values("Cache-Control"))
}

// Add a snippetCreate handler function.
func snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method now allowed", http.StatusMethodNotAllowed)
		/*w.WriteHeader(405)
		w.Write([]byte("Method not allowed"))*/
	}
	w.Write([]byte("Create a new snippet..."))
}

func main() {
	// Use the http.NewServeMux() function to initialize a new servemux, then
	// register the home function as the handler for the "/" URL pattern.
	mux := http.NewServeMux()

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)
	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)

}
