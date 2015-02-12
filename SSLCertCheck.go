package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

// Global declaration for Colored output
var green = color.New(color.FgGreen).SprintFunc()

// URL for the Alexa top sites
const URL string = "http://s3.amazonaws.com/alexa-static/top-1m.csv.zip"

// UnzipFile unzips a file
func UnzipFile(_name string) (unzipped string) {
	zipfile := _name
	fmt.Println("Opening to unzip: ", green(zipfile))

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
			fmt.Println("Decompressing: ", green(path))
		}
	}
	return unzipped
}

// DownloadFromURL downloads a file from a URL
// Saves the file in the same path as the executable
func DownloadFromUrl(url string) (name string) {
	tokens := strings.Split(URL, "/")
	filename := tokens[len(tokens)-1]
	fmt.Println("Downloading ",url, "to", filename)

	output, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		log.Fatal (err)
	}

	fmt.Println(n, " bytes downloaded")
	name = filename
	return name
}

// Main works as follows:
// 1. Delete the target filename, before any download takes place
// 2. Download the file using DownloadToFile
// 3. Unzips the Alexa file using UnzipFile
func main() {
	alexaFile := DownloadFromUrl(URL)
	unzipped := UnzipFile(alexaFile)
	fmt.Println("From main - unzipped: ", unzipped)
	os.Remove(alexaFile)
	fmt.Println("Removed: ", alexaFile)
}
