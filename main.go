package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/ma314smith/signedxml"
)

// Validate an enveloped xml signature
func main() {
	if len(os.Args) < 2 {
		log.Fatal("URL or filepath required as the first argument")
	}

	path := os.Args[1]

	xml, err := BodyOf(path)
	if err != nil {
		log.Fatal(err)
	}

	// Validate the xml
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

func BodyOf(path string) ([]byte, error) {
	if IsUrl(path) {
		return XMLFromUrl(path)
	}
	return XMLFromFile(path)
}

// get xm from file
func XMLFromFile(path string) ([]byte, error) {
	xml, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return xml, nil
}

// get body as bytes from url
func XMLFromUrl(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// Is valid url
func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
