package models

import (
	"errors"
	"time"
)

//the top-
//level data types that our database model will use and return

var (
	ErrNoRecord = errors.New("models: no matching records found")
	//Add a new ErrInvalidacredentials error. we will use this later if a user
	//tries to login with a n incorrect email address or password
	ErrinvalidCredentials = errors.New("models: invalid credentials")
	//Add a new ErrDuplicateEmail error. we will use this later if a user
	//tries to sign up with an email that's already in use.
	ErrDuplicateEmail = errors.New("models: duplicate email")
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

//Define a  new User type. Notice how the field names types align
//with the column in the databases 'users table'

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
	Active         bool
}
