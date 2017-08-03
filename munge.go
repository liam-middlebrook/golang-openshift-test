package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	csh_auth "github.com/liam-middlebrook/csh-auth"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func protectedProfile(c *gin.Context) {
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

	r := gin.Default()
	csh := csh_auth.CSHAuth{}
	log.Info("csh auth")
	csh.Init(
		os.Getenv("csh_auth_client_id"),
		os.Getenv("csh_auth_client_secret"),
		os.Getenv("csh_auth_jwt_secret"),
		os.Getenv("csh_auth_state"),
		os.Getenv("csh_auth_server_host"),
		os.Getenv("csh_auth_redirect_uri"),
		"/auth/login",
	)
	log.Info("csh auth INIT")
	r.GET("/", csh.AuthWrapper(func(c *gin.Context) { c.String(http.StatusOK, "spooky data") }))
	r.GET("/test", csh.AuthWrapper(protectedProfile))
	r.GET("/auth/login", csh.AuthRequest)
	r.GET("/auth/redir", csh.AuthCallback)
	r.GET("/auth/logout", csh.AuthLogout)
	r.Run()
}
