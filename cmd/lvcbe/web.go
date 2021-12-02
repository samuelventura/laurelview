package main

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/samuelventura/go-tools"
	"github.com/samuelventura/go-tree"
)

func web(node tree.Node) {
	api := node.GetValue("api").(*apiDso)
	endpoint := node.GetValue("webep").(string)
	gin.SetMode(gin.ReleaseMode) //remove debug warning
	router := gin.New()          //remove default logger
	router.Use(gin.Recovery())   //catches panics
	router.Use(static)
	rapi := router.Group("/api")
	rapi.POST("/signup", func(c *gin.Context) {
		aid, ok := c.GetQuery("aid")
		if !ok {
			c.String(http.StatusBadRequest, "missing aid")
			return
		}
		if !strings.Contains(aid, "@") {
			c.String(http.StatusBadRequest, "invalid id format "+aid)
			return
		}
		_, err := api.post_signup(aid)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		//FIXME send password thru email
		msg := "Your password was sent to your email."
		c.JSON(http.StatusOK, msg)
	})
	rapi.POST("/signin", func(c *gin.Context) {
		aid, ok := c.GetQuery("aid")
		if !ok {
			c.String(http.StatusBadRequest, "missing aid")
			return
		}
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		dro, err := api.post_signin(aid, string(body))
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		c.JSON(http.StatusOK, dro)
	})
	rapi.GET("/signout", func(c *gin.Context) {
		sid, ok := c.GetQuery("sid")
		if !ok {
			c.String(http.StatusBadRequest, "missing sid")
			return
		}
		dro, err := api.get_signout(sid)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		c.JSON(http.StatusOK, dro)
	})
	rapi.POST("/recover", func(c *gin.Context) {
		aid, ok := c.GetQuery("aid")
		if !ok {
			c.String(http.StatusBadRequest, "missing aid")
			return
		}
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		_, err = api.post_recover(aid, string(body))
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		//FIXME send password thru email
		msg := "A recovering password was sent to your email."
		c.JSON(http.StatusOK, msg)
	})
	listen, err := net.Listen("tcp", endpoint)
	tools.PanicIfError(err)
	node.AddCloser("listen", listen.Close)
	port := listen.Addr().(*net.TCPAddr).Port
	node.SetValue("port", port)
	log.Println("port", port)
	server := &http.Server{
		Addr:    endpoint,
		Handler: router,
	}
	node.AddProcess("server", func() {
		err = server.Serve(listen)
		if err != nil {
			log.Println(endpoint, port, err)
		}
	})
}
