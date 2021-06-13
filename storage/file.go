/*
Copyright © 2021 John Arroyo

storage File package

write retrieved data to a file
*/

package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// File interactions
type File struct {
	Path     string
	Folder   string
	Filename string
	Verbose  bool
}

// Init initialize File
func (f *File) Init(path string, folder string, filename string, verbose bool) {
	f.Path = path
	f.Folder = folder
	f.Filename = filename
	f.Verbose = verbose
}

func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}

// Check to see if the folder exists, if not create it.
func prepFolder(folder string) {
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		os.MkdirAll(folder, os.ModePerm)
	}
}

// WriteFileJSON write the supplied json data to a file.
// Full filepath comes from the vars in the struct
// Path + Folder + Filename
// Path & Folder can be empty, but filename cannot be
func (f *File) WriteFileJSON(data map[string]interface{}) {
	// @todo check if folder empty?
	// @todo make sure filepath is not empty
	var fullpath = fmt.Sprintf("%v%v", f.Path, f.Folder)
	var filepath = fmt.Sprintf("%v/%v", fullpath, f.Filename)

	fileData, err := json.MarshalIndent(data, "", " ")
	checkErr(err)

	prepFolder(fullpath)

	file, err := os.Create(filepath)
	checkErr(err)

	err = ioutil.WriteFile(filepath, fileData, 0644)
	checkErr(err)

	if f.Verbose {
		log.Println("Write file: " + filepath)
	}

	file.Sync()
	file.Close()
}

// WriteFile write the supplied data string to a file.
// Full filepath comes from the vars in the struct
// Path + Folder + Filename
// Path & Folder can be empty, but filename cannot be
func (f *File) WriteFile(data string) {
	// @todo check if folder empty?
	// @todo make sure filepath is not empty
	var fullpath = fmt.Sprintf("%v%v", f.Path, f.Folder)
	var filepath = fmt.Sprintf("%v/%v", fullpath, f.Filename)
	prepFolder(fullpath)

	file, err := os.Create(filepath)
	checkErr(err)

	d1 := []byte("Hello\nWorld!!!\n")
	err = ioutil.WriteFile(filepath, d1, 0644)
	checkErr(err)

	file.Sync()
	file.Close()
}

// DownloadFile from URL
func (f *File) DownloadFile(url string, filename string) error {
	if f.Verbose {
		log.Println("Download file " + url)
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Make sure folder exists
	folderpath := fmt.Sprintf("%v%v", f.Path, f.Folder)
	prepFolder(folderpath)

	// Create the file
	filepath := fmt.Sprintf("%v/%v", folderpath, filename)
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
