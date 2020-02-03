/*
Copyright © 2020 John Arroyo

cms File package

write retrieved data to a file
*/

package cms

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type File struct {
	Path     string
	Folder   string
	Filename string
}

func (f *File) FileInit(path string, folder string, filename string) {
	f.Path = path
	f.Folder = folder
	f.Filename = filename
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

// Write the supplied json data to a file.
// Full filepath comes from the vars in the struct
// Path + Folder + Filename
// Path & Folder can be empty, but filename cannot be
func (f *File) WriteFileJson(data []interface{}) {
	// @todo check if folder empty?
	// @todo make sure filepath is not empty
	var fullpath = fmt.Sprintf("%v%v", f.Path, f.Folder)
	var filepath = fmt.Sprintf("%v/%v", fullpath, f.Filename)

	fmt.Println("Write file: " + filepath)
	fileData, err := json.MarshalIndent(data, "", " ")
	checkErr(err)

	prepFolder(fullpath)

	file, err := os.Create(filepath)
	checkErr(err)

	err = ioutil.WriteFile(filepath, fileData, 0644)
	checkErr(err)

	file.Sync()
	file.Close()
}

// Write the supplied data string to a file.
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

func (f *File) DownloadFile(url string, filename string) error {
	fmt.Println("Download file " + url)

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
