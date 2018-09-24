package storagemanagementcli

import (
	"errors"
	"testing"
)

func TestCheckErr(t *testing.T) {
	fatal := false
	logFatal = func(...interface{}) {
		fatal = true
	}
	checkErr(errors.New("error alert"))
	if fatal != true {
		t.Error("Expected checkErr to cause a fatal error, but it did not")
	}
}
