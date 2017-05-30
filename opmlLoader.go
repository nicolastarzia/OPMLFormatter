package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/SlyMarbo/rss"
	"github.com/gilliek/go-opml/opml"
)

func requestHTTP(outlin opml.Outline, bufWRT *bufio.Writer) {
	currentLink := outlin.XMLURL
	feed, err := rss.Fetch(currentLink)

	if err == nil {
		if len(feed.Items) == 0 {
			fmt.Printf("%s (DEPRECATED - NO POST)\n", currentLink)
			return
		}
		if feed.Items[0].Date.Before(time.Now().AddDate(0, -5, 0)) {
			fmt.Printf("%s - %s (DEPRECATED - older)\n", currentLink, feed.Items[0].Date)
		} else {
			fmt.Printf("%s - %s (added)\n", currentLink, feed.Refresh)
			bufWRT.WriteString("<outline type=\"rss\" text=\"" + outlin.Text + "\" title=\"" + outlin.Title + "\" htmlUrl=\"" + outlin.HTMLURL + "\" xmlUrl=\"" + outlin.XMLURL + "\"/>\n")
		}
	} else {
		fmt.Printf("%q - %s\n", currentLink, err)
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
