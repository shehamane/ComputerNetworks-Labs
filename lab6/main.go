package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func handleHTTP(w http.ResponseWriter, req *http.Request) {
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
func copyHeader(dst, src http.Header) {
	for k, vv := range src {
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
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}

//curl -Lv —proxy http://localhost:3000 http://www.google.com