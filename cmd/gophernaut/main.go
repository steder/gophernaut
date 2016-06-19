package main

import (
	"fmt"
	"log"
	"net/http"

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
	log.Fatal(http.ListenAndServe(":8483", http.HandlerFunc(gophernaut.GetGopherHandler(hostnames))))
}
