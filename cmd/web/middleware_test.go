package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecureHeaders(t *testing.T) {
	//initializing a new httptest.ResponseRecorder and dummy.Request

	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)

	if err != nil {
		t.Fatal(err)
	}
	// create a mock HTTP handler that we can pass to our secureHeaders
	//middleware, which writes  a 200 status code and "Ok" reponse body

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Ok"))
	})

	//pass the mock HTTP handler to our secureHeaders middleware.Because
	// secureHeaders * returns a http.Handler we can call its ServeHTTP()
	//method, passing in the http.ResponseRecorder and dummy http.Request to
	//executes
	secureHeaders(next).ServeHTTP(rr, r)
	// call the Result middleware() method on the http.ResponseRecorder to get the results
	// of the test
	rs := rr.Result()

	//check that the middleware has correctly has set the x-Frame-options header
	// on the response
	frameOptions := rs.Header.Get("x-Frame-Options")

	if frameOptions != "deny" {
		t.Errorf("want %q; got %q", "deny", frameOptions)
	}
	//check that the middleware has correctly set the x-xss-Protection header
	//on the response
	xssProtection := rs.Header.Get("X-XSS-Protection")
	if xssProtection != "1; mode=block" {
		t.Errorf("want %q; got %q", "1; mode=block", xssProtection)
	}
	//check that the middleware has correctly called the next handler in line
	//and the response status code and body as expected
	if rs.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
	}
	defer rs.Body.Close()
	body, err := ioutil.ReadAll(rs.Body)

	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "Ok" {
		t.Errorf("want body to equal %q", "Ok")
	}
}
