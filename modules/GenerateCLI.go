package storagemanagementcli

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"code.cloudfoundry.org/bytefmt"
	"github.com/olekukonko/tablewriter"
)

// CLIConfig contains all the configurations needed to run
type CLIConfig struct {
	Selection   Files
	Files       *FileMap
	SortedFiles *Files
	ResultCount *int
	Status      *string
	Active      *bool
	Ignore      *string
}

var platform = runtime.GOOS

// GenerateCLI prints a console log interface of the top files passed in
func GenerateCLI(config CLIConfig) {
	// Configure CLI table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Bytes", "Path"})
	table.SetRowLine(true)

	// Build CLI table
	for _, v := range config.Selection {
		// Convert byte size to human readable sizes
		byteSize := bytefmt.ByteSize(uint64(v.Size))
		table.Append([]string{fmt.Sprint(v.ID), fmt.Sprint(byteSize), v.Path})
	}

	// Print out CLI table
	table.Render() // Send output

	// Print CLI Status
	if len(*config.Status) > 0 {
		fmt.Println(*config.Status)
	}

	// Request input
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Commands: cd <PATH> | delete <ID> | open <ID> | more <NUMBER> | exit")
	t, _ := reader.ReadString('\n')
	te := strings.Split(t, "\n")
	text := te[0]
	entry := strings.Split(text, " ")
	command := entry[0]

	// Exit
	if command == "exit" {
		*config.Active = false
	}

	if len(entry) >= 2 {
		param := entry[1]
		// If command "cd" use the param as a string and Change directories
		if command == "cd" {
			fmt.Println(fmt.Sprintf("Calculating all files in %v...", param))
			*config.Status = ""
			changeDirectory(param, config.SortedFiles, config.Ignore)
		} else {
			// Otherwise parse the param as an int and Set the selected file to that int
			id, err := strconv.Atoi(param)
			checkErr(err)
			selectedFile := (*config.Files)[id]

			// Delete file
			if command == "delete" {
				deleteFile(selectedFile.Path, selectedFile.ID, config.Files, config.SortedFiles)
				paths := strings.Split(selectedFile.Path, "/")
				fileName := paths[len(paths)-1]
				*config.Status = fmt.Sprintf("%v moved to the Trash", fileName)
			}

			// Open file
			if command == "open" {
				*config.Status = fmt.Sprintf("Opening %v", selectedFile.Path)
				openFile(selectedFile.Path)
			}

			// Show more results
			if command == "more" {
				*config.ResultCount = id
				*config.Status = fmt.Sprintf("Now showing the top %v largest paths", id)
			}
		}
	} else {
		*config.Status = "Invalid Command"
	}
}
