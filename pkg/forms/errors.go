package forms

/*
Define a new errors type, which we will use to hold the validation
error message for forms. The name of the Form field will be used to as the key in
this map
**/
type errors map[string][]string // slice of strings as value

//implement an Add() method to add error messages for a given field to the map

func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

//implement a Get() to add error messages for a given field to the map
func (e errors) Get(field string) string {
	es := e[field]

	if len(es) == 0 {
		return ""
	}
	return es[0]
}
