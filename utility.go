package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"golang.org/x/net/html"
)

func FileCreate(name, content string) error {
	//fh, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	fh, err := os.Create(name)
	if err != nil {
		log.Println(err)
		return err
	}
	defer fh.Close()
	_, err = fh.WriteString(content)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func tagAttrMap(attrs []html.Attribute) map[string]string {
	var omap = make(map[string]string)
	for _, attr := range attrs {
		omap[attr.Key] = attr.Val
	}
	return omap
}

func parse_html(n *html.Node) {

	if n.Type == html.ElementNode {

		for _, element := range n.Attr {
			if element.Key == "class" {
				//fmt.Printf("class: %s\n", element.Val)
				classList = append(classList, element.Val)
			} else {
				//fmt.Println(">", element.Namespace, element.Key, element.Val)
			}
		}

		if n.Data == "link" {
			tmap := tagAttrMap(n.Attr)
			linkRows = append(linkRows, tmap)
			//fmt.Println("link=>", n.Data, tmap, tmap["href"])
		} else {
			//fmt.Println("=>", n.Data, len(n.Attr))
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parse_html(c)
	}
}

func templateParser(txtTemplate string, tmap map[string]interface{}) string {

	var tplOutput bytes.Buffer
	tpl := template.New("parser")
	tmplt, err := tpl.Parse(txtTemplate)
	if err != nil {
		log.Println("templateParserERR", err)
		return ""
	}
	err = tmplt.Execute(&tplOutput, tmap)
	if err != nil {
		log.Println("templateExecuteERR", err)
		return ""
	}
	return tplOutput.String()
}

//d:flex
func cssOutput(tmap map[string]interface{}) string {

	value, isExist := tmap["value"].(string)
	if isExist {
		tmap["value"] = valueReplacer(value)
	}
	cssTemplate := ".{{.identifier}}{\n  {{.key}}: {{.value}};\n}"
	return templateParser(cssTemplate, tmap)
}

func valueReplacer(value string) string {

	var output string
	for _, char := range value {
		if char == '_' {
			output += strings.Replace(string(char), "_", " ", 1)
		} else {
			output += fmt.Sprintf("%c", char)
		}
	}
	return output
}

func specialCharReplacer(name string) string {

	//bt\:1px_solid_\#e2e2e3
	var output string
	for _, char := range name {

		if char == ':' {
			output += strings.Replace(string(char), ":", "\\:", 1)
			//fmt.Printf("%c %v\n", char, char)
		} else if char == '#' {
			output += strings.Replace(string(char), "#", "\\#", 1)

		} else if char == '.' {
			output += strings.Replace(string(char), ".", "\\.", 1)

		} else if char == '[' {
			output += strings.Replace(string(char), "[", "\\[", 1)

		} else if char == ']' {
			output += strings.Replace(string(char), "]", "\\]", 1)

		} else if char == '%' {
			output += strings.Replace(string(char), "%", "\\%", 1)

		} else if char == '!' {
			output += strings.Replace(string(char), "!", "\\!", 1)

		} else {
			output += fmt.Sprintf("%c", char)
		}
		//xlg:translateX_-50%
		//fmt.Println(output)
	}

	return output

}

var pMap = map[string]string{
	"b":           "border", //1
	"c":           "color",
	"d":           "display",
	"f":           "font",  //
	"g":           "grid",  //5
	"i":           "inset", //
	"m":           "margin",
	"o":           "overflow",
	"p":           "padding",
	"s":           "scale", //10
	"h":           "height",
	"v":           "visibility", //
	"w":           "width",
	"z":           "zoom", //14
	"bg":          "background",
	"bs":          "border-style",
	"fs":          "font-size",
	"fw":          "font-weight",
	"ml":          "margin-left",
	"mr":          "margin-right", //20
	"mt":          "margin-top",
	"mb":          "margin-bottom",
	"pt":          "padding-top",
	"pb":          "padding-bottom",
	"pl":          "padding-left",
	"pr":          "padding-right",
	"ls":          "letter-spacing", //27
	"lh":          "line-height",
	"bt":          "border-top",
	"bb":          "border-bottom",
	"fd":          "flex-direction",
	"fg":          "flex-grow",
	"jc":          "justify-content",
	"ai":          "align-items",
	"mw":          "max-width",
	"mh":          "max-height",
	"ws":          "white-space",
	"va":          "vertical-align",
	"zi":          "z-index",
	"td":          "text-decoration", //40
	"bc":          "border-color",
	"ta":          "text-align",
	"br":          "border-radius", //43
	"bts":         "border-top-style",
	"pos":         "position",
	"top":         "top",
	"gap":         "gap",
	"btc":         "border-top-color",
	"bgc":         "background-color",
	"bgi":         "background-image", //50
	"gtc":         "grid-template-columns",
	"flex":        "flex",
	"left":        "left",
	"right":       "right",
	"list":        "list-style", //li-st 55
	"minh":        "min-height",
	"maxh":        "max-height",
	"order":       "order",
	"cursor":      "cursor",
	"flex-shrink": "flex-shrink",
	"flex-wrap":   "flex-wrap",
	"opacity":     "opacity",
	"transition":  "transition",
	"filter":      "filter",
	"transf":      "transform",
	"content":     "content", //66
}

func cssClassParser(singleClassName string) string {

	var key, value string
	slc := strings.Split(singleClassName, ":")
	identifier := specialCharReplacer(singleClassName) //bt\:1px_solid_\#e2e2e3,d:flex, xlg\:translateX_-50\%
	//fmt.Println(identifier)
	if len(slc) == 2 {
		key = slc[0]
		value = slc[1]
	}

	tmap := map[string]interface{}{
		"identifier": identifier,
		"key":        pMap[key],
		"value":      value,
	}
	css := cssOutput(tmap)
	//fmt.Println(css)
	return css

}
