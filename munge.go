package main

import (
    "flag"
    "fmt"
    "net/http"
    "runtime"
    log "github.com/sirupsen/logrus"
)

func handler(w http.ResponseWriter, r *http.Request) {
    log.WithFields(log.Fields{
    "remote": r.RemoteAddr,
    "url": r.URL
    }).Info("Received request")

    h := w.Header()
    h.Set("Content-Type", "text/plain")

    fmt.Fprint(w, "Hello world!\n\n")
    fmt.Fprintf(w, "Go version: %s\n", runtime.Version())
    fmt.Fprint(w, "MUNGE!\n\n")
}

func main() {
    flag.Parse()

    log.Info("Starting server...")
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}
