package mysql

import (
	"database/sql"
	"io/ioutil"
	"testing"
)

func newTestDB(t *testing.T) (*sql.DB, func()) {
	//Establish a sql.DB connection pool for our test database. because our
	//setup and the teardown scripts contains mulitiple SQL statement, we need to use the 'multiStatements=true' parameter in our DSN. This instructs
	//our MYSQL database driver to support executing multiple SQL statements
	//							tcp(0.0.0.0:8083)
	db, err := sql.Open("mysql", "root:pass@tcp(127.0.0.1:3036)/test_snippetbox?parseTime=true&multiStatements=true")

	if err != nil {
		t.Fatal(err)

	}

	//Read the setup SQL script from file and execute the statements
	script, err := ioutil.ReadFile("./testdata/setup.sql")

	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
	}

	//Return the connection pool and an anonymous function which reads and
	//executes the teardown script, and closes the connection pool. we can
	//assign this anonymous function and call it later once test has
	//completed
	return db, func() {
		script, err := ioutil.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}
		_, err = db.Exec(string(script))
		if err != nil {
			t.Fatal(err)
		}
		db.Close()
	}

}
