package main

import (
	"fmt"
	"net/http"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("x-xss-protection", "1; mode=block")
		w.Header().Set("X-Frame-options", "deny")

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		/*Create a deferred function (which will always be run in the event)
		of a panic as go unwinds the stack
		**/

		defer func() {
			// Use a builtin recover function to check if there has been a
			// a panic or not. if there has
			if err := recover(); err != nil {
				//set the connection: close header on the response.
				w.Header().Set("connection", "close")
				//call the app.serverError helper method to return a 500
				// Internal Server response.
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func requireAuthentication(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//if the user is not authenticated, redirect them to the login
		//page and return from the middleware chain so that no subject handlers in
		//
	})

}
