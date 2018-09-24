package storagemanagementcli

import (
	"testing"
)

func TestLimitFiles(t *testing.T) {
	var files = Files{
		File{ID: 1, Size: int64(1), Path: ""},
		File{ID: 2, Size: int64(2), Path: ""},
	}

	t.Log("Limits # of Files to value passed in")
	limit := 1
	count := len(LimitFiles(files, limit))

	if count != 1 {
		t.Errorf("Expected count of 1 file, but it was %d instead.", count)
	}

	t.Log("Should not affect counts of files lower than the limit")
	limit = 3
	count = len(LimitFiles(files, limit))

	if count != 2 {
		t.Errorf("Expected count of 2 file, but it was %d instead.", count)
	}
}
