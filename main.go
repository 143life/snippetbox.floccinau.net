package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// Define a home handler function which writes a byte slice containing
// "Hello from Snippetbox" as the response body.
func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Hello from Snippetbox"))
}

// snippetView handler function
func snippetView(w http.ResponseWriter, r *http.Request) {
	// Extract the value of the id parameter from the query string and try to
	// convert it to an integer using the strconv.Atoi() function. If it can't
	// be converted to an integer, or the value is less than 1, we return a 404 page
	// not found response
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	// Use the fmt.Fprintf() function to interpolate the id value with our response
	// and write it to the http.ResponseWriter
	fmt.Fprintf(w, "Display a specific snippet with ID %d", id)

	w.Write([]byte("\nDisplay a specific snippet..."))
}

// snippetCreate handler function
func snippetCreate(w http.ResponseWriter, r *http.Request) {
	// Use r.Method to check whether the request is using POST or not.
	if r.Method != "POST" {
		// Use the Header().Set() method to add an 'Allow: POST' header to the
		// response header map. The first parameter is the header name, and
		// the second parameter is the header value.
		// call this method BEFORE calling WriteHeader() and Write()
		w.Header().Set("Allow", http.MethodPost)

		// Set a new cache-control neader. If an existing "Cache-Control" header exists
		// it will be overwritten
		w.Header().Set("Cache-Control", "public, max-age=31536000")
		// In contrast, the Add() method appends a new "Cache-Control" header and can
		// be called multiple times.
		w.Header().Add("Cache-Control", "public")
		w.Header().Add("Cache-Control", "max-age=31536000")
		// Delete all values for the "Cache-Control" header.
		w.Header().Del("Cache-Control")
		//Retrieve the first value for the "Cache-Control" header.
		w.Header().Get("Cache-Control")
		// Retrieve a slice of all values for the "Cache-Control" header.
		w.Header().Values("Cache-Control")
		// If need to avoid textproto.CanonicalMIMEHeaderKey()
		w.Header()["X-XSS-Protection"] = []string{"1; mode=block"}
		// Remove system generated headers
		w.Header()["Date"] = nil

		// // if it's not, use the w.WriteHeader() method to send a 405 status
		// // code and the w.Write() method to write a "Method Not Allowed"
		// // response body. We then return from the function so that the
		// // subsequent code is not executed
		// //w.WriteHeader(405)
		// //w.Write([]byte("Method Not Allowed"))

		// // //w.Header().Set("Content-Type", "application/json")
		// // //w.Write([]byte(`{"name":"Alex"}`))

		// Use the http.Error() function to send a 405 status code and "Method Not
		// Allowed" string as the response body.
		// There we're now passing our http.ResponseWriter to another function,
		// which sends a response to the user of us
		http.Error(w, "MethodNotAllowed", http.StatusMethodNotAllowed)
		return
	}

	w.Write([]byte("Create a new snippet..."))
}

func main() {
	// Use the http.NewServeMux() function to initialize a new servemux, then
	// register the ... function as the handler for the "..." URL pattern
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	// Use the http.ListenAndServe() function to start a new web server. We pass in
	// two parameters: the TCP network address to listen on (in this case ":4000")
	// and the servemux we just created. If http.ListenAndServe() returns an error
	// we use the log.Fatal() function to log the error message and exit. Note
	// that any error returned by http.ListenAndServe() is always non-nil
	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
