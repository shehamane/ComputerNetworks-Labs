package main

import (
	"fmt"
	"github.com/mmcdole/gofeed"
	"log"
	"net/http"
	"sort"
)

func getFeed(url string) *gofeed.Feed {
	//response, _ := http.Get("http://" + url)
	//if response == nil {
	//	return nil
	//}
	//defer response.Body.Close()
	//RSSdata, _ := ioutil.ReadAll(response.Body)
	//RSSparser := rss.Parser{}
	//feed, _ := RSSparser.Parse(strings.NewReader(string(RSSdata)))
	parser := gofeed.NewParser()
	feed, _ := parser.ParseURL("http://" + url)
	return feed
}

func handleFeed(url string) *gofeed.Feed{
	if url == "" {
		sourceURLs := [3]string{"lenta.ru/rss", "vz.ru/rss.xml", "news.mail.ru/rss/90/"}
		feed := new(gofeed.Feed)
		for _, sourceURL := range sourceURLs {
			feedTMP := getFeed(sourceURL)
			for _, item := range feedTMP.Items[:10] {
				feed.Items = append(feed.Items, item)
			}
		}
		sort.Sort(feed)
		return feed
	} else {
		feed := getFeed(url)
		if feed == nil {
			return nil
		}
		return feed
	}
}

func showNews(w http.ResponseWriter, url string) {
	feed := handleFeed(url)
	if feed==nil{
		fmt.Fprintf(w,"<h1>ERROR404</h1>")
		return
	}
	items := feed.Items
	for _, news := range items {
		fmt.Fprintf(w, `
	<div style = "border: 1px dashed black; padding: 0 10px;">
		<h3>%s</h3>
		<p>%s</p><br>
		<span style="color: grey;">%s</span>
		<span style="color: grey;">%s</span>
	</div>`, news.Title, news.Description, news.Published, news.Link)
	}
}

func showLink(w http.ResponseWriter, href string, text string) {
	fmt.Fprintf(w,
		`<a href="%s" style="
		margin-left: 10px;
		font-size: 35px;
		text-decoration: none;
		color: blue;
		padding: 10px 0;">
			%s
		</a>`, href, text)
}

func showLinks(w http.ResponseWriter) {
	fmt.Fprintf(w, `
	<div style="display: flex; justify-content: space-around;">
`)
	showLink(w, "/lenta.ru/rss", "lenta.ru")
	showLink(w, "/vz.ru/rss.xml", "vz.ru")
	showLink(w, "/news.mail.ru/rss/90/", "news.mail.ru")
	fmt.Fprintf(w, "</div>")
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path[1:]
	showNews(w, url)
}

func showHyperText(w http.ResponseWriter, r *http.Request) {
	showLinks(w)
	requestHandler(w, r)
}

func HomeRouterHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.URL.Path == "/favicon.ico" {
		return
	}
	showHyperText(w, r)
}

func main() {
	http.HandleFunc("/", HomeRouterHandler)
	err := http.ListenAndServe(":3029", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
