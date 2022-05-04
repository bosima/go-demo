package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	httpServer "github.com/go-micro/plugins/v4/server/http"
	"go-micro.dev/v4"
	"go-micro.dev/v4/server"
)

type Hello struct {
	Language string
}

func (a *Hello) InitRouter(router *gin.Engine) {
	router.GET("/hello/say", a.Say)
	router.POST("/hello/set", a.Set)
}

func (h *Hello) getHelloString() string {
	if h.Language == "cn" {
		return "你好"
	} else {
		return "Hello"
	}
}

func (h *Hello) Say(c *gin.Context) {
	name, ok := c.GetQuery("name")

	if !ok {
		c.String(http.StatusBadRequest, "no query parameter 'name'")
		return
	}

	c.JSON(200, gin.H{
		"Message": fmt.Sprintf("%s %s", h.getHelloString(), name),
	})
}

func (h *Hello) Set(c *gin.Context) {
	var json map[string]interface{}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if lang, ok := json["Language"]; ok {
		h.Language = lang.(string)
		c.JSON(200, gin.H{
			"errcode": 0,
		})
		return
	} else {
		c.String(http.StatusBadRequest, "not found 'Language'")
		return
	}
}

func main() {
	srv := httpServer.NewServer(
		server.Name("go-micro-gin-demo"),
		server.Address(":8002"),
	)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	hello := &Hello{}
	hello.InitRouter(router)

	hd := srv.NewHandler(router)
	if err := srv.Handle(hd); err != nil {
		log.Fatalln(err)
	}

	service := micro.NewService(
		micro.Server(srv),
	)
	service.Init()
	service.Run()
}
