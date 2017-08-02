package main

import (
    "flag"
    "net/http"
    log "github.com/sirupsen/logrus"
    "github.com/gin-gonic/gin"
    csh_auth "github.com/liam-middlebrook/csh-auth"
)

func protectedProfile(c *gin.Context){
    claims, ok := c.Value(csh_auth.AuthKey).(csh_auth.CSHClaims)
    if !ok {
        log.Fatal("error finding claims")
        return
    }
    c.String(http.StatusOK, "uid %s email %s name %s uuid %s", claims.UserInfo.Username, claims.UserInfo.Email, claims.UserInfo.FullName, claims.UserInfo.Subject)
}

func main() {
    flag.Parse()

    log.Info("Starting server...")

    // needs to be declared here not inline so provider is global XXX FIXME
    r := gin.Default()
    csh_auth.Init()
    r.GET("/", csh_auth.AuthWrapper(func (c *gin.Context) { c.String(http.StatusOK, "spooky data") }))
    r.GET("/test", csh_auth.AuthWrapper(protectedProfile))
    r.GET("/authenticate", csh_auth.AuthRequest)
    r.GET("/authredir", csh_auth.AuthCallback)
    r.GET("/logout", csh_auth.AuthLogout)
    r.Run()
}
