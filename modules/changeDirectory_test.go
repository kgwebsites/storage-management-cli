package storagemanagementcli

import "testing"

func TestChangeDirectory(t *testing.T) {
	var testParam1, testParam2 string
	var testParam3 FileMap
	var getFilesResponse = FileMap{}
	getFilesResponse[1] = File{ID: 1, Size: 1, Path: ""}
	sortedResponse := Files{getFilesResponse[1]}

	GetFiles = func(param1 string, param2 string) FileMap {
		testParam1 = param1
		testParam2 = param2
		return getFilesResponse
	}
	SortFiles = func(param3 FileMap) Files {
		testParam3 = param3
		return sortedResponse
	}
	sortedFiles := Files{}
	ignore := "ignore"

	changeDirectory("dir", &sortedFiles, &ignore)

	t.Log("ChangeDirectory should pass GetFiles the dir string and ignore string")
	if testParam1 != "dir" {
		t.Errorf("Expected testParam1 to be 'dir', but it was %s instead.", testParam1)
	}
	if testParam2 != "ignore" {
		t.Errorf("Expected testParam2 to be 'ignore', but it was %s instead.", testParam2)
	}

	t.Log("ChangeDirectory should pass SortFiles the response from GetFiles")
	if testParam3[0] != getFilesResponse[0] {
		t.Errorf("Expected testParam3 to be getFilesResponse but it was %v instead.", getFilesResponse)
	}

	t.Log("ChangeDirectory should mutate referenced sortedFiles and ignore variables")
	if sortedFiles[0] != sortedResponse[0] {
		t.Errorf("Expected sortedFiles to be mutated to the sortedResponse, but it was %v instead", sortedResponse)
	}
}
