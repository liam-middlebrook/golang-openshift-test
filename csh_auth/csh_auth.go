package csh_auth

import (
    "fmt"
    "os"
    "net/http"
    "time"
    log "github.com/sirupsen/logrus"
    jwt "github.com/dgrijalva/jwt-go"
    oidc "github.com/coreos/go-oidc"
    "golang.org/x/net/context"
    "golang.org/x/oauth2"
    "github.com/gin-gonic/gin"
)


var clientID = os.Getenv("csh_auth_client_id")
var clientSecret = os.Getenv("csh_auth_client_secret")
var secret = os.Getenv("csh_auth_jwt_secret")
var state = os.Getenv("csh_auth_state")
var server_host = os.Getenv("csh_auth_server_host")
var redirect_uri = os.getenv("csh_auth_redirect_uri")

var AuthKey = "cshauth"
var provider_uri = "https://sso.csh.rit.edu/realms/csh"

// globals eww
var config oauth2.Config // this guy changes a bit, weird
var ctx context.Context
var provider *oidc.Provider

// =================
//      structs
// =================

type CSHClaims struct {
    Token string `json:"token"`
    UserInfo CSHUserInfo `"json:user_info"`
    jwt.StandardClaims
}

type CSHUserInfo struct {
    Subject       string `json:"sub"`
    Profile       string `json:"profile"`
    Email         string `json:"email"`
    EmailVerified bool   `json:"email_verified"`
    // contains filtered or unexported fields
    Username       string `json:"preferred_username"`
    FullName       string `json:"name"`
}

// =================
//    auth helper
// =================

func AuthWrapper(page gin.HandlerFunc) gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {

        cookie, err := c.Cookie("Auth")
        if err != nil || cookie == "" {
            log.Info("cookie not found")
            c.Redirect(http.StatusFound, "/authenticate?referer=" + c.Request.URL.String())
            return
        }

        token, err := jwt.ParseWithClaims(cookie, &CSHClaims{}, func(token *jwt.Token) (interface{}, error){
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method")
            }
            return []byte(secret), nil
        })
        if err != nil {
            log.Fatal("token failure")
            return
        }

        if claims, ok := token.Claims.(*CSHClaims); ok && token.Valid {
            // add in user info data
            c.Set(AuthKey, *claims)
            // call the wrapped func
            page(c)
        } else {
            log.Fatal("error on auth layer?")
            return
        }
    })
}

func AuthRequest(c *gin.Context) {
    // Thrash this so we don't get additive weirdness
    config.RedirectURL = redirect_uri + "?referer=" + c.Query("referer")
    c.Redirect(http.StatusFound, config.AuthCodeURL(state))
}

func AuthCallback(c *gin.Context) {
    if c.Query("state") != state {
        log.Fatal("error state not match\n")
        return
    }
    oauth2Token, err := config.Exchange(ctx, c.Query("code"))
    if err != nil {
        log.Fatal("failed to exchange token\n")
        return
    }
    userInfo:= &CSHUserInfo{}
    oidcUserInfo, err := provider.UserInfo(ctx, oauth2.StaticTokenSource(oauth2Token))
    if err != nil {
        log.Fatal("failed to get userinfo\n")
    }
    oidcUserInfo.Claims(userInfo)
    if err != nil {
        log.Fatal("failed to get userinfo marshal\n")
    }

    expireToken := time.Now().Add(time.Hour * 1).Unix()
    expireCookie := 0 // session
    claims := CSHClaims {
        oauth2Token.AccessToken,
        *userInfo,
        jwt.StandardClaims {
            ExpiresAt: expireToken,
            Issuer: server_host,
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signedToken, err := token.SignedString([]byte(secret))

    c.SetCookie("Auth", signedToken, int(expireCookie), "", "", false, true)
    c.Redirect(http.StatusFound, c.Query("referer"))
}

func Init() {
    var err error
    ctx = context.Background()
    provider, err = oidc.NewProvider(ctx, provider_uri)
    if err != nil {
        log.Fatal(err)
    }
    config = oauth2.Config{
        ClientID: clientID,
        ClientSecret: clientSecret,
        Endpoint: provider.Endpoint(),
        RedirectURL: redirect_uri,
        Scopes: []string{oidc.ScopeOpenID, "profile", "email", "preferred_username", "name"},
    }
}

func AuthLogout(c *gin.Context) {
    c.SetCookie("Auth", "", 0, "", "", false, true)
    c.Redirect(http.StatusFound, provider_uri + "/protocol/openid-connect/logout?redirect_uri=http://" + server_host + "/")
}
