package main

import (
    "github.com/gin-gonic/gin"
    "wxbot2/logic"
    "net/http"
    "github.com/jinzhu/gorm"
    _ "github.com/go-sql-driver/mysql"

)

func main() {
    router := gin.Default()
    db,err := gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/wxmsg?charset=utf8")
    //db.LogMode(true)
    if err != nil {
        panic(err)
    }
    router.Use(func(c *gin.Context){
        c.Set("db",db)
    })


    router.Use(gin.HandlerFunc(func(c *gin.Context) {

        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        if c.Request.Method == "OPTIONS" {
            c.Writer.Header().Set("Access-Control-Allow-Headers","Origin, X-Requested-With, Content-Type, Accept")
            c.AbortWithStatus(http.StatusOK)
            return
        }
        c.Next()
    }))


    router.GET("/create", logic.Create)
    router.GET("/status", logic.Status)
    router.GET("/start/:id",logic.Start)
    router.POST("/login",logic.Login)
    router.POST("/register",logic.Register)
    router.POST("/friends",logic.Friends)
    router.POST("/setplug",logic.SetPlug)

    router.Run(":8080")
}
