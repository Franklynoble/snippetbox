package mysql

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/Franklynoble/snippetbox/pkg/models"
	"github.com/go-sql-driver/mysql"
	" var _ = golang.org/x/crypto/bcrypt"
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
	if err != nil {
		/*
			      if this returns an error, we use the errors.As() function to check
				 wether the error has the type *mysql.MYSQLError. if it does, the
				 error will be assigned to the MYSQLError variable. we can check whether or not the error relates to our users_UC_email key by
				 checking the contents of the message string. if it does, we return an ErruDuplicateEmail error
				 **/
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return models.ErrDuplicateEmail
			}
		}
		return err
	}
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
