package main

import (
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"

	"github.com/SudiMango/stratus/util/config"
)

func main() {
	cfg, err := config.GetConfig("config.yaml")
	if err != nil {
		log.Fatalf("Error loading config file: %v", err)
	}

	err = cfg.Validate()
	if err != nil {
		log.Fatalf("Error validating config: \n%v", err)
	}

	mux := http.NewServeMux()
	httpHandler := http.Handler(mux)
	allowedHosts := make(map[string]struct{})

	for _, route := range cfg.Routes {
		allowedHosts[strings.ToLower(route.Host)] = struct{}{}
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

	if cfg.Proxy.Http.RedirectToHttps {
		httpHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			host := r.Host

			if hostname, _, err := net.SplitHostPort(host); err == nil {
				host = hostname
			}

			host = strings.ToLower(strings.TrimSuffix(host, "."))
			if _, exists := allowedHosts[host]; !exists {
				http.Error(w, "unknown host", http.StatusMisdirectedRequest)
				return
			}

			// For local dev where port isnt automatically 443
			if cfg.Proxy.Https.Port != 443 {
				host = net.JoinHostPort(host, strconv.Itoa(cfg.Proxy.Https.Port))
			}

			target := url.URL{
				Scheme:   "https",
				Host:     host,
				Path:     r.URL.Path,
				RawPath:  r.URL.RawPath,
				RawQuery: r.URL.RawQuery,
			}

			http.Redirect(w, r, target.String(), http.StatusPermanentRedirect)
		})
	}

	if cfg.Proxy.Http.Enabled {
		go func() {
			httpServer := &http.Server{
				Addr:              ":" + strconv.Itoa(cfg.Proxy.Http.Port),
				Handler:           httpHandler,
				ReadTimeout:       cfg.Proxy.ReadTimeout,
				ReadHeaderTimeout: cfg.Proxy.ReadHeaderTimeout,
				WriteTimeout:      cfg.Proxy.WriteTimeout,
				IdleTimeout:       cfg.Proxy.IdleTimeout,
				MaxHeaderBytes:    cfg.Proxy.MaxHeaderBytes,
			}

			log.Printf("HTTP proxy running on port %s", strconv.Itoa(cfg.Proxy.Http.Port))
			if err := httpServer.ListenAndServe(); err != nil {
				log.Fatalf("HTTP server failed: %v", err)
			}
		}()
	}

	if cfg.Proxy.Https.Enabled {
		httpsServer := &http.Server{
			Addr:              ":" + strconv.Itoa(cfg.Proxy.Https.Port),
			Handler:           mux,
			ReadTimeout:       cfg.Proxy.ReadTimeout,
			ReadHeaderTimeout: cfg.Proxy.ReadHeaderTimeout,
			WriteTimeout:      cfg.Proxy.WriteTimeout,
			IdleTimeout:       cfg.Proxy.IdleTimeout,
			MaxHeaderBytes:    cfg.Proxy.MaxHeaderBytes,
		}

		log.Printf("HTTPS proxy running on port %s", strconv.Itoa(cfg.Proxy.Https.Port))
		if err := httpsServer.ListenAndServeTLS(cfg.Proxy.Https.Cert, cfg.Proxy.Https.Key); err != nil {
			log.Fatalf("HTTPS server failed: %v", err)
		}
	} else if cfg.Proxy.Http.Enabled {
		select {}
	}
}
