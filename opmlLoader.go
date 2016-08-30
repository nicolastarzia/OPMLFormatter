package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gilliek/go-opml/opml"
)

func requestHTTP(outlin opml.Outline, bufWRT *bufio.Writer) {
	currentLink := outlin.XMLURL
	resp, err := http.Get(currentLink)
	if err == nil && resp.StatusCode == 200 {
		bufWRT.WriteString("<outline type=\"rss\" text=\"" + outlin.Text + "\" title=\"" + outlin.Title + "\" htmlUrl=\"" + outlin.HTMLURL + "\" xmlUrl=\"" + outlin.XMLURL + "\"/>\n")
		fmt.Printf("%q - %q\n", currentLink, resp.Status)
	} else {
		fmt.Printf("%q - %q\n", currentLink, " inválido")
	}
	bufWRT.Flush()
}

func recursiveOutline(outlin opml.Outline, bufWRT *bufio.Writer) {
	currentLink := outlin.XMLURL
	if outlin.Type == "rss" && len(currentLink) > 0 {
		requestHTTP(outlin, bufWRT)
	} else {
		bufWRT.WriteString("<outline text=\"" + outlin.Text + "\" title=\"" + outlin.Title + "\">\n")
	}
	if outlin.Outlines != nil {
		for myOuts := range outlin.Outlines {
			recursiveOutline(outlin.Outlines[myOuts], bufWRT)
		}
		bufWRT.WriteString("</outline>\n")
	}
}

func main() {
	filename := "feedly.opml"
	fmt.Printf("Filename: %s\n", filename)
	myFile, err := filepath.Abs(filename)
	if err != nil {
		log.Fatal(err)
	}
	doc, err := opml.NewOPMLFromFile(myFile)
	if err != nil {
		log.Fatal(err)
	}

	fil, err := os.Create("output.txt")
	if err != nil {
		log.Fatal(err)
	}

	w := bufio.NewWriter(fil)
	w.WriteString(`<?xml version="1.0" encoding="UTF-8"?>
<opml version="1.0">
    <head>
        <title>Nicolas subscriptions in feedly Cloud</title>
    </head>
    <body>
    `)
	for outlin := range doc.Body.Outlines {
		recursiveOutline(doc.Body.Outlines[outlin], w)
	}
	w.WriteString(`
    </body>
    </opml>`)
	w.Flush()
	fmt.Printf("END!!")
}
