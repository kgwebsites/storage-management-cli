package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
	"os"
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

var files []file

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

	// Create tmp dir/files
	tf, err := createTmp()
	logErr(err)
	defer tf.Close()

	// Write biggest 10 file meta data to tmp file
	err = writeToTmp(files, tf)
	logErr(err)

	// Configure CLI table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Bytes", "Path"})
	table.SetRowLine(true)

	// Build CLI table
	for i, v := range files {
		// Convert byte size to human readable sizes
		byteSize := bytefmt.ByteSize(uint64(v.size))
		table.Append([]string{fmt.Sprint(i + 1), fmt.Sprint(byteSize), trimPath(v.path)})
	}
	table.Render() // Send output

	// Request input
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Commands: delete <ID> | open <ID> | path <ID> | more <NUMBER> | cd <PATH>")
	text, _ := reader.ReadString('\n')
	entry := strings.Split(text, " ")
	command := entry[0]
	param := entry[1]
	fmt.Println(command + " " + param)
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

func writeToTmp(files []file, tf *os.File) error {
	for i, f := range files {
		dataID := make([]byte, 4)
		binary.LittleEndian.PutUint32(dataID, uint32(i))
		dataSize := make([]byte, 4)
		binary.LittleEndian.PutUint32(dataSize, uint32(f.size))
		dataPath := []byte(f.path)
		data := append(dataID[:], dataSize[:]...)
		data = append(data[:], dataPath[:]...)

		// Write file data to temp file
		_, err := tf.Write(data)
		if err != nil {
			return err
		}

		tf.Sync()
	}
	return nil
}

func logErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
