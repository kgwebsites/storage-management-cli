package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"code.cloudfoundry.org/bytefmt"
	"github.com/olekukonko/tablewriter"
)

type file struct {
	size int64
	path string
}

type fileID struct {
	ID   string `json:"ID"`
	Size int64  `json:"Size"`
	Path string `json:"Path"`
}

var files []file
var selection []fileID

var wg sync.WaitGroup
var root = "/Users/KyleGoss/Documents/projects"

func main() {
	// Walk path concurrently recursively and retrieve file meta data
	wg.Add(1)
	walkDir(root)
	wg.Wait()

	// Sort files by size
	sort.Slice(files, func(i, j int) bool {
		return files[i].size > files[j].size
	})

	// Get the biggest 10 largest files
	files = files[0:10]
	for i, f := range files {
		selection = append(selection, fileID{ID: fmt.Sprint(i + 1), Size: f.size, Path: f.path})
	}

	// Create tmp dir/files
	tf, err := createTmp()
	checkErr(err)
	defer tf.Close()

	// Write biggest 10 file meta data to tmp file
	err = writeToTmp(selection, tf)
	checkErr(err)

	// Configure CLI table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Bytes", "Path"})
	table.SetRowLine(true)

	// Build CLI table
	for _, v := range selection {
		// Convert byte size to human readable sizes
		byteSize := bytefmt.ByteSize(uint64(v.Size))
		table.Append([]string{fmt.Sprint(v.ID), fmt.Sprint(byteSize), trimPath(v.Path)})
	}
	table.Render() // Send output

	// Request input
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Commands: delete <ID> | open <ID> | path <ID> | more <NUMBER> | cd <PATH>")
	text, _ := reader.ReadString('\n')
	entry := strings.Split(text, " ")
	command := entry[0]
	param := entry[1]
	if command == "delete" {
		deleteFile(param)
	}
	// Delete file by id
	// Open file by id
	// Show full path by id
	// Show more results
	// Change directories

}

func walkDir(dir string) {
	defer wg.Done()

	visit := func(path string, f os.FileInfo, err error) error {
		files = append(files, file{path: path, size: f.Size()})
		if f.IsDir() && path != dir {
			wg.Add(1)
			go walkDir(path)
			return filepath.SkipDir
		}
		return nil
	}

	filepath.Walk(dir, visit)
}

func trimPath(path string) string {
	path = strings.Trim(path, "\r")
	sp := strings.Split(path, "/")
	rs := strings.Split(root, "/")
	rootBase := rs[len(rs)-1]
	if strings.Split(root, "")[len(root)-1] == "/" {
		rootBase = rs[len(rs)-2]
	}
	if len(sp) >= 2 {
		return strings.Join([]string{rootBase, "...", sp[len(sp)-1]}, "/")
	}
	return path
}

func createTmp() (*os.File, error) {
	err := os.MkdirAll(".tmp", 0777)
	if err != nil {
		return nil, err
	}
	// Make temp file
	tf, err := os.Create(".tmp/cache")
	if err != nil {
		return nil, err
	}
	return tf, nil
}

func writeToTmp(files []fileID, tf *os.File) error {
	data, err := json.Marshal(files)
	if err != nil {
		return err
	}

	// Write file data to temp file
	_, err = tf.Write(data)
	if err != nil {
		return err
	}

	tf.Sync()
	return nil
}

func deleteFile(id string) {
	var cached []fileID
	data, err := ioutil.ReadFile(".tmp/cache")
	checkErr(err)
	err = json.Unmarshal(data, &cached)
	checkErr(err)
	for _, file := range cached {
		if strings.TrimSpace(file.ID) == strings.TrimSpace(id) {
			paths := strings.Split(file.Path, "/")
			fileName := paths[len(paths)-1]
			// Mac OS
			// Get user meta data
			usr, err := user.Current()
			checkErr(err)
			// Place the file in the trash
			err = os.Rename(file.Path, fmt.Sprintf("%v/.Trash/%v", usr.HomeDir, fileName))
			checkErr(err)
			// Rewrite the cache

			// Reprint the table
			fmt.Println(fileName, "moved to the Trash")
		}
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
