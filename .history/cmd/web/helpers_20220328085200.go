package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

//The serverError helper writes an error message and stack trace to the errorLogger
//then sends a generic 500 Internal Server Error response to the user.

func (app *application) serverError(w http.ResponseWriter, err error) {

	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())

	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

}

// the ClientError helper sends a specific Status code and Corresponding description
//to the User. We'll use this later in the Book to send responses like 4000 "Bad Request
// when there is a problem with  the request that  the user sent

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	/*
	   Retrieve the appropriate template set from the cache based on the page
	   (like 'home.page.tmpl'). if no entry exists in the cache with the
	   provided name, call the serverError helper method  that we made earlier

	*/
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("the template %s does not exist", name))
		return
	}
	// Execute the template set,  passing in any dynamic data.
	// err := ts.Execute(w, td)
	// if err != nil {
	// app.serverError(w, err)
	//}
	// Initialize a  new buffer.
	buf := new(bytes.Buffer)

	/*
		write the template to the buffer, instead of straight  to the
		http.responseWriter if there is an error, call our serverError helper
		 and return
		**/
	err := ts.Execute(buf, app.addDefaultData(td, r))

	if err != nil {
		app.serverError(w, err)
		return
	}
	buf.WriteTo(w)
}

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {

	if td == nil {
		td = &templateData{}
	}
	td.CurrentYear = time.Now().Year()
	//Add the flash message to the Template data if one Exists
	td.Flash = app.session.PopString(r, "flash")
	return td
}
