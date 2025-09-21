package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	// Chapter 4.5: Designing a database model |
	// Import the models package that we just created. You need to prefix this with
	// whatever module path you set up back in chapter 02.01 (Project Setup and Creating
	// a Module) so that the import statement looks like this:
	// "{your-module-path}/internal/models". If you can't remember what module path you
	// used, you can find it at the top of the go.mod file.
	"snippetbox.floccinau.net/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

// Define an application struct to hold the application-wide dependencies for the
// web application. For now we'll only include fields for the two custom loggers, but
// we'll add more to it as the build progresses.
// Chapter 4.5: Designing a database model |
// Add a snippets field to the application struct. This will allow us to
// make the SnippetModel object available to our handlers.
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *models.SnippetModel
}

func main() {
	// Chapter 3.1: Command-line flags |
	// Define a new command-line flag with the name 'addr', a default value of ":4000"
	// and some short help text explaining what the flag controls. The value of the
	// flag will be stored in the addr variable at runtime.
	// example: go run ./cmd/web -addr=":9999"
	// Note: you may use the -help flag to list all the avaliable command-line flags
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Chapter 4.4 Creating a database connection pool |
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")

	// Chapter 3.1: Command-line flags |
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

	// Chapter 4.4: Creating a database connection pool |
	// To keep the main() function tidy I've put the code for creating a connection
	// pool into the separate openDB() function below.We pass openDB() the DSN
	// from the command-line flag.
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	// Chapter 4.4: Creating a database connection pool |
	// We also defer a call to db.Close(), so that the connection pool is closed
	// before the main() function exits.
	defer db.Close()

	// Chapter 3.3: Dependency injection |
	// Initialize a new instance of our application struct, containing the
	// dependencies.
	// Chapter 4.5: Designing a database model |
	// Initialize a models.SnippetModel instance and add it to the application
	// dependecnies.
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		snippets: &models.SnippetModel{DB: db},
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

	// Chapter 4.4: Creating a database connection pool |
	// Because the err variable is now already declared in the code above, we need
	// to use the assignment operator = here, instead of the := 'declare and adsign'
	// operator
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

// Chapter 4.4: Creating a database connection pool |
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
