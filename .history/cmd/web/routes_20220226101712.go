package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

/*
update the signature for the routes method so that it returns a
 http.Handler instead of *http.serverMux.
**/
func (app *application) routes() http.Handler {

	/*
		       Create a middleware chain containing our 'standard' middleware
			  which will be used for every request application recieves
			**/
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	mux := pat.New()
	mux.Get("/", http.HandlerFunc(app.home))
	mux.Get("/snippet/create", http.HandlerFunc(app.createSnippetForm))
	mux.Post("/snippet/create", http.HandlerFunc(app.createSnippet))
	mux.Get("/snippet/:id", http.HandlerFunc(app.showSnippet))

	//create a file server which serves files out of the "./ui/static" directory
	// Note that the Path given to the http.Dir function is relative to the Object
	//Directory root

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	//use the mux.Handle() function to register the files server as the handler for
	// all URL  paths that start with "/static/." for matching paths, we strip the
	// "/static" prefix before the request reaches the file Server

	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	/*
		Pass the servermux as the 'next' parameter to the scureHeaders middleware
		Because secureHeaders is just a function, and the function returns a http.Handler we don't to do anything else.
		  **/

	// Wrap the existing chain with the logRequest middleware
	return standardMiddleware.Then(mux)
}
