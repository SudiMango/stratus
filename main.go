package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"

	"github.com/SudiMango/stratus/util/config"
)

func main() {
	cfg, err := config.GetConfig("config.yaml")
	if err != nil {
		log.Fatalf("Error loading config file: %v", err)
	}

	mux := http.NewServeMux()

	for _, route := range cfg.Routes {
		target, err := url.Parse(config.BuildURL(route))
		if err != nil {
			log.Fatalf("Invalid target url for %q: %v", route.Path, err)
		}

		proxy := &httputil.ReverseProxy{
			Rewrite: func(r *httputil.ProxyRequest) {
				r.SetURL(target)
				r.SetXForwarded()
			},
		}
		mux.HandleFunc(route.Host+route.Path, func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Proxying request: %s %s %s %s -> %s", r.Method, r.Host, r.URL.Path, r.URL.RawQuery, target.String())
			proxy.ServeHTTP(w, r)
		})
	}

	log.Printf("Reverse proxy running on port %s", strconv.Itoa(cfg.Proxy.Port))
	if err := http.ListenAndServe(":"+strconv.Itoa(cfg.Proxy.Port), mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
