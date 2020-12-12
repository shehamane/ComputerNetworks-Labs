package main

import (
	"fmt"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"net/url"
	"time"
)

func showForm(w http.ResponseWriter) {
	fmt.Fprintf(w, `
	<div>
		<form action="/proxy" method="get">
			<input type="text" name="url"/>
			<input type="submit"/> 
		</form>
	</div>`)
}

func getAttr(node *html.Node, key string) string {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func changeAttr(node *html.Node, key string, url string) {
	for i, attr := range node.Attr {
		if attr.Key == key {
			node.Attr[i].Val = url + node.Attr[i].Val
		}
	}
}

func isElem(node *html.Node, tag string) bool {
	return node != nil && node.Type == html.ElementNode && node.Data == tag
}

func findAndFix(node *html.Node, URL string) {
	if node == nil{
		return
	}
	if isElem(node, "a") {
		parsedURL, _ := url.Parse(getAttr(node, "href"))
		if parsedURL.Scheme == "" && parsedURL.Host == "" {
			changeAttr(node, "href", URL+"/")
		}
		changeAttr(node, "href", "/proxy?url=")

	}
	if isElem(node, "img") {
		parsedURL, _ := url.Parse(getAttr(node, "src"))
		if parsedURL.Scheme == "" && parsedURL.Host == "" {
			changeAttr(node, "src", URL+"/")
		}
	}
	if isElem(node, "link") {
		parsedURL, _ := url.Parse(getAttr(node, "href"))
		if parsedURL.Scheme == "" && parsedURL.Host == "" {
			changeAttr(node, "href", URL+"/")
		}
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		findAndFix(child, URL)
	}
}

func fixLinks(resp *http.Response, url string) *html.Node {
	doc, err := html.Parse(resp.Body)
	if err != nil {
		panic(err)
	}
	findAndFix(doc, url)
	return doc
}

func proxy(w http.ResponseWriter, URL string) {
	resp, err := http.Get(URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	parsedURL, _ := url.Parse(URL)
	mainURL := parsedURL.Scheme + "://" + parsedURL.Host + "/"
	doc := fixLinks(resp, mainURL)
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	html.Render(w, doc)
}

func handleHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.URL.Path == "/favicon.ico" {
		return
	} else if r.URL.Path == "/proxy" {
		proxy(w, r.Form.Get("url"))
	} else {
		showForm(w)
	}
}
func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		if k == "Content-Security-Policy"{
			continue
		}
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func main() {
	server := &http.Server{
		Addr: ":3000",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handleHTTP(w, r)
			fmt.Printf("Подключение к %s...", r.URL.Path)
		}),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}

//curl -Lv —proxy http://localhost:3000 http://www.google.com
