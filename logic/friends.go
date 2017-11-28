package logic

import (
	"github.com/gin-gonic/gin"
	"wxbot2/manager"
	"encoding/json"
	"fmt"
)

func Friends(c *gin.Context){
	uuid := c.PostForm("uuid")
	session :=  manager.GlobalSessionManager.Get(uuid)
	fmt.Println(manager.GlobalSessionManager)
	b,_ := json.Marshal(session.UserList)
	c.Writer.Write(b)
}
