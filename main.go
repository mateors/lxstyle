package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

var classList []string
var linkRows []map[string]string

func htmlNode(htmlFilePath string) error {

	var styleSheetName, cssOutput string

	osF, err := os.Open(htmlFilePath)
	if err != nil {
		log.Println(err)
		return err
	}
	doc, err := html.Parse(osF)
	if err != nil {
		log.Println(err)
		return err
	}

	//scan|parse the html document
	parse_html(doc)

	//where is the style sheet code generated?
	if len(linkRows) > 0 {
		styleSheetName = linkRows[0]["href"] //get the first style sheet
	}

	for _, cls := range classList {

		slc := strings.Split(cls, " ")
		for _, sClass := range slc {
			cssOutput += cssClassParser(sClass) + "\n"
		}

	}

	//generate style sheet code
	os.MkdirAll(filepath.Dir(styleSheetName), 0755)
	err = FileCreate(styleSheetName, cssOutput)
	return err
}

func main() {

	err := htmlNode("index.html")
	if err != nil {
		log.Println(err)
	}

}
