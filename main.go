package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
    "sort"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type page struct {
    Title    string
    Link     string
    Children []string
}

type Hierarchy struct {
    Links []string
    Pages map[string]page
}

var currentPage page = *new(page)
var rootPage string = "https://www.peanuts.com"
var visitedLinks map[string]bool = make(map[string]bool)
var unvisitedLinks map[string]bool = make(map[string]bool)
//var Hierarchy map[string]page = make(map[string]page)

func main() {
    pageMap := make(map[string]page)

    linksFile, err := os.Create("links.html")
    checkError(err)
    defer linksFile.Close()

    resp, err := http.Get(rootPage)
    checkError(err)

    log.Println("Visited page: " + rootPage)
    currentPage.Link = rootPage
    add(visitedLinks, rootPage)

    doc, err := html.Parse(resp.Body)
    checkError(err)

    if title, present := getPageTitle(doc); present {

        currentPage.Title = title
    }
    findLinks(rootPage, doc)

    //Hierarchy[currentPage.Link] = currentPage
    pageMap[currentPage.Link] = currentPage
    currentPage = *new(page)

    for !setIsEmpty(unvisitedLinks) {
        for link := range unvisitedLinks {
            if has(visitedLinks, link) {

                remove(unvisitedLinks, link)

                continue
            }

            currentPage.Link = link

            time.Sleep(time.Second)
            resp, err := http.Get(link)
            checkError(err)

            log.Printf("Visited page: %v", link)
            add(visitedLinks, link)

            doc, err := html.Parse(resp.Body)
            checkError(err)

            title, pageHasTitle := getPageTitle(doc)
            
            if pageHasTitle {

                currentPage.Title = title
            } else {

                currentPage.Title = link
            }

            findLinks(link, doc)

            //Hierarchy[currentPage.Link] = currentPage
            pageMap[currentPage.Link] = currentPage

            currentPage = *new(page)
        }
    }

    fmt.Printf("\n\nFound %d unique links\n\n", len(visitedLinks))

    links := make([]string, 0, len(pageMap))
    for link := range pageMap {
        links = append(links, pageMap[link].Link)
    }

    sort.Strings(links)
    h := Hierarchy{Links: links, Pages: pageMap}

    t, err := template.ParseFiles("templ.html")
    t.Execute(linksFile, h)
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
                    currentPage.Children = append(currentPage.Children, link)

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

