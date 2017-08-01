package main

import (
    "flag"
    "fmt"
    "github.com/golang/glog"
    "net/http"
    "runtime"
)

func handler(w http.ResponseWriter, r *http.Request) {
    glog.Infof("Received request from %s for %s", r.RemoteAddr, r.URL)
    glog.Flush()

    h := w.Header()
    h.Set("Content-Type", "text/plain")

    fmt.Fprint(w, "Hello world!\n\n")
    fmt.Fprintf(w, "Go version: %s\n", runtime.Version())
}

func main() {
    flag.Parse()

    glog.Info("Starting server...")
    glog.Flush()
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}
