package forms

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

/* Create a custom Form struct, which anonymously embeds a url.values object
(to hold the Form data) and  an Erros field to hold any validation errors
for the form data
**/

type Form struct {
	url.Values
	Errors errors
}

/*Define a new function to initialize a custom Form struct. notice that
this takes the form data as  the parameter?
**/

func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// implement a Required method to check that specific fields in the form
//data are present and not blank. if any field fails this check,
// add the approppriate message form errors
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

//Implement a MaxLength nethod to check that a specific field in the form
//contains a maximum number of characters. if the check fails then add the
//appropriate message to the form errors.

func (f *Form) MaxLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) > d {
		f.Errors.Add(field, fmt.Sprintf("This field is too long(maximum is %d characters)", d))

	}
}

//implement a PermittedValues method to check that a specific field in the form
//matches one of a set of specific permitted values. if chek fials
// then add the appropriate message to the form errors

func (f *Form) PermittedValues(field string, opts ...string) {
	value := f.Get(field)

	if value == "" {
		return
	}
	for _, opt := range opts {

		if value == opt {
			return
		}

	}
	f.Errors.Add(field, "This field is invalid")
}

// implement a valid method which returns true if there are no errors
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}


/*
  Use the regex.MustCompile() function to parse pattern and compile a regular expression for sanity checking the format of an email address
  This returns *regex.Regexp object, or panics in the event of an error
  Doing this once at runtime, and storing the compile regular expression 
  Object in a variable, is more performant than recompiling the pattern 
   every request
**/
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")