package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
	"os"
	"io"
)

func Dummy1Middleware(c *gin.Context) {

	fmt.Println("Im a dummy1111!")

	// Pass on to the next-in-chain

	c.Next()
	fmt.Println("exit dummy 1")

}

func Dummy2Middleware(c *gin.Context) {

	fmt.Println("Im a dummy2222!")

	// Pass on to the next-in-chain

	c.Next()
	fmt.Println("exit dummy 2")
}


func Dummy3Middleware() gin.HandlerFunc {

	fmt.Println("Im a dummy333!")

	// Pass on to the next-in-chain

	return func(c *gin.Context) {

		c.Next()

	}

}

func Dummy4Middleware() gin.HandlerFunc {

	fmt.Println("Im a dummy444!")

	// Pass on to the next-in-chain

	return func(c *gin.Context) {

		c.Next()

	}

}

func start() {
	router := gin.Default()
	router.NoRoute(func(c *gin.Context) {
		c.IndentedJSON(http.StatusNotFound, nil)
	})

	router.Use(Dummy1Middleware)
	router.Use(Dummy2Middleware)
	//router.Use(Dummy3Middleware())
	//router.Use(Dummy4Middleware())
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
		//err := c.Request.ParseForm()
		//if err != nil {
		//	fmt.Println(err)
		//}
		//for k, v := range c.Request.PostForm {
		//	fmt.Printf("k:%v\n", k)
		//	fmt.Printf("v:%v\n", v)
		//}

		file, handler, _ := c.Request.FormFile("file")
		filename := handler.Filename
		fmt.Println("Received file:", filename)

		out, err := os.Create("/home/qydev/haha")
		if err != nil {
			fmt.Println("error: ", err)
			c.String(http.StatusExpectationFailed, "error open")
			return
		}
		defer out.Close()
		_, err = io.Copy(out, file)
		if err != nil {
			fmt.Println("error: ", err)
			c.String(http.StatusExpectationFailed, "error save")
			return
		}

		c.String(http.StatusOK, "Hello print")
	})

	router.POST("/multi", func(c *gin.Context) {
		form, _ := c.MultipartForm()
		files := form.File["files"]
		for _, f := range files {
			fmt.Println(f.Filename)
			file, e := f.Open()
			if e == nil {
				filename := fmt.Sprint("/home/qydev/", f.Filename)
				fmt.Println(filename)
				out, err := os.Create(filename)
				if err != nil {
					fmt.Println("error: ", err)
					c.String(http.StatusExpectationFailed, "error open")
					return
				}
				defer out.Close()
				_, err = io.Copy(out, file)
				if err != nil {
					fmt.Println("error: ", err)
					c.String(http.StatusExpectationFailed, "error save")
					return
				}
			} else {
				fmt.Println("eeeeeeee happen", e)
				c.String(http.StatusExpectationFailed, "error get file")
				return
			}
		}
		c.String(http.StatusOK, "Hello multi")
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