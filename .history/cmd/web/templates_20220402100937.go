package main

import (
	"html/template"
	"net/url"
	"path/filepath"
	"time"

	"github.com/Franklynoble/snippetbox/pkg/forms"
	"github.com/Franklynoble/snippetbox/pkg/models"
)

/*
Define a  templateData type to act as the hlding structure for
any dynamic data that we want to pass to our HTML template

at the Moment, it only contains one field, but we'll add more to its as the build progresses

**/
//Add FormData and FormErrors field to the templateData struct

//Add a new IsAuthenticated
type templateData struct {
	CSRFToken       string
	CurrentYear     int
	Flash           string
	FormData        url.Values
	Form            *forms.Form
	IsAuthenticated bool
	FormErrors      map[string]string
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
}

/*Create a humanDate which returns
representation of a time.Time object
**/
func humanDate(t time.Time) string {
	//Return Empty string if time has zero value
	if t.IsZero() {
		return ""
	}
	// Convert the time to UTC before formatting it
	return t.UTC().Format("02 Jan 2006 at 15:04")

	return t.Format("02 Jan 2006 at 15:04")
}

/*
Initialize a template.FuncMap object and store it in global variable. this is
essentially a string-keyed map which acts as a lookup between the names of our
custom template functions and the functions themselves
**/
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {

	// Initialize a new map to act the cache
	cache := map[string]*template.Template{}

	//use the filepath.Glob function to get a slice of all filepaths with
	// the extension '.page.tmpl' This essentially gives us a slice of all the
	// 'page' templates for the application

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))

	if err != nil {
		return nil, err
	}
	// loop through the pages one-by-one
	for _, page := range pages {
		//extract the file name (like "home.page.tmpl") from the full file path
		//and assign it to the name variable.

		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)

		//parse the page template file in to a template set
		//ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		/*
			Use the ParseGlob method to add any 'layout' template to the
			template set (in our case, it's just the 'base' layout at the  moment)		**/
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))

		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}
		//Add the template set to the cache, using the name of the page
		//(like 'home.page.tmpl')
		cache[name] = ts
	}
	return cache, nil

}
