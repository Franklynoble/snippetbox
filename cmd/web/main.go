package main

import (
	// "crypto/tls"
	// "database/sql"
	// "flag"
	// "fmt"
	// "html/template"
	// "log"
	// "net/http"
	// "os"
	// "time"

	"crypto/tls" // New import
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Franklynoble/snippetbox/pkg/models"
	"github.com/Franklynoble/snippetbox/pkg/models/mysql"
	"github.com/golangcollege/sessions"
	// "github.com/Franklynoble/snippetbox/pkg/models/mysql"
	// _ "github.com/go-sql-driver/mysql"
	// "github.com/golangcollege/sessions"
)

// Define an application struct to hold the Application-wide dependencies for the
//web application. for now we'll only include field for the two custom logger
// we'll add More to it as the Build progresses.
// Add a snippets field to the application struct. This will allow us to
// make the SnippetModel object available to our handlers

type contextKey string

const contextKeyIsAuthenticated = contextKey("isAuthenticated")

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	session  *sessions.Session
	snippets interface {
		Insert(string, string, string) (int, error)
		Get(int) (*models.Snippet, error)
		Latest() ([]*models.Snippet, error)
	}
	//snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
	users         interface {
		Insert(string, string, string) error
		Authenticate(string, string) (int, error)
		Get(int) (*models.User, error)
		ChangePassword(int, string, string) error
	}
	//users *mysql.UserModel //add new users field to the application struct

}

func main() {

	// Define a new command-line flag with the name 'addr' a default value of ":4000"
	// and some short help text explaining what the flag controls.
	// The Value of the flag will be stored in the addr variable runtime.
	addr := flag.String("addr", ":4000", "HTTP network address")

	dsn := flag.String("dsn", "secret:root@tcp(0.0.0.0:3306)/snippetbox?parseTime=true", "MYSQL data Source")

	//dsn := flag.String("dsn", "root:secret@localhost(0.0.0.0:8083)/snippetbox?parseTime=true", "MYSQL data Source")

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

	// Use the session.New() function to initialize a new manager,
	// passing in the  secret key as the parameter. Then we configure it so
	//sessions Always expire after 12 hours

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true //set the secure flag on our session cookie

	//And add the session Manager to our application dependencies

	// Initialize a mysql.SnippetModel instance and add it to the application
	// dependencies.
	// Add add it to the application dependencies
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		session:       session,
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
		users:         &mysql.UserModel{DB: db},
	}

	//initialize a new http.Server struct. We set Addr and Handler fields
	//That the server uses the same newtwork address and routes as before, and
	//the Errorlog fields so That the server now uses the custom errorlog Logger
	// the Event of  any problems

	//Initialize a list.config struct to hold the Most non-fault TLS settings we want
	// the server to use.

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// Set the server's TLSConfig field to USe the tlsConfig variable we just
	//created

	srv := &http.Server{
		Addr:      *addr,
		ErrorLog:  errorLog,
		Handler:   app.routes(),
		TLSConfig: tlsConfig,
		//Add Idle, Read and write timeout
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", *addr)
	// call the listenAndserver() method on our new http.Server Struct.

	//Use the listenAndServeTLS() method to start the https Server. we
	//pass in the paths to the TLS certificate and corresponding private key as
	// the two parameters

	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
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
