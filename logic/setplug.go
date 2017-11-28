package logic

import (
	"github.com/gin-gonic/gin"
	"wxbot2/manager"
	"fmt"
	"wxbot2/plug"
	"wxbot2/response"
	"wxbot2/wxbot"
)

func SetPlug(c *gin.Context){

	uuid := c.PostForm("uuid")
	q := c.PostForm("q")
	//uuid := qu.Uuid
	fmt.Println(uuid)
	fmt.Println(q)
	bot := manager.GlobalSessionManager.Get(uuid)
	fmt.Println(bot)
	switch q {
		case "weather":
			bot.MsgHander.Response = plug.Weather
		default:
			bot.MsgHander.Response = Ra1
	}

	c.Writer.WriteHeader(200)
}


func Ra1(a response.Msg, s *wxbot.Wxbot) {
	s.SendText(a.FromUserName,"ra1"+a.Content)
}
