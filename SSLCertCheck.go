package main

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

// Global declaration for Colored output
var green = color.New(color.FgGreen).SprintFunc()

// URL for Alexa file declaration
const URL string = "http://s3.amazonaws.com/alexa-static/top-1m.csv.zip"

// CA files  declaration
const CA string = "http://curl.haxx.se/ca/cacert.pem"

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
func DownloadFromURL(url string) (name string) {
	tokens := strings.Split(url, "/")
	filename := tokens[len(tokens)-1]
	fmt.Println("Downloading ", url, "to", filename)

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
		log.Fatal(err)
	}

	fmt.Println(n, " bytes downloaded")
	name = filename
	return name
}

// CsvParse splits the Alexa file
func CsvParse(filename string) (hostnames []string) {
	csvfile, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer csvfile.Close()

	reader := csv.NewReader(csvfile)
	rawCsvData, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var hostname string
	var allRecords []string
	for _, each := range rawCsvData {
		hostname = each[1]
		allRecords = append(allRecords, hostname)
	}
	return allRecords

}

// SwedishTopSites returns only domains under.se or .nu
func SwedishTopSites(hostnames []string) (swedishHosts []string) {
	swedishDomains := regexp.MustCompile(`.*.se$|.*.nu$`)
	for _, host := range hostnames {
		isSwedish := swedishDomains.Match([]byte(host))
		if isSwedish {
			swedishHosts = append(swedishHosts, host)
		}
	}
	return swedishHosts
}

// PrintHosts prints a list of hostnames
func PrintHosts(hostnames []string) {
	index := 0
	for _, host := range hostnames {
		fmt.Println("Host is: ", host)
		index++
	}
	fmt.Println("Total number found:", green(index))
}

// Main connects to Alexa and downloads the zip file.
// Then it unzips the file and processes the CSV
func main() {
	//alexaFile := DownloadFromURL(URL)
	//unzipped := UnzipFile(alexaFile)
	//fmt.Println("From main - unzipped: ", unzipped)
	//os.Remove(alexaFile)
	//fmt.Println("Removed: ", alexaFile)
	caFile := DownloadFromURL(CA)
	fmt.Println(caFile)
	fmt.Println("Testing CSV")
	hosts := CsvParse("koko.csv")
	swedishHosts := SwedishTopSites(hosts)
	PrintHosts(swedishHosts)
}
