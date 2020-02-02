/*
Copyright © 2020 John Arroyo

cms File package

write retrieved data to a file
*/

package cms

import (
	"fmt"
)

type File struct {
	Path string
	Folder string
	Filename string
}

func (f *File) WriteFile(data string) {
	// check if folder empty?
	// make sure filepath is not empty

	var filepath = fmt.Sprintf("%v%v/%v", f.Path, f.Folder, f.Filename)
	
	fmt.Println("cms package: WriteFile: "+filepath)
}
