package main

import (
	"bufio"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/steder/gophernaut"
)

var hostname = fmt.Sprintf("http://127.0.0.1:%d", 8080)
var executable = fmt.Sprintf("python -m SimpleHTTPServer %d", 8080)

func copyToLog(dst *log.Logger, src io.Reader) {
	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		dst.Print(scanner.Text())
	}
}

func startProcess(events chan int) {
	procLog := log.New(os.Stdout, "gopher-worker ", log.Ldate|log.Ltime)
	commandParts := strings.Split(executable, " ")
	command := exec.Command(commandParts[0], commandParts[1:]...)
	fmt.Printf("Command: %v\n", command)

	stdout, err := command.StdoutPipe()
	if err != nil {
		procLog.Fatalln("Unable to connect to stdout from command...")
	}
	stderr, err := command.StderrPipe()
	if err != nil {
		procLog.Fatalln("Unable to connect to stderr from command...")
	}

	//go io.Copy(os.Stdout, stdout)
	//go io.Copy(os.Stderr, stderr)
	go copyToLog(procLog, stdout)
	go copyToLog(procLog, stderr)
	command.Start()

	for {
		_, ok := <-events
		if !ok {
			command.Process.Kill()
		}
	}
}

func myHandler(w http.ResponseWriter, myReq *http.Request) {
	requestPath := myReq.URL.Path

	targetURL, _ := url.Parse(hostname)
	director := func(req *http.Request) {
		targetQuery := targetURL.RawQuery
		req.URL.Scheme = targetURL.Scheme
		// TODO: adjust request host to assign the request to the appropriate child process
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

	//fmt.Printf("path: %s\n", request_path)
	switch {
	case requestPath == "/admin":
		//fmt.Printf("admin path...\n")
		adminHandler(w, myReq)
		return
	case strings.HasPrefix(requestPath, "/static"):
		//fmt.Printf("static path...\n")
		staticHandler.ServeHTTP(w, myReq)
		return
	}
	//fmt.Printf("proxy path...\n")
	proxy.ServeHTTP(w, myReq)
}

func main() {
	// Test reading a config yaml:
	c := gophernaut.ReadConfig()
	fmt.Printf("Host %s and Port %d\n", c.Host, c.Port)

	eventsChannel := make(chan int)
	go startProcess(eventsChannel) // TODO MANY PROCESSES, MUCH POOLS
	fmt.Printf("Gophernaut is gopher launch!\n")
	http.ListenAndServe(":8483", http.HandlerFunc(myHandler))
	// TODO: our own ReverseProxy implementation of at least, ServeHTTP so that we can
	// monitor the response codes to track successes and failures
}
