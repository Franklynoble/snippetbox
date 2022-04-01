package main

import (
	"errors"
	"fmt"

	//"text/template"

	"net/http"
	"strconv"

	"github.com/Franklynoble/snippetbox/pkg/forms"
	"github.com/Franklynoble/snippetbox/pkg/models"
)

//  Define a  home handler function which writes a  byte slice containing
//Hello from SnippetBox as the response body

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	//Because Pat matches the "/" path exactly, we can remove the manual check
	// of r.URL.Path != "/" from this handler
	/*
		if r.URL.Path != "/" {
			app.notFound(w) // use the notfound()  helper
			return
		}
		**/
	//panic("oops! something went wrong") //delibrate panic
	s, err := app.snippets.Lattest()

	if err != nil {
		app.serverError(w, err)
		return
	}
	//
	// for _, snippet := range s {
	// fmt.Fprintf(w, "%v\n", snippet)
	// }
	//
	//use the new render helper func
	app.render(w, r, "home.page.tmpl", &templateData{
		Snippets: s, // show all snippets
	})
}

/*
	// Create an instance of a templateData struct holding the slice of
	//Snippets.
	data := &templateData{Snippets: s}

	//initialize a slice containing the path to the two files. Note that the

	// home.gae.tmpl file must be  the *first* file in the slice

		files := []string{
			"./ui/html/home.page.tmpl",
			"./ui/html/base.layout.tmpl",
			"./ui/html/footer.partial.tmpl",
		}

		//use the template.ParseFiles() to read the files and store the template in a template set. Notice we can pass the slice of files paths
		//as a variadic parameter
		tp, err := template.ParseFiles(files...)

		if err != nil {
			app.errorLog.Println(err.Error())
			app.serverError(w, err)
			return
		}

		//pass in the Template Struct why executing  the Template
		err = tp.Execute(w, data)

		if err != nil {
			app.errorLog.Println(err.Error())
			app.serverError(w, err)

		}

		//w.Write([]byte("Hello from Snippetbox"))

	 **/

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	// Pat does not strip the colon from the named capture key, so we need to
	//get the value of ":id" from the query string instead of the "id"

	id, err := strconv.Atoi(r.URL.Query().Get(":id"))

	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	// Use the SnippetModel object's Get Method to retrieve the Data for a specific record based on its ID. if no matching record is found,
	// return a 404 NOT Found response.
	// this returns  Models. and  err
	s, err := app.snippets.Get(id)

	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	/*
		    Use the PopString() method to retrieve the value for the "flash" key.
		    PopString() also deletes the key and value from the session data, so it acts
		     like a one-time fetch. if there is no matching key in the session
			data this will return the empty string

		   **/
	//	flash := app.session.PopString(r, "flash")

	// pass  the flash message  to the template
	// now use the new render helper.
	app.render(w, r, "show.page.tmpl", &templateData{
		Snippet: s, // show single snippet
		//	Flash:   flash,
	})
}

/*
		     initialize a slice containing the print to the show.page.tmpl file
			plus the base layout and footer partial that we made earlier

	data := &templateData{Snippet: s}

	files := []string{
		"./ui/html/show.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}
	//parse the template file
	ts, err := template.ParseFiles(files...)

	if err != nil {
		app.serverError(w, err)
		return
	}
	 And then execute them. Notice how we are passing in the snippet
	data (a models.Snippet struct) as the final parameter

	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err)

	}
	// write the snippet data as a plain-text HTTP response body.
	//	fmt.Fprintf(w, "Displaying a specific snippet with ID %v...", s)
	//w.Write([]byte("Display a specific snippet ..."))
	 **/

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {

	app.render(w, r, "create.page.tmpl", &templateData{
		//pass empty forms.Form object to the template
		Form: forms.New(nil),
	})
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {

	// User r.method to check whether the method is  POST or  not. Note that
	//http.MethodPost is a constant equal to the string  "POST"

	//checking if the request is a POST is now superflous and can be removed
	/*
		if r.Method != http.MethodPost {
			// if it's not, use the w.WriteHeader() method to send a  405 status
			// code, and  the w.writerHead() method to send "Method Not allowed"
			//response Body. we Then return from the function so that the
			// sunsequent code is not executed.
			//w.Header().Set("Allow", http.MethodPost)
			//w.WriteHeader(405)
			w.Header().Set("Allow", http.MethodPost)
			//use the http.Error() function  to see a 405 a status code and "Method Not allowed string as the Method b
			app.clientError(w, http.StatusMethodNotAllowed) // Use the clientError() from the helper
			return
		}
		**/ /*
		title := `O snail`
		content := `O snail \n Climb Mount Fuji, \nBut slowly, slowly!\n Kobayshit`
		expires := `7`
		id, err := app.snippets.Insert(title, content, expires)

		if err != nil {
			app.errorLog.Print(err)
			app.serverError(w, err)
			return}
		**/
	// First we call r.ParseForm() which adds any  data to the POST request bodies
	// to the r.PostForm map. this also works in the same way for PUT and PATCH
	//requests if there are any errors, we use our app.ClientError helper to send
	// a 400 Bad Request response to the user.
	err := r.ParseForm()

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}
	// Use the r.PostForm.Get() method to retrieve relevent data fields
	//from the r.PostForm map

	//recent errors handling
	// create a new forms.Form struct containing the POSTed data from the
	//form, then use the validation methods to check the content

	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	//if the form is not valid, redisplay the template passing in the
	//form.Form object as the data.

	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
		return
	}

	/*Because the form (with type url.Values) has been anonymously embedded
	in the form.Form struct, we can use the Get() method to retrieve
	// the validated value for a particular form field
	**/
	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))

	if err != nil {
		app.serverError(w, err)
		return
	}
	/*Use the Put() method to add a string value ("Your snippet was saved
	successfully!") and the corresponding key ("flash") to the session
	//data. Note that there is no existing session for the current user
	or their session has expired) then a new , empty session for them
	will be automatically be created by the session middleware

	**/
	app.session.Put(r, "flash", "snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)

}

/*
	errors := make(map[string]string)
//
	//check that the title is not blank and it is not morethan  100 characters
	if strings.TrimSpace(title) == "" {
		errors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		errors["title"] = "This is too long(maximum is 100 characters)"
	}

	//check that the content field is not blank
	if strings.TrimSpace(content) == "" {
		errors["content"] = "This field can not be Blank"
	}
	//check the expires field is not blank and  matches one of the permitted
	//values("1", "7", or "365")

	if strings.TrimSpace(expires) == "" {
		errors["expires"] = "this field cannot be blank"
	} else if expires != "365" && expires != "7" && expires != "1" {
		errors["expires"] = "This field is invalid"
	}

	//if there any errors, dum it in the plane text and return
	// if tdhere are any validation errors, re-display the create.page.tmpl
	//template passing in the validation errors and previously submitted
	//r.PostForm
	if len(errors) > 0 {
		app.render(w, r, "create.page.tmpl", &templateData{
			FormErrors: errors,
			FormData:   r.PostForm,
		})

		return
	}
	// Create a new snippet record in the  database using the form data
	// Change the redirect to use the new semantic URL of /snippet/:id
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	//w.Write([]byte("Create a  new snippet.."))
}
**/

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {

	app.render(w, r, "signup.page.tmpl", &templateData{
		Form: forms.New(nil),
	})

	//	fmt.Fprintf(w, "Display the user signup form...")
}
func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {

	// pass the Form
	err := r.ParseForm()

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	//validate the Form contents using the Form helper we made earlier
	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.MaxLength("name", 255)
	form.MaxLength("email", 255)
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 10)

	//if there are any errors redisplay the signUp Form
	if !form.Valid() {
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	}
	// Try to create a new user record in the database. if the email already exist
	//add an error message to the form and re-display it

	err = app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))

	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.Errors.Add("email", "addres is already in use")
			app.render(w, r, "signup.page.tmpl", &templateData{Form: form})

		} else {
			app.serverError(w, err)
		}
		return
	}
	//otherwise add a confirmation flash message to the session confirming that
	//their signUp worked and asking them to log in.
	app.session.Put(r, "flash", "your sign up was successful. Please log in.")

	// add redirect to the user to the login page

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {

	app.render(w, r, "login.page.tmpl", &templateData{
		Form: forms.New(nil),
	})

	//fmt.Fprintf(w, "Display the user login form ...")

}
func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Check whether the credentials are valid. if they are not, add a generic error
	//message to the Form failures map and redisplay the Login
	// This field is only available after ParseForm is called
	//The HTTP client ignores PostForm and uses Body instead.

	forms := forms.New(r.PostForm) // PostForm contains the parsed form data from PATCH, POST	or PUT body parameters.

	id, err := app.users.Authenticate(form.get("email"), form.get("password"))
	if err != nil {
		if errors.Is(err, models.ErrinvalidCredentials) {
			form.Errors.Add("generic", "email of pass is incorrect")
			app.render(w, r, "login.page.tmpl", &templateData{Form: form})
		} else {
			app.serverError(w, err)
		}
		return
	}

	//Add the ID of the current user to the session, so that they are now 'logged' in
	app.session.Put(r, "authenticateUserID", id)
}

//fmt.Fprintf(w, "Authenticate and  login the user ...")

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Logout  the User ...")
}
