package mysql

import (
	"database/sql"

	"github.com/Franklynoble/snippetbox/pkg/models"
)

//Define a snippetModel type wraps a sql.DB connection pool.
type SnippetModel struct {
	DB *sql.DB
}

//this will insert a  new snippet into the Database
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {

	//write the SQL  statement we want to execute. i have split it over two line
	//for readability (which is why it's surrounded with backquotes instead
	// of normal double quotes

	stmt := `INSERT INTO snippets (title, content, created, expires)
          VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	//use the Exec() method on the embedded connection pool to execute the
	// statment. The first parameter is the SQL statement, followed by the
	//title, content and sql.Result object, which contains some basic
	//information about what happend when the statement was executed
	result, err := m.DB.Exec(stmt, title, content, expires)

	if err != nil {
		return 0, err
	}

	// Use the LastInsertedID() method on the result object to get the iD of Our
	// newly Inserted record in the snippets  table
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

//this will return a specific snippet based on  its id.
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {

	// write the SQL statement we want to execute, Again, i have split it over two
	// line for readability

	stmt := `SELECT id, title, content, created, expires FROM snippets
	 		WHERE expires > UTC_TIMESTAMP() AND id = ?`

	// Use the QueryRo() method on the connection pool to execute our
	//SQL statement ,  passing in the unstrusted id variable as the value for the
	//place holder parameter. this returns to a sql.Row Object which
	//holds the result from the Database.
	row := m.DB.QueryRow(stmt, id)

	// initialize a pointer to new zeroed Snippet struct.

	s := &models.Snippet{}

	//Use row.Scan() to copy the values from each field in sql.Row to the
	// corresponding field in the snippets struct. Notice that the arguments
	// to row.Scan are *pointers* to the place you want to copy the data into,
	// and the number of  argument must be exactly thesame as the  number  of  columns returned by your statement. if  the query returns no rows, then
	//rows.Scan() will return a sql.ErrNoRows error. we check for that  and return
	//Our Models.ErroNoRecord error instead of a snippet object

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	// if everything went Ok then return the Snippet Object.
	return s, nil
}

//this will return the 10 most recently created snippets.
func (m *SnippetModel) Lattest() ([]*models.Snippet, error) {

	stmt := `SELECT id, title, content, created, expires FROM snippets
	          WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)

	if err != nil {

		return nil, err
	}
	defer rows.Close()

	snippets := []*models.Snippet{}

	for rows.Next() {

		s := &models.Snippet{}

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	return snippets, nil

}
