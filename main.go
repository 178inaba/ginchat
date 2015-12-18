package main

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"strconv"

	"golang.org/x/net/websocket"

	"github.com/178inaba/ginchat/model"
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/ugorji/go/codec"
)

var (
	mh = &codec.MsgpackHandle{}
	mp = websocket.Codec{Marshal: mpMarshal, Unmarshal: mpUnmarshal}
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
	var err error
	for {
		var data string
		err = mp.Receive(conn, &data)
		if err != nil {
			log.Error("receive error: ", err)
			break
		}

		log.Debug("data: ", data)

		err = mp.Send(conn, data)
		if err != nil {
			log.Error("send error: ", err)
		}
	}
}

func mpMarshal(v interface{}) (msg []byte, payloadType byte, err error) {
	err = codec.NewEncoderBytes(&msg, mh).Encode(v)
	return msg, websocket.BinaryFrame, err
}

func mpUnmarshal(msg []byte, payloadType byte, v interface{}) (err error) {
	return codec.NewDecoderBytes(msg, mh).Decode(v)
}
