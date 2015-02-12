package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

// Main works as follows:
// 1. Delete the target filename, before any download takes place
// 2. Download the file using DownloadToFile
func main() {
	// Cleanup before download
	if _, err := os.Stat(os.Args[2]); err == nil {
		fmt.Println("Removing old file")
		err2 := os.Remove(os.Args[2])
		if err2 != nil {
			fmt.Printf("Could not remove old file: %s", err2)
		}
	}

	DownloadToFile(os.Args[1], os.Args[2], os.Args[3])
}
