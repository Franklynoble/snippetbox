package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/Franklynoble/snippetbox/pkg/models/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
)

// Define an application struct to hold the Application-wide dependencies for the
//web application. for now we'll only include field for the two custom logger
// we'll add More to it as the Build progresses.
// Add a snippets field to the application struct. This will allow us to
// make the SnippetModel object available to our handlers

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	session       *sessions.Session
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
}

func main() {

	// Define a new command-line flag with the name 'addr' a default value of ":4000"
	// and some short help text explaining what the flag controls.
	// The Value of the flag will be stored in the addr variable runtime.
	addr := flag.String("addr", ":4000", "HTTP network address")

	dsn := flag.String("dsn", "root:pwd@tcp(0.0.0.0:8083)/snippetbox?parseTime=true", "MYSQL data Source")

	// will be used encrypt and authenticate session cookies) it should be 32
	//bytes long.

	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")

	// importantly, we use the flag.parse() function to parse the command-line flag
	//this reads in the command-line flag value and assigns it to the addr
	//otherwise it will always contain the default value of ":4000" if any errors are
	//encountered during parsing the application will be terminted
	flag.Parse()

	//use  the log.New() to create a  logger for writing information message. This takes
	// three paramertes: the destination to write the  logs to (os.Stdout), a string
	//prefix for message (INFO followed by a tab), and flags to indicate what
	//additional information to include (local date and time). Note that the flags
	//are joined using the bitwise OR operator |.

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	//create a logger for writting error message in the same way, but use stderr as
	// the destination and use the log.Lshortfile flag to include the relevant
	//file name and line number.
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	//initialize a new template cache ...
	templateCache, err := newTemplateCache("./ui/html/")

	if err != nil {
		errorLog.Fatal(err)
	}

	// Initialize a mysql.SnippetModel instance and add it to the application
	// dependencies.
	// Add add it to the application dependencies
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	//initialize a new http.Server struct. We set Addr and Handler fields
	//That the server uses the same newtwork address and routes as before, and
	//the Errorlog fields so That the server now uses the custom errorlog Logger
	// the Event of  any problems

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	// call the listenAndserver() method on our new http.Server Struct.
	err = srv.ListenAndServe()
	errorLog.Fatal(err)

	// err := http.ListenAndServe(*addr, mux)
	// errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	fmt.Println("Connected")
	return db, nil
}
