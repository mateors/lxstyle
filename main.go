package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/net/html"
)

var classList []string
var linkRows []map[string]string

// github.com/andybalholm/cascadia
func htmlNode(htmlFilePath string) error {

	var styleSheetName string
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
	//fmt.Println(doc.Attr)
	parse_html(doc)

	fmt.Println(">>", len(classList))
	for _, cls := range classList {

		slc := strings.Split(cls, " ")
		for _, sClass := range slc {
			fmt.Println(cssClassParser(sClass))
		}

	}

	if len(linkRows) > 0 {
		styleSheetName = linkRows[0]["href"] //get the first style sheet
	}
	fmt.Println(styleSheetName, len(linkRows))

	//cssClassParser("gap:1.5rem") //gap\:1\.5rem

	// sel, err := cascadia.Parse("div.invoice-details")
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// fmt.Println(sel.String())
	return nil
}

func main() {

	htmlNode("index.html")

}
