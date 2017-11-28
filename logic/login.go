package logic

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"time"
	"encoding/json"
	"wxbot2/model"
	"fmt"
)



func Login(c *gin.Context){
	//fmt.Println(c.Request.Body)\
	user:= new(model.User)
	c.BindJSON(&user)
	fmt.Println(user)
	db,_ := c.Get("db")
	dbc := db.(*gorm.DB)

	dbc.Debug().Where(user).Find(&user)
	var ret = make(map[string]interface{})

	if user.Id == 0{
		c.Writer.WriteHeader(500)
		return
		//c.Writer.Header().Set("status","500")
	}
	c.Writer.WriteHeader(200)
	ret["id"] = user.Id
	b,_ := json.Marshal(ret)
	c.Writer.Write(b)
}

func Register( c *gin.Context){
	name := c.PostForm("name")
	password := c.PostForm("password")
	email := c.PostForm("email")
	//var ret = make(map[string]interface{})
	//ret["msg"] = ""

	if name == "" ||  password == "" || email == ""{
		c.Writer.Header().Set("status","500")
		return
	}
	user:= new(model.User)
	//Mon Jan 2 15:04:05
	user.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	user.Name = name
	user.Password = password
	user.Email = email
	user.Status = 1
	user.UpdateTime = user.CreateTime
	db,_ := c.Get("db")
	dbc := db.(*gorm.DB)
	if dbc.Save(&user).Error != nil{
		c.Writer.WriteHeader(500)
		//c.Writer.Header().Set("status","500")
		return
	}
	c.Writer.WriteHeader(200)
	//c.Writer.Header().Set("status","200")
}


