// Dimitrios Stergiou <dstergiou@gmail.com
// Checks Alexa top 1M sites per country domain

package main

import (
	"archive/zip"
	"encoding/csv"
	"flag"
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
	fmt.Println("Downloading Alexa CSV to", green(filename))

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

	_, err = io.Copy(output, response.Body)
	if err != nil {
		log.Fatal(err)
	}

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

	var allRecords []string
	for _, each := range rawCsvData {
		hostname := each[1]
		allRecords = append(allRecords, hostname)
	}
	return allRecords

}

// HostsTopSites returns only domains under the domains the user inputted
func HostsTopSites(hostnames []string, domain string) (matchedHosts []string) {
	var regexpBuild string
	if domain == "se" {
		regexpBuild = ".*.se$|.*.nu$"
	} else {
		regexpBuild = ".*." + domain + "$"
	}
	domainToMatch := regexp.MustCompile(regexpBuild)
	for _, host := range hostnames {
		isSwedish := domainToMatch.Match([]byte(host))
		if isSwedish {
			matchedHosts = append(matchedHosts, host)
		}
	}
	return matchedHosts
}

// PrintHosts prints a list of hostnames
func PrintHosts(hostnames []string, amount int) {
	index := 0
	for _, host := range hostnames {
		fmt.Println(host)
		index++
		if index == amount {
			break
		}
	}
	fmt.Println("Hosts listed:", green(index))
}

// Main connects to Alexa and downloads the zip file.
// Then it unzips the file and processes the CSV
func main() {
	domainFlag := flag.String("domain", "", "Internet domain to use")
	amountFlag := flag.Int("amount", 10, "Number of hosts to display")
	downloadFlag := flag.Bool("dload", false, "Set the flag for a new download for the CSV file")
	flag.Parse()

	if *domainFlag == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *downloadFlag {
		alexaFile := DownloadFromURL(URL)
		_ = UnzipFile(alexaFile)
		os.Remove(alexaFile)
	}
	hosts := CsvParse("top-1m.csv")
	fmt.Println("For domain: ", green(*domainFlag))
	topHosts := HostsTopSites(hosts, *domainFlag)
	PrintHosts(topHosts, *amountFlag)
}
