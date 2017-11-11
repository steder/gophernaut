package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/trace"
	"time"

	"github.com/steder/gophernaut"
)

const debug = true

func main() {
	if debug {
		log.Printf("Enabling tracing")
		traceFile, err := os.Create("trace.out")
		if err != nil {
			panic(err)
		}
		defer traceFile.Close()
		if err := trace.Start(traceFile); err != nil {
			panic(err)
		}
		defer trace.Stop()
	}

	log.SetPrefix("gophernaut ")
	log.SetFlags(log.Ldate | log.Ltime)
	c := gophernaut.ReadConfig()
	log.Printf("Host %s and Port %d\n", c.Host, c.Port)
	log.Printf("Pool %d, %s, %s, %s\n",
		c.Pool.Size,
		c.Pool.Template.Name,
		c.Pool.Template.Executable,
		c.Pool.Template.Hostname,
	)
	log.Printf("Creating pool...\n")

	pool := gophernaut.Pool{Executables: c.GetExecutables(), Hostnames: c.GetHostnames(), Size: c.Pool.Size}
	// TODO: this should block until startup has completed.
	log.Printf("Starting pool...\n")
	pool.Start()

	// TODO: our own ReverseProxy implementation of at least, ServeHTTP so that we can
	// monitor the response codes to track successes and failures
	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", c.Port),
		Handler:        http.HandlerFunc(gophernaut.GetGopherHandler(pool, c)),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go pool.ManageProcesses(server)

	log.Printf("Gophernaut is gopher launch!\n")

	// can't use log.Fatal here as deferred trace.Stop is never called
	log.Print(server.ListenAndServe())
	// On the plus side get to say something on the way out...
	log.Printf("Gophernaut is needed elsewhere...\n")
}
