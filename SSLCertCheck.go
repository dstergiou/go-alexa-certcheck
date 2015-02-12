package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// ReadFile reads a file from a URL
// No checks at this point - generic errors will be returned
func ReadFile(_url string) (_bytes []byte, _err error) {
	fmt.Printf("Reading file from: %s \n", _url)
	var res *http.Response
	res, _err = http.Get(_url)
	if _err != nil {
		log.Fatal(_err)
	}

	_bytes, _err = ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if _err != nil {
		log.Fatal(_err)
	}

	//fmt.Printf("ReadFile: %s", string(_bytes))
	fmt.Printf("WriteFile: Size of download: %d\n", len(_bytes))
	return
}

// WriteFile writes a file to the disk
// No error checks at this point
func WriteFile(_target string, _bytes []byte) (_err error) {
	fmt.Printf("WriteFile: Size of download: %d\n", len(_bytes))
	if _err = ioutil.WriteFile(_target, _bytes, 0444); _err != nil {
		log.Fatal(_err)
	}
	return
}

// DownloadToFile access a URL, a target filename and a final filename
// Downloads the file from the URL and writes it to the _name target
// No decent error checking at this point
func DownloadToFile(_url string, _target string, _name string) {
	fmt.Printf("DownloadToFile from: %s\n", _url)
	if bytes, err := ReadFile(_url); err == nil {
		fmt.Printf("%s is now downloaded\n", _name)
		if WriteFile(_target, bytes) == nil {
			fmt.Printf("%s is now copied: %s\n", _name, _target)
		}
	}
}

// UnzipFile unzips a file
func UnzipFile(_name string) (unzipped string) {
	zipfile := _name
	fmt.Println("Opening to unzip: ", zipfile)

	reader, err := zip.OpenReader(zipfile)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	for _, f := range reader.Reader.File {
		unzipped = f.Name
		zipped, err := f.Open()
		if err != nil {
			log.Fatal(err)
		}
		defer zipped.Close()

		path := filepath.Join("./", f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
			fmt.Println("Creating directory ", path)
		} else {
			writer, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, f.Mode())
			if err != nil {
				log.Fatal(err)
			}
			defer writer.Close()
			if _, err = io.Copy(writer, zipped); err != nil {
				log.Fatal(err)
			}
			fmt.Println("Decompressing: ", path)
		}
	}
	return unzipped
}

// Main works as follows:
// 1. Delete the target filename, before any download takes place
// 2. Download the file using DownloadToFile
// 3. Unzips the Alexa file using UnzipFile
func main() {

	// Variable assignment
	var url = os.Args[1]
	var file = os.Args[2]
	var tempfile = os.Args[3]

	// Cleanup before download
	if _, err := os.Stat(file); err == nil {
		fmt.Println("Removing old file: ", file)
		err2 := os.Remove(tempfile)
		if err2 != nil {
			fmt.Printf("Could not remove old file: %s", err2)
		}
	}

	DownloadToFile(url, tempfile, file)
	unzipped := UnzipFile(file)
	fmt.Println("From main - unzipped: ", unzipped)
}
