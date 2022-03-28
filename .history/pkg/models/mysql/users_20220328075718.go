package mysql

import (
	"database/sql"

	"github.com/Franklynoble/snippetbox/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

// we will use the Insert method to add new record to the users table
func (m *UserModel) Insert(name, email, password string) error {
	// Create a bcrypt hash of the plain-text password

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stmt := `INSERT INTO users(name, email, hashed_password,created)
            VALUES(?,?,?,UTC_TIMESTAMP())`

	//Use the exec() method to insert the user details and hashed password
	// into the users table

	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))

	return nil

}

/*
we will use  the Authenticate method to  verify whether a user exist with
the provided email address and password. this will return the relevant

**/
func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

//we will use  the Get method to fetch details for a specific user based on their user ID
func (m *UserModel) Get(id int) (*models.User, error) {

	return nil, nil
}
