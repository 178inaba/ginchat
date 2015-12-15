package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"golang.org/x/net/websocket"

	"github.com/178inaba/ginapp/model"
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func main() {
	r := gin.Default()
	r.Use(static.Serve("/", static.LocalFile("static", false)))
	r.LoadHTMLGlob("tpl/*")

	r.GET("/", root)

	r.GET("/join", join)
	r.POST("/join", postJoin)

	r.GET("/login", login)
	r.POST("/login", postLogin)

	r.GET("/chat", chat)

	// websocket
	r.GET("/ws_chat", wsChat)

	r.NoRoute(notFound)

	r.Run(":8080")
}

func root(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func join(c *gin.Context) {
	c.HTML(http.StatusOK, "join.html", nil)
}

func postJoin(c *gin.Context) {
	screenName := c.PostForm("id")
	pass := c.PostForm("pass")
	name := c.PostForm("name")
	strAge := c.PostForm("age")
	intro := c.PostForm("intro")

	// pass to crypto
	cryptoPass := fmt.Sprintf("%x", sha256.Sum256([]byte(pass)))

	// age to int
	age, err := strconv.Atoi(strAge)
	if err != nil {
		log.Error("conv err: ", err)
		c.HTML(http.StatusOK, "join.html", nil)
		return
	}

	model.CreateUser(screenName, cryptoPass, name, intro, age)

	c.Redirect(http.StatusFound, "/chat")
}

func login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func postLogin(c *gin.Context) {
	screenName := c.PostForm("id")
	pass := c.PostForm("pass")

	user := model.GetUser(screenName)

	// pass to crypto
	cryptoPass := fmt.Sprintf("%x", sha256.Sum256([]byte(pass)))

	if user.CryptoPass != cryptoPass {
		log.Error("login err: screen name: ", screenName)
		c.HTML(http.StatusOK, "login.html", nil)
		return
	}

	c.Redirect(http.StatusFound, "/chat")
}

func chat(c *gin.Context) {
	c.HTML(http.StatusOK, "chat.html", nil)
}

func wsChat(c *gin.Context) {
	h := websocket.Handler(chatHandler)
	h.ServeHTTP(c.Writer, c.Request)
}

func notFound(c *gin.Context) {
	c.HTML(http.StatusNotFound, "404.html", nil)
}

// websocket
func chatHandler(conn *websocket.Conn) {
	io.Copy(conn, conn)
}
