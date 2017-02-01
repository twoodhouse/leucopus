package ioManager

import "testing"

func TestIoManager(t *testing.T) {
	iom := New([]string{}, []string{"http://buttonlight01.env:5001/checklight"}) //https://golang.org
	// iom := New([]string{}, []string{"https://golang.org"}) //

	iom.GetIncomingDataRowMap()
}
