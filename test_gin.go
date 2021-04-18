package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
)

func main()  {
	fmt.Println("hello ")
	r :=gin.Default()
	// get无参数传递
	r.GET("/hello/test", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// get有参数传递 该方式/hello/test 不可匹配
	r.GET("/hello/test/:id", func(context *gin.Context) {
		id :=context.Param("id")
		context.JSON(200,gin.H{
			"id":id,
		})
	})
   //路径中的id参数可容忍缺失 :/hello/default如果没有严格匹配的路由 会被重定向为/hello/default/ id的取值为/123
	r.GET("/hello/default/*id", func(context *gin.Context) {
		var  id  =context.Param("id")
		fmt.Println(id)
		if id == "" {
          context.JSON(400,gin.H{
          	"id":"",
          	"error":"id为null",
		  })
		} else {
			context.JSON(200,gin.H{
				"id":id,
			})
		}

	})
	//url parameter中取值
	r.GET("/test", func(context *gin.Context) {
		firstName :=context.DefaultQuery("firstname","he")
		lastName :=context.Query("lastname")
		context.JSON(200,gin.H{
			"name":firstName+lastName,
		})
	})

	// 以表单的形式读取数据
	r.POST("/hello/post", func(c *gin.Context) {
		firstName := c.PostForm("fistname")
		lastName := c.DefaultPostForm("lastname", "anonymous") // 此方法可以设置默认值
		c.JSON(200, gin.H{
			"status":  "success",
			"name":firstName+lastName,
		})
	})

	type PersonName struct{
		FirstName string
		LastName string
	}
    // 获取前端以json的形式的传值
	r.POST("/hello/post/json", func(c *gin.Context) {
		personName :=&PersonName{}
		if c.Bind(personName) ==nil{
			firstName := personName.FirstName
			lastName := personName.LastName // 此方法可以设置默认值
			c.JSON(200, gin.H{
				"status":  "success",
				"name":firstName+lastName,
			})
		}else {
			c.JSON(200, gin.H{
				"status":  "error",
				"message":nil,
			})
		}

	})

	r.Run("localhost:8081")


}
