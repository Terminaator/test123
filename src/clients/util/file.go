package util

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type File struct {
	file    *string
	Clients []string
}

func (f *File) Read() {
	jsonFile, err := os.Open(*f.file)

	if err != nil {
		log.Fatal(err)
	}

	defer jsonFile.Close()

	bytes, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(bytes, f)
}

func NewFile(file *string) *File {
	f := &File{file: file}
	f.Read()
	return f
}
