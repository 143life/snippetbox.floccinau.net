package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// Chapter 3.4: Centralized handling |
// The serverError helper writes an error message and stack trace to the errorLog,
// then sends a generic 500 Internal Server Error response to the user.
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	// 2 cause we need error message from file when error appeared,
	// not from this file.
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// Chapter 3.4: Centralized error handling |
// The clientError helper sends a specific status code and corresponding description
// to the user. We'll use this later to send responses like 400 "Bad Request"
// when there's a problem with the request that the user sent.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// Chapter 3.4: Centralized error handling |
// For consistency, we'll also implement a notFound helper. This is simply a
// convenience wrapper around clientError which sends a 404 Not Found response to
// the user.
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
