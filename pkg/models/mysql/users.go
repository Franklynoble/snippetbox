package mysql

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/Franklynoble/snippetbox/pkg/models"
	"github.com/go-sql-driver/mysql"
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
	/*
		Retrieve the id and hashed password associated with the given email. if no
		match email exists, or the user is not active, we return the ErrInvalidCredential error
		**/
	var id int
	var hashePassword []byte

	stmt := "SELECT id, hashed_password FROM users WHERE email= ? AND active = TRUE"

	row := m.DB.QueryRow(stmt, email)
	err := row.Scan(&id, &hashePassword) // this would scan and put info to this variable created from above at line 61,62
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrinvalidCredentials
		} else {
			return 0, err
		}
	}
	//Check whether the hashed password and plain-text password provided match
	//if they do not, we return the ErrInvalidCredentials error.
	err = bcrypt.CompareHashAndPassword(hashePassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrinvalidCredentials
		} else {
			return 0, err
		}
	}
	// otherwise, the password is correct. Return the user ID.
	return id, nil

}

//we will use  the Get method to fetch details for a specific user based on their user ID
func (m *UserModel) Get(id int) (*models.User, error) {

	u := &models.User{}

	stmt := `SELECT id, name, email, created, active FROM users WHERE id = ?`

	err := m.DB.QueryRow(stmt, id).Scan(&u.ID, &u.Name, &u.Email, &u.Created, &u.Active)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return u, nil

}

func (m *UserModel) ChangePassword(id int, currentpassword, newpassword string) error {
	var currentHashedPassword []byte

	row := m.DB.QueryRow(`SELECT hashed_password FROM users WHERE id = ?`, id)
	err := row.Scan(&currentHashedPassword)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword(currentHashedPassword, []byte(currentpassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return models.ErrinvalidCredentials
		} else {
			return err
		}
	}
	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newpassword), 12)
	if err != nil {
		return err
	}
	stmt := `UPDATE users SET hashed_password = ? WHERE id = ?`
	_, err = m.DB.Exec(stmt, string(newHashedPassword), id)
	return err

}
