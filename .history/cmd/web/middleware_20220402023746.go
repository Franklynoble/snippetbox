package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Franklynoble/snippetbox/pkg/models"
	"github.com/justinas/nosurf"
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

func (app *application) requireAuthentication(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//if the user is not authenticated, redirect them to the login
		//page and return from the middleware chain so that no subject handlers in
		//

		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		//otherwise set the "Cache-Control: no store" header so that pages
		// require authentication are not stored in the users browser cache(or other intermediary cache)
		w.Header().Add("Cache-Control", "no-store")

		//And call the next handler in the chain
		next.ServeHTTP(w, r)
	})

}

// create a new NoSurf middleware function which uses a customize CRF cookie with
// the secure, Oath and  HttpOnly flags set.

func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})
	return csrfHandler
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Check if a authenticatedUserID value exists in the session. if this isn't
		// present* then call the next handler in the chain as normal
		exists := app.session.Exists(r, "authenticatedUserID")
		if !exists {
			next.ServeHTTP(w, r)
			return
		}
		/* Fetch the details of the current user from the database. if no matching
		record is found, or the current user is has been deactivated, remove the
		(invalid) authenticatedID value from their session and call the next
		handler in the chain as normal
		 **/
		user, err := app.users.Get(app.session.GetInt(r, "authenticatedUserID"))
		if errors.Is(err, models.ErrNoRecord) || !user.Active {
			app.session.Remove(r, "authenticatedUserID")
			next.ServeHTTP(w, r)
			return
		} else if err != nil {
			app.serverError(w, err)

		}
	})
}
