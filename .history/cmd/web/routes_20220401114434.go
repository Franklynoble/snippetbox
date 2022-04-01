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

	/*
		Create a new middleware chain containing the middleware specific to our
		dynamic application routes, For now, this chain will Only contain
		 the session middleware but we will add More to it latter
		**/
	dynamicMiddleWare := alice.New(app.session.Enable)

	mux := pat.New()
	// update these routes to use the new dynamic middleware chain followed by the appropriate handler function

	mux.Get("/", dynamicMiddleWare.ThenFunc(app.home))

	//Add the requireAuthentication middleware to the chain
	mux.Get("/snippet/create", dynamicMiddleWare.Append(app.requireAuthentication).ThenFunc(app.createSnippetForm))
	//Add the requireAuthentication middleware to the chain
	mux.Post("/snippet/create", dynamicMiddleWare.Append(app.requireAuthentication).ThenFunc(app.createSnippet))
	mux.Get("/snippet/:id", dynamicMiddleWare.ThenFunc(app.showSnippet))

	//add all the five  routes  from handlers
	mux.Get("/user/signup", dynamicMiddleWare.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleWare.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleWare.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleWare.ThenFunc(app.loginUser))
	//Add the requireAuthentication middleware to the chain
	mux.Post("/user/logout", dynamicMiddleWare.ThenFunc(app.logoutUser))
	//create a file server which serves files out of the "./ui/static" directory
	// Note that the Path given to the http.Dir function is relative to the Object
	//Directory root

	// leave the static files unchanged, as we do not need static file for our dynamic data
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
