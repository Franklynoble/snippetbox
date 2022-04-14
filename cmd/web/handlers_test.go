package main

import (
	"bytes"
	"net/http"
	"testing"
)

func TestPing(t *testing.T) {
	//update testPing to use

	// create a new instance of our application struct. for now, this just
	//contains a couple of mock loggers(which discard anything written to them)
	// app := &application{
	// errorLog: log.New(io.Discard, "", 0),
	// infoLog:  log.New(ioutil.Discard, "", 0),
	// }

	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()
	// We then use the httptest.NewTLSServer() function to create a new  test
	//server, passing in the  value returned by our app.routes() method as the  handler
	// for the server. this starts up a HTTPS server which listens on new
	//randomly-chosen port of your local machine for the duraation of the test
	// notice that we defer a call to ts.close() to shutdown the server when the test finishes
	// ts := httptest.NewTLSServer(app.routes())
	// defer ts.Close()

	// the network address that the test server is listening on is contained
	// in the ts.URL field. we can use  this along with the ts.Client().Get()
	// method to make a GET/ ping request against the test server.This returns a http.response struct containing the respponse
	code, _, body := ts.get(t, "/ping")
	// if err != nil {
	// t.Fatal(err)
	// }
	//  we can then check the value of the response of the status code and body using
	if code != http.StatusOK {
		t.Errorf("wants %d; got %d", http.StatusOK, code)
	}
	// defer rs.Body.Close()
	// body, err := ioutil.ReadAll(rs.Body)
	// if err != nil {
	// t.Fatal(err)
	// }
	if string(body) != "Ok" {
		t.Errorf("wants body equal %q", "Ok")
	}
	//initializing a new
	//Initializing a  new httptest.ResponseRecorder.
	//
	// rr := httptest.NewRecorder()
	//
	//initialize a new dummy http.Request.
	//
	// r, err := http.NewRequest(http.MethodGet, "/", nil)
	//
	// if err != nil {
	// t.Fatal(err)
	// }
	//call the  ping handler finction , passing in the httptest.ResponseRecorder and
	// the  httptest.ResponseRecorder and the  http.request

	//	ping(rr, r)
	//call the result() method on the htt.ResponseRecorder to get the
	//http.response generated by the ping handler
	// rs = rr.Result()
	//
	// we can then examine the http.response to check that the status code
	// written by the  ping handler was 200.
	// if rs.StatusCode != http.StatusOK {
	// t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
	// }
	// we can check that the response body written by the ping handler
	// equals status "Ok"
	// defer rs.Body.Close()
	// body, err = ioutil.ReadAll(rs.Body)
	// if err != nil {
	// t.Fatal(err)
	// }
	// if string(body) != "Ok" {
	// t.Errorf("want body to equal %q", "OK")
	// }
}

func TestShowSnippet(t *testing.T) {
	// create a new instance of our application struct which uses mocked
	//dependencies

	app := newTestApplication(t)
	//establish a new test server for runing end-to-end tests.
	ts := newTestServer(t, app.routes())

	defer ts.Close()
	// Set up some table-driven tests to ckeck the response sent  by our
	//application for different
	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody []byte
	}{
		{"Valid ID", "/snippet/1", http.StatusOK, []byte("An old silent pond...")},
		{"Non-existent ID", "/snippet/2", http.StatusNotFound, nil},
		{"Negative ID", "/snippet/-1", http.StatusNotFound, nil},
		{"Decimal ID", "/snippet/1.23", http.StatusNotFound, nil},
		{"String ID", "/snippet/foo", http.StatusNotFound, nil},
		{"Empty ID", "/snippet/", http.StatusNotFound, nil},
		{"Trailing slash", "/snippet/1/", http.StatusNotFound, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.urlPath)

			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}
			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body to contain %q", tt.wantBody)
			}
		})
	}
}
