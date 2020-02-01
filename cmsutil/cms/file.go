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
	Folder string
	Filename string
}

func (f *File) WriteFile(path string) {
	fmt.Println("cms package: File: "+path+f.Folder+"/"+f.Filename)
}
