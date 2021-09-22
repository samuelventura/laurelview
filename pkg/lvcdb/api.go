package lvcdb

import (
	"context"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

// dbDriver = sqlite3
// dbSource = :memory:
// curl --header "token: abc123" -D - http://localhost:5001/api/login
// curl --header "token: abc123" -D - http://localhost:5001/api/auth/nodes
// curl --header "token: abc123" -D - http://localhost:5001/api/auth/nodes/mac123

func NewApi(ctx Context) Dispatch {
	log := ctx.PrefixLog("api")
	driver := ctx.GetValue("dbDriver").(string)
	source := ctx.GetValue("dbSource").(string)
	endpoint := ctx.GetValue("apiEndpoint").(string)
	dao := NewDao(driver, source)
	gin.SetMode(gin.ReleaseMode) //remove debug warning
	router := gin.New()          //remove default logger
	router.Use(gin.Recovery())   //looks important
	router.GET("/api/login", func(c *gin.Context) {
		key := http.CanonicalHeaderKey("token")
		token := c.Request.Header[key]
		log.Info("/api/login", token)
		if len(token) == 1 {
			c.Header("token", token[0])
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	})
	authorized := router.Group("/api/auth", func(c *gin.Context) {
		key := http.CanonicalHeaderKey("token")
		token := c.Request.Header[key]
		log.Info("/api/auth", token)
		if len(token) != 1 {
			c.Header("WWW-Authenticate", "Authorization Required")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set("token", token[0])
	})
	authorized.GET("/nodes", func(c *gin.Context) {
		token := c.MustGet("token").(string)
		list := dao.ListNodes(token)
		c.JSON(http.StatusOK, gin.H{"nodes": list})
	})
	authorized.GET("/nodes/:mac", func(c *gin.Context) {
		token := c.MustGet("token").(string)
		mac := c.Param("mac")
		node := dao.GetNode(token, mac)
		c.JSON(http.StatusOK, gin.H{"node": node})
	})
	server := &http.Server{
		Addr:    endpoint,
		Handler: router,
	}
	go func() {
		defer TraceRecover(log.Debug)
		listen, err := net.Listen("tcp", endpoint)
		PanicIfError(err)
		defer listen.Close()
		port := listen.Addr().(*net.TCPAddr).Port
		log.Info("port", port)
		err = server.Serve(listen)
		TraceIfError(log.Debug, err)
	}()
	dispatchs := make(map[string]Dispatch)
	dispatchs["dispose"] = func(mut Mutation) {
		ClearDispatch(dispatchs)
		server.Shutdown(context.Background())
		DisposeArgs(mut.Args)
	}
	return MapDispatch(log, dispatchs)
}
