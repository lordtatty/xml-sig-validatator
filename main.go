package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/ma314smith/signedxml"
)

// Supplied by goreleaser https://goreleaser.com/cookbooks/using-main.version
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func verStr() string {
	if version == "dev" {
		return "Local development build"
	}
	return fmt.Sprintf("version: %s\ncommit: %s\nbuilt: at %s", version, commit, date)
}

func main() {
	// if "--version" is supplied, print version and exit
	verFlag := flag.Bool("version", false, "print version and exit")
	flag.Parse()
	if *verFlag {
		fmt.Println(verStr())
		return
	}

	if len(os.Args) < 2 {
		fmt.Println("URL or filepath required as the first argument")
		fmt.Println("Or use --version to print version and build info")
		return
	}

	path := os.Args[1]

	xml, err := BodyOf(path)
	if err != nil {
		log.Fatal(err)
	}

	_, err = Validate(xml)
	if err != nil {
		log.Println("SIGNATURE FAILED to validate")
		log.Fatal(err)
	}
	fmt.Println("Signature is Valid")
}

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

func XMLFromFile(path string) ([]byte, error) {
	xml, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return xml, nil
}

func XMLFromUrl(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
