package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
)

func start() {
	router := gin.Default()
	router.NoRoute(func(c *gin.Context) {
		c.IndentedJSON(http.StatusNotFound, nil)
	})

	router.GET("/welcome", func(c *gin.Context) {
		firstname := c.DefaultQuery("firstname", "Guest")
		lastname := c.Query("lastname")

		c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
	})

	router.GET("/shou/:firstname/:lastname", func(c *gin.Context) {
		//firstname := c.DefaultQuery("firstname", "Guest")
		//lastname := c.Query("lastname")
		firstname := c.Param("firstname")
		lastname := c.Param("lastname")

		c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
	})

	router.POST("/print", func(c *gin.Context) {
		fmt.Println(c.Request.URL.Path)
		fmt.Println(c.Request.ContentLength, c.Request.Form)
		err := c.Request.ParseForm()
		if err != nil {
			fmt.Println(err)
		}
		for k, v := range c.Request.PostForm {
			fmt.Printf("k:%v\n", k)
			fmt.Printf("v:%v\n", v)
		}
		c.String(http.StatusOK, "Hello print")
	})
	router.Run()
}

func main() {
	//router := gin.Default()
	//router.GET("/welcome", func(c *gin.Context) {
	//	firstname := c.DefaultQuery("firstname", "Guest")
	//	lastname := c.Query("lastname")
	//
	//	c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
	//})
	//
	//router.GET("/shou/:firstname/:lastname", func(c *gin.Context) {
	//	//firstname := c.DefaultQuery("firstname", "Guest")
	//	//lastname := c.Query("lastname")
	//	firstname := c.Param("firstname")
	//	lastname := c.Param("lastname")
	//
	//	c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
	//})
	//
	//router.Run()
	start()
}