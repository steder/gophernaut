package gophernaut

import (
	"fmt"
	"html/template"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"strings"
)

// GetGopherHandler ...
func GetGopherHandler(pool Pool) func(w http.ResponseWriter, r *http.Request) {
	var requestCount = 0

	staticHandler := http.StripPrefix("/static", http.FileServer(http.Dir("static")))
	adminTemplate := template.Must(template.ParseFiles("templates/admin.html"))
	adminHandler := func(w http.ResponseWriter, req *http.Request) {
		adminTemplate.Execute(w, nil)
	}

	myHandler := func(responseWriter http.ResponseWriter, myReq *http.Request) {
		requestPath := myReq.URL.Path
		// DONE: adjust request host to assign the request to the appropriate child process
		// TODO: multiprocess, pick one of n hostnames based on pool status
		worker := pool.GetWorker()
		hostname := worker.Hostname
		requestCount++
		targetURL, _ := url.Parse(hostname)

		director := func(req *http.Request) {
			targetQuery := targetURL.RawQuery
			req.URL.Scheme = targetURL.Scheme
			req.URL.Host = targetURL.Host

			// clean up but preserve trailing slash:
			trailing := strings.HasSuffix(req.URL.Path, "/")
			req.URL.Path = path.Join(targetURL.Path, req.URL.Path)
			if trailing && !strings.HasSuffix(req.URL.Path, "/") {
				req.URL.Path += "/"
			}

			// preserve query string:
			if targetQuery == "" || req.URL.RawQuery == "" {
				req.URL.RawQuery = targetQuery + req.URL.RawQuery
			} else {
				req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
			}
		}

		proxy := &httputil.ReverseProxy{Director: director}

		worker.StartRequest()
		defer worker.CompleteRequest()

		switch {
		case requestPath == "/admin":
			adminHandler(responseWriter, myReq)
			return
		case strings.HasPrefix(requestPath, "/static"):
			staticHandler.ServeHTTP(responseWriter, myReq)
			return
		}

		proxy.ServeHTTP(responseWriter, myReq)
	}

	return myHandler
}
