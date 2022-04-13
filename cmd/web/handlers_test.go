package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {
	//Initializing a  new httptest.ResponseRecorder.
	rr := httptest.NewRecorder()

	//initialize a new dummy http.Request.

	r, err := http.NewRequest(http.MethodGet, "/", nil)

	if err != nil {
		t.Fatal(err)
	}
	//call the  ping handler finction , passing in the httptest.ResponseRecorder and
	// the  http0test.ResponseRecorder and the  http.request

	ping(rr, r)
	//call the result() method on the htt.ResponseRecorder to get the
	//http.response generated by the ping handler
	rs := rr.Result()

	// we can then examine the http.response to check that the status code
	// written by the  ping handler was 200.
	if rs.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
	}
	// we can check that the response body written by the ping handler
	// equals status "Ok"
	defer rs.Body.Close()
	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "Ok" {
		t.Errorf("want body to equal %q", "OK")
	}
}
