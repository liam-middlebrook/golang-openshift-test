package main

import (
    "flag"
    "fmt"
    "net/http"
    "runtime"
    log "github.com/sirupsen/logrus"
)

import "github.com/gin-gonic/gin"

func handler(w http.ResponseWriter, r *http.Request) {
    log.WithFields(log.Fields{
    "remote": r.RemoteAddr,
    "url": r.URL,
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

    r := gin.Default()

    r.GET("/", func (c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "MUNGE",
        })
    })

    r.Run()
}
