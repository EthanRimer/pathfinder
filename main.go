package main

import (
	"fmt"
	"golang.org/x/net/html"
    "log"
	"net/http"
	"net/url"
    "strings"
    "time"
)

var rootPage string = "https://www.peanuts.com"
var visitedLinks map[string]bool = make(map[string]bool)
var unvisitedLinks map[string]bool = make(map[string]bool)
var linkQueue []string

func main() {

    resp, err := http.Get(rootPage)
    checkError(err)

    log.Println("Visited page: " + rootPage)
    add(visitedLinks, rootPage)

    doc, err := html.Parse(resp.Body)
    checkError(err)

    findLinks(rootPage, doc)
    for !isEmpty() {

        var link string = dequeue()
        if _, present := visitedLinks[link]; present {

            continue
        }

        time.Sleep(time.Second)
        resp, err := http.Get(link)
        checkError(err)

        log.Println("Visited page: " + link)
        add(visitedLinks, link)

        doc, err := html.Parse(resp.Body)
        checkError(err)

        findLinks(link, doc)
    }

    fmt.Printf("\n\nFound %v unique links\n", len(visitedLinks))
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
                    enqueue(link)

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

// *** QUEUE FUNCTIONS ***
func enqueue(s string) {

    linkQueue = append(linkQueue, s)
}

func dequeue() string {

    var s string = linkQueue[0]
    linkQueue = linkQueue[1:]

    return s
}

func front() string {
    return linkQueue[0]
}

func isEmpty() bool {
    return len(linkQueue) == 0
}

func isInQueue(link string) bool {
    for _, s := range linkQueue {
        if s == link {
            return true
        }
    }

    return false
}

