package logic

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wxbot2/wxbot"
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"encoding/json"
	"time"
	"wxbot2/model"
	"wxbot2/response"
	"wxbot2/manager"
)

func Create(c *gin.Context) {
	// create session
	bot := wxbot.Init(nil)
	uid := c.Query("uid")
	db, _ := c.Get("db")
	dbc := db.(*gorm.DB)
	useid, _ := strconv.Atoi(uid)
	bot.InitOver = func() {
		config, _ := json.Marshal(bot.Data)
		nt := time.Now().Format("2006-01-02 15:04:05")
		wxb := new(model.Wxbot)
		wxb.Uid = useid
		wxb.Uin = bot.XmlConfig.Wxuin
		wxb.Status = 1
		wxb.UpdateTime = nt
		wxb.CreateTime = nt
		wxb.Config = string(config)
		wxb.Uuid = bot.Uuid
		if dbc.Save(wxb).Error != nil{
			dbc.Debug().Table("wxbot").Where("uin=?",wxb.Uin).Update(wxb)
		}
		manager.GlobalSessionManager.Set(bot.Uuid,bot)
	}
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

	go bot.Run(false)

	c.String(http.StatusOK, bot.Uuid)
}
