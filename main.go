package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	target, err := url.Parse("http://localhost:3000")
	if err != nil {
		log.Fatalf("Invalid target url: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Proxying request: %s %s -> %s", r.Method, r.URL.Path, target.String())
		proxy.ServeHTTP(w, r)
	})

	log.Println("Reverse proxy running on port 9090")
	if err := http.ListenAndServe(":9090", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
