package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"strings"

	"github.com/steder/gophernaut"
)

var hostnames = []string{
	fmt.Sprintf("http://127.0.0.1:%d", 8080),
	fmt.Sprintf("http://127.0.0.1:%d", 8081),
}

var executables = []string{
	fmt.Sprintf("python -m SimpleHTTPServer %d", 8080),
	fmt.Sprintf("python -m SimpleHTTPServer %d", 8081),
}

var requestCount = 0

func myHandler(w http.ResponseWriter, myReq *http.Request) {
	requestPath := myReq.URL.Path
	// DONE: adjust request host to assign the request to the appropriate child process
	// TODO: multiprocess, pick one of n hostnames based on pool status
	hostname := hostnames[requestCount%len(hostnames)]
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

	staticHandler := http.StripPrefix("/static", http.FileServer(http.Dir("static")))
	adminTemplate := template.Must(template.ParseFiles("templates/admin.html"))
	adminHandler := func(w http.ResponseWriter, req *http.Request) {
		adminTemplate.Execute(w, nil)
	}

	switch {
	case requestPath == "/admin":
		adminHandler(w, myReq)
		return
	case strings.HasPrefix(requestPath, "/static"):
		staticHandler.ServeHTTP(w, myReq)
		return
	}
	proxy.ServeHTTP(w, myReq)
}

func main() {
	log.SetPrefix("gophernaut ")
	log.SetFlags(log.Ldate | log.Ltime)
	c := gophernaut.ReadConfig()
	log.Printf("Host %s and Port %d\n", c.Host, c.Port)

	pool := gophernaut.Pool{Executables: executables}
	pool.Start()
	go pool.ManageProcesses()

	log.Printf("Gophernaut is gopher launch!\n")
	// TODO: our own ReverseProxy implementation of at least, ServeHTTP so that we can
	// monitor the response codes to track successes and failures
	log.Fatal(http.ListenAndServe(":8483", http.HandlerFunc(myHandler)))
}
