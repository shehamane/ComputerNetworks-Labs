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

func isText(node *html.Node) bool {
	return node != nil && node.Type == html.TextNode
}

func getAttr(node *html.Node, key string) string {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func isElem(node *html.Node, tag string) bool {
	return node != nil && node.Type == html.ElementNode && node.Data == tag
}

func findFirst(node *html.Node, tag string) *html.Node {
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Data == tag {
			return c
		}
	}
	return nil
}

func findFirstClass(node *html.Node, tag string, class string) *html.Node {
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Data == tag && getAttr(c, "class") == class {
			return c
		}
	}
	return nil
}

type Item struct {
	time, title, text, img, src string
}

func makeItem(node *html.Node) *Item {
	var item Item

	item.img = "https://mipt.ru" + getAttr(findFirst(findFirstClass(node, "a", "post-picture"), "img"), "src")
	title := findFirstClass(findFirstClass(node, "div", "post-contents"), "div", "post-title-block")
	item.time = findFirstClass(title, "span", "post-date").FirstChild.Data
	item.src = getAttr(findFirstClass(title, "a", "post-title post-link"), "href")
	item.title = findFirstClass(title, "a", "post-title post-link").FirstChild.Data
	post_summary := findFirstClass(findFirstClass(node, "div", "post-contents"), "div", "post-summary")
	if isText(post_summary.FirstChild) {
		item.text = post_summary.FirstChild.Data
	} else {
		texts := getChildren(post_summary.FirstChild)
		item.text = ""
		for _, text := range texts {
			if isText(text) {
				item.text += text.Data
			} else {
				item.text += text.FirstChild.Data
			}
		}
	}

	return &item
}

func takeItems(node *html.Node) []*Item {
	if isElem(node, "div") && getAttr(node, "class") == "post-list" {
		var items []*Item
		for col := node.FirstChild; col != nil; col = col.NextSibling {
			if !(isElem(col, "div") && getAttr(col, "class") == "col") {
				continue
			}
			for post := col.FirstChild; post != nil; post = post.NextSibling {
				if !(isElem(post, "div") && getAttr(post, "class") == "post post-type-news") {
					continue
				}
				items = append(items, makeItem(post))
			}
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

func getNewsList(url string) []*Item {
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

func showNewsList(w http.ResponseWriter, url string) {
	items := getNewsList(url)
	if items == nil {
		fmt.Fprintf(w, "<h1>ERROR404</h1>")
		return
	}
	for _, news := range items {
		fmt.Fprintf(w, `
	<div style = "border: 1px solid black; padding: 0; display: flex; flex-direction: row;">
		<img src="%s">
		<div style = "display: flex; flex-direction: column; margin-left: 25px;'">
			<a href="%s"><h3>%s</h3></a>
			<p>%s</p>
			<div style="color: grey;">%s</div>
		</div>
	</div>`, news.img, news.src, news.title, news.text, news.time)
	}
}

func takeNews(w http.ResponseWriter, node *html.Node) {
	if isElem(node, "div") && getAttr(node, "class") == "post-details news-details" {
		html.Render(w, node)
		if findFirstClass(node, "div", "post-picture") != nil {
			fmt.Fprintf(w, `
	<img src="https://mipt.ru%s">
`, getAttr(getChildren(findFirstClass(node, "div", "post-picture"))[1], "src"))
		}
	} else {
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			takeNews(w, child)
		}
	}
}

func showNews(w http.ResponseWriter, url string) {
	response, err := http.Get("https://mipt.ru" + url)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusOK {
		doc, err := html.Parse(response.Body)
		if err != nil {
			panic(err)
		}
		takeNews(w, doc)
	}
}

func HomeRouterHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.URL.Path == "/favicon.ico" {
		return
	}
	if r.URL.Path == "/" {
		showNewsList(w, "https://mipt.ru/")
	} else {
		showNews(w, r.URL.Path)
	}
}

func main() {
	http.HandleFunc("/", HomeRouterHandler)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
