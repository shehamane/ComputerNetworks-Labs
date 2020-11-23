package main

import (
	"fmt"
	"golang.org/x/net/html"
	"log"
	"net/http"
)

func getChildren(node *html.Node) []*html.Node {
	var children []*html.Node
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		children = append(children, c)
	}
	return children
}

func getAttr(node *html.Node, key string) string {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func isText(node *html.Node) bool {
	return node != nil && node.Type == html.TextNode
}

func isElem(node *html.Node, tag string) bool {
	return node != nil && node.Type == html.ElementNode && node.Data == tag
}

func findFirst(node *html.Node, tag string) *html.Node{
	for c := node.FirstChild; c!=nil; c = c.NextSibling{
		if c.Data == tag{
			return c
		}
	}
	return nil
}

func findAll(node *html.Node, tag string) []*html.Node{
	var ans []*html.Node
	for c := node.FirstChild; c!=nil; c = c.NextSibling{
		if c.Data == tag{
			ans = append(ans, c)
		}
	}
	return ans
}

type Item struct {
	time, title, text, href string
}

func makeItem(node *html.Node) *Item{
	var item Item
	item.text = ""
	item.time = getChildren(getChildren(findFirst(node, "div"))[1])[0].Data
	if findFirst(node, "h2") != nil {
		item.title = getChildren(getChildren(findFirst(node, "h2"))[0])[0].Data
		item.href = getAttr(getChildren(findFirst(node, "h2"))[0], "href")
	}
	for _,p := range findAll(node, "p"){
		if isText(p.FirstChild){
			item.text += p.FirstChild.Data
		}else if (p.FirstChild != nil){
			item.text += findFirst(p, "a").FirstChild.Data
			item.text += ";\n"
		}
	}

	return &item
}

func takeItems(node *html.Node) []*Item {
	if isElem(node, "div") && getAttr(node, "class") == "newsline" {
		var items []*Item
		for _, itemNode := range getChildren(node) {
			if !isElem(itemNode, "div") || !(getAttr(itemNode, "class")=="nl-item"){
				continue
			}
			items = append(items, makeItem(itemNode))
		}
		return items
	} else {
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			result := takeItems(child)
			if result != nil {
				return result
			}
		}
	}
	return nil
}

func getNews(url string) []*Item {
	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusOK {
		doc, err := html.Parse(response.Body)
		if err != nil {
			panic(err)
		}
		return takeItems(doc)
	}
	return nil
}

func showNews(w http.ResponseWriter, url string) {
	items := getNews(url)
	if items==nil{
		fmt.Fprintf(w,"<h1>ERROR404</h1>")
		return
	}
	for _, news := range items {
		fmt.Fprintf(w, `
	<div style = "border: 1px dashed black; padding: 0 10px;">
		<a href="%s" style="text-decoration: none;"><h3>%s</h3></a>
		<p>%s</p><br>
		<span style="color: grey;">%s</span>
	</div>`, news.href, news.title, news.text, news.time)
	}
}

func HomeRouterHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.URL.Path == "/favicon.ico" {
		return
	}
	showNews(w, "https://www.sports.ru/ufc/")
}

func main() {
	http.HandleFunc("/", HomeRouterHandler)
	err := http.ListenAndServe(":3029", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}