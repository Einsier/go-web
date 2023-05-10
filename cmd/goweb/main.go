package main

import (
	"log"
	"net/http"
	"time"

	"github.com/einsier/go-web/web"
)

func onlyForV2() web.HandlerFunc {
	return func(c *web.Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		c.Fail(500, "Internal Server Error")
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	r := web.New()
	r.Use(web.Logger(), web.Recovery()) // global middleware

	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./static")

	r.GET("/", func(c *web.Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})

	// index out of range for testing Recovery()
	r.GET("/panic", func(c *web.Context) {
		names := []string{"webktutu"}
		c.String(http.StatusOK, names[100])
	})

	v1 := r.Group("/v1")
	{
		v1.GET("/hello", func(c *web.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}

	v2 := r.Group("/v2")
	v2.Use(onlyForV2()) // v2 group middleware
	{
		v2.GET("/hello/:name", func(c *web.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
		v2.POST("/login", func(c *web.Context) {
			c.JSON(http.StatusOK, web.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})

	}

	r.Run(":9999")
}
