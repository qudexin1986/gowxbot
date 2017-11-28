package logic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"wxbot2/manager"
	"net/http"
)


func Status(c *gin.Context) {
	uuid := c.Query("uuid")
	session := manager.GlobalSessionManager.Get(uuid)
	if session == nil {
		c.String(http.StatusNotFound, fmt.Sprintf("no session for %s", uuid))
		return
	}
	if session.Cookies == nil {
		c.JSON(200, gin.H{
			"status":    "CREATED",
		})
		return
	}



	c.JSON(200, gin.H{
		"status":    "SERVING",
		"bot" :session,
	})
	return
}
