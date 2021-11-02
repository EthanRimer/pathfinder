package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type page struct {
    title    string
    link     string
    children []string
}

var currentPage page = *new(page)
var rootPage string = "https://www.peanuts.com"
var visitedLinks map[string]bool = make(map[string]bool)
var unvisitedLinks map[string]bool = make(map[string]bool)
var hierarchy map[string]page = make(map[string]page)

func main() {

    resp, err := http.Get(rootPage)
    checkError(err)

    log.Println("Visited page: " + rootPage)
    currentPage.link = rootPage
    add(visitedLinks, rootPage)

    doc, err := html.Parse(resp.Body)
    checkError(err)

    if title, present := getPageTitle(doc); present {

        currentPage.title = title
    }
    findLinks(rootPage, doc)

    hierarchy[currentPage.title] = currentPage
    currentPage = *new(page)

    for !setIsEmpty(unvisitedLinks) {
        for link := range unvisitedLinks {
            if has(visitedLinks, link) {

                remove(unvisitedLinks, link)

                continue
            }

            currentPage.link = link

            time.Sleep(time.Second)
            resp, err := http.Get(link)
            checkError(err)

            log.Printf("Visited page: %v", link)
            add(visitedLinks, link)

            doc, err := html.Parse(resp.Body)
            checkError(err)

            title, pageHasTitle := getPageTitle(doc)
            _, present := hierarchy[title]
            
            if pageHasTitle && !present {

                currentPage.title = title
            } else {

                currentPage.title = link
            }

            findLinks(link, doc)

            hierarchy[currentPage.title] = currentPage

            currentPage = *new(page)
        }
    }

    fmt.Printf("\n\nFound %d unique links\n\n", len(visitedLinks))

    for title, webpage := range hierarchy {

        fmt.Printf("Page Title: %v\nPage Link: %v", title, webpage.link)
        if len(webpage.children) != 0 {
            fmt.Println("\nChildren:")
            for _, link := range webpage.children {

                fmt.Println(link)
            }
        }

        fmt.Println("\n")
    }
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
                    currentPage.children = append(currentPage.children, link)

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

func getPageTitle(n *html.Node) (string, bool) {
    if n.Type == html.ElementNode && n.Data == "title" {

        return n.FirstChild.Data, true
    }

    for c := n.FirstChild; c != nil; c = c.NextSibling {
        if title, present := getPageTitle(c); present {
            
            return title, present
        }
    }

    return "", false
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

    if u.Fragment != "" {
        
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

