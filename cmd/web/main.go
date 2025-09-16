package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// Define an application struct to hold the application-wide dependencies for the
// web application. For now we'll only include fields for the two custom loggers, but
// we'll add more to it as the build progresses.
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	// Chapter 3.1: Command-line flags
	// Define a new command-line flag with the name 'addr', a default value of ":4000"
	// and some short help text explaining what the flag controls. The value of the
	// flag will be stored in the addr variable at runtime.
	// example: go run ./cmd/web -addr=":9999"
	// Note: you may use the -help flag to list all the avaliable command-line flags
	addr := flag.String("addr", ":4000", "HTTP network address")
	// Importantly, we use the flag.Parse() function to parse the command-line flag.
	// This reads in the command-line flag value and assigns it to the addr
	// variable. You need to call this *before* you use the addr variable
	// otherwise it vill always contain the default value of ":4000". If any errors are
	// encountered during parsing the application will be terminated.
	flag.Parse()

	// Chapter 3.2: Leveled logging
	// Use log.New() to create a logger for writing information messages. This takes
	// three parameters: the destination to write the logs to (os.Stdout), a string
	// prefix for message (INFO followed by a tab), and flags to indicate what
	// additional information to include (local date and time). Note that the flags
	// are joined using the bitwise OR operator |.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	// Create a logger for writing error messages in the same way, but use stderr as
	// the destination and use the log.Lshortfile flag to include the relevant
	// file name and line number.
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Chapter 3.3: Dependency injection
	// Initialize a new instance of our application struct, containing the
	// dependencies.
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	// Chapter 3.2: The http.Server error log
	// Initialize a new http.Server struct. We set the Addr and Handler fields so
	// that the server uses the same network address and routes before, and set
	// the ErrorLog field so that the server now uses the custom errorLog logger in
	// the event of any problems.
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		// Chapter 3.5: Isolating the application routes |
		Handler: app.routes(),
	}

	// The value returned from the flag.String() is a pointer to the flag
	// value, not the value itself. So we need to dereference the pointer (i.e.
	// prefix it with the * symbol) before using it. Note that we're using the
	// log.Printf() function to interpolate the address with the log message.
	infoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
