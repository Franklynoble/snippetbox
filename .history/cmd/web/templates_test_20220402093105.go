package main

import (
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {

	//initialize a new time.Time object and  pass it to the humanDate function
	// here  the time is in yyyy-mm-dd hh:mm:ss + nsec nanoseconds
	//tm := time.Date(2020, 12, 17, 10, 0, 0, 0, time.UTC)
	//hd := humanDate(tm)

	//check that the output from the humanDate function is in the format we
	// expect. if it is not what we expect, use the t.Error() function to
	// indicate that the test has failed and log the expected and actual values
	//if hd != "17 Dec 2020 at 10:00" {
	//	t.Errorf("want %q; got %q", "17 Dec 2020 at 10:00", hd)
	//}

	//create a slice of anonymous struct containing the test case name,
	//input to our  humanDate() function (the tm field), and the expected output
	//(they want field)
	test := []struct {
		name string
		tm   time.Time
		want string
	}{
		name: "UTC",
		tm:   time.Date(2020, 12, 17, 10, 0, 0, 0, time.UTC),
	}

}
