package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strings"
)

var hostname = fmt.Sprintf("http://127.0.0.1:%d", 8080)
var executable string = fmt.Sprintf("python -m SimpleHTTPServer %d", 8080)

func start_process(events chan int) {
	command_parts := strings.Split(executable, " ")
	command := exec.Command(command_parts[0], command_parts[1:]...)
	fmt.Printf("Command: %v\n", command)
	command.Start()
	_, ok := <-events
	if !ok {
		command.Process.Kill()
	}
}

func myHandler(w http.ResponseWriter, req *http.Request) {
	request_path := req.URL.Path

	target_url, _ := url.Parse(hostname)
	director := func(req *http.Request) {
		targetQuery := target_url.RawQuery
		req.URL.Scheme = target_url.Scheme
		// TODO: adjust request host to assign the request to the appropriate child process
		req.URL.Host = target_url.Host
		req.URL.Path = path.Join(target_url.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
	}

	proxy := &httputil.ReverseProxy{Director: director}

	static_handler := http.StripPrefix("/static", http.FileServer(http.Dir("static")))
	admin_template := template.Must(template.ParseFiles("templates/admin.html"))
	admin_handler := func(w http.ResponseWriter, req *http.Request) {
		admin_template.Execute(w, nil)
	}

	fmt.Printf("path: %s\n", request_path)
	switch {
	case request_path == "/admin":
		fmt.Printf("admin path...\n")
		admin_handler(w, req)
		return
	case strings.HasPrefix(request_path, "/static"):
		fmt.Printf("static path...\n")
		static_handler.ServeHTTP(w, req)
		return
	}
	fmt.Printf("proxy path...\n")
	proxy.ServeHTTP(w, req)
}

func main() {
	events_channel := make(chan int)
	go start_process(events_channel)
	fmt.Printf("Gophernaut is ready for eBusiness!\n")
	http.ListenAndServe(":8483", http.HandlerFunc(myHandler))
	// TODO: our own ReverseProxy implementation of at least, ServeHTTP so that we can
	// monitor the response codes to track successes and failures
}
