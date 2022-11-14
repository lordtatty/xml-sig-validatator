package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/ma314smith/signedxml"
)

// Validate an enveloped xml signature
func main() {
	// command line flag for url
	url := flag.String("url", "https://www.w3.org/TR/xmldsig-core/xml-stylesheet.txt", "url to validate")
	flag.Parse()

	if *url == "" {
		log.Fatal("--url is required")
	}

	fmt.Println("Checking url:", *url)

	// Download the xml file
	resp, err := http.Get(*url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	xml, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Validate the xml file
	_, err = Validate(xml)
	if err != nil {
		log.Println("SIGNATURE FAILED to validate")
		log.Fatal(err)
	}
	fmt.Println("Signature is Valid")
}

// validate the xml file
func Validate(xml []byte) (bool, error) {
	v, err := signedxml.NewValidator(string(xml))
	if err != nil {
		return false, err
	}
	_, err = v.ValidateReferences()
	if err != nil {
		return false, err
	}
	return true, nil
}
