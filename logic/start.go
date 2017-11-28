package logic

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"wxbot2/model"
	//"fmt"
	"encoding/json"

	"net/http"

	"wxbot2/wxbot"
	"time"
	"wxbot2/response"
	"wxbot2/manager"
	"wxbot2/plug"
)

func Start(c *gin.Context){
	id := c.Param("id")
	db,_ := c.Get("db")
	dbc := db.(*gorm.DB)
	user := new(model.Wxbot)
	dbc.Where("uid=? and status = 1",id).Find(user)
	//fmt.Println(users)
	var data = new(wxbot.Data)
	if manager.GlobalSessionManager.Get(user.Uuid) != nil{

		return
	}
	json.Unmarshal([]byte(user.Config),data)
	bot := wxbot.Init(data)
	bot.InitOver = func() {
		config, _ := json.Marshal(bot.Data)
		nt := time.Now().Format("2006-01-02 15:04:05")
		wxb := new(model.Wxbot)
		wxb.UpdateTime = nt
		wxb.Config = string(config)
		dbc.Debug().Table("wxbot").Where("uin=?",bot.XmlConfig.Wxuin).Updates(wxb)
		manager.GlobalSessionManager.Set(bot.Uuid,bot)
	}

	bot.MsgHander.Uid = user.Uid
	bot.MsgHander.Response = plug.Weather // plug

	bot.Handler = func( msg []byte) {
		config, _ := json.Marshal(bot.Data)
		nt := time.Now().Format("2006-01-02 15:04:05")
		wxb := new(model.Wxbot)
		wxb.UpdateTime = nt
		wxb.Config = string(config)
		dbc.Debug().Table("wxbot").Where("uin=?",bot.XmlConfig.Wxuin).Updates(wxb)

		dta := new(response.Webwxsync)
		json.Unmarshal(msg, dta)
		if dta.AddMsgCount > 0{
			for _,ms := range dta.AddMsgList{
				if bot.SpecialContact[ms.FromUserName]{
					continue
				}
				if ms.FromUserName == bot.User.UserName{
					continue
				}
				bot.SendText(ms.FromUserName,ms.Content)
			}
		}
	}
	go bot.Run(true)

	c.String(http.StatusOK, user.Uuid)
}
