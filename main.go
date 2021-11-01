package main

import (
    "encoding/csv"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
    "strings"
	"time"

	"golang.org/x/net/html"
)

var rootPage string = "https://www.peanuts.com"
var visitedLinks map[string]bool = make(map[string]bool)
var unvisitedLinks map[string]bool = make(map[string]bool)
var hierarchy map[string][]string = make(map[string][]string)
var children []string

func main() {

    linksFile, err := os.Create("links.csv")
    checkError(err)
    defer linksFile.Close()

    resp, err := http.Get(rootPage)
    checkError(err)

    log.Println("Visited page: " + rootPage)
    add(visitedLinks, rootPage)

    doc, err := html.Parse(resp.Body)
    checkError(err)

    csvWriter := csv.NewWriter(linksFile)
    csvWriter.Write([]string{"link", "children"})
    findLinks(rootPage, doc)

    for !setIsEmpty(unvisitedLinks) {
        for link := range unvisitedLinks {
            if has(visitedLinks, link) {
                remove(unvisitedLinks, link)
                continue
            }

            time.Sleep(time.Second)
            resp, err := http.Get(link)
            checkError(err)

            log.Printf("Visited page: %v", link)
            add(visitedLinks, link)

            doc, err := html.Parse(resp.Body)
            checkError(err)

            children = make([]string, 0)
            findLinks(link, doc)
            hierarchy[link] = children
        }
    }

    fmt.Printf("\n\nFound %d unique links\n\n", len(visitedLinks))

    for link, children := range hierarchy {

        csvWriter.Write([]string{link, strings.Join(children, " ")})
    }

    csvWriter.Flush()
}

func findLinks(rootLink string, n *html.Node) {

    if n.Type == html.ElementNode && n.Data == "a" {
        for _, a := range n.Attr {
            if a.Key == "href" {

                var link string = a.Val
                
                if strings.Contains(link, rootLink) || string(link[0]) == "/" {
                    
                    link = formatLink(link)

                    if !validLink(link) {

                        break
                    }

                    add(unvisitedLinks, link)
                    children = append(children, link)

                    log.Println("Found page: " + link)
                } 

                break
            }
        }
    }

    for c := n.FirstChild; c != nil; c = c.NextSibling {

        findLinks(rootLink, c)
    }
}


func checkError(err error) {
    if err != nil {

        log.Fatalln(err)
    }
}

func formatLink(link string) string {
    if string(link[0]) == "/" {

        link = rootPage + link
    }
    if string(link[len(link) - 1]) == "/" {

        link = link[:len(link) - 1]
    }

    return link
}

func validLink(link string) bool {
    u, err := url.Parse(link)
    checkError(err)

    if u.Fragment != "" || len(u.Query()) != 0 {
        
        return false
    }
    if has(unvisitedLinks, link) {

        return false
    }
    if has(visitedLinks, link) {

        return false
    }

    return true
}

// *** SET FUNCTIONS ***
func add(links map[string]bool, str string) {
    if _, present := links[str]; !present {

        links[str] = true
    }
}

func remove(links map[string]bool, str string) {
    if _, present := links[str]; present {

        delete(links, str)
    }
}

func has(links map[string]bool, str string) bool {

    _, present := links[str]

    return present
}

func setIsEmpty(links map[string]bool) bool {

    return len(links) == 0
}

