package wxbot

import (
	"wxbot2/response"
	"encoding/json"
	"reflect"
	//"fmt"
)

type MsgHander struct {

	Response func(in response.Msg, session *Wxbot)
	MsgRec  []byte
	Wxbot *Wxbot
	Uid int
}

func (m *MsgHander)Init()  {
	SpecialContact := map[string]bool{
		"filehelper":            true,
		"newsapp":               true,
		"fmessage":              true,
		"weibo":                 true,
		"qqmail":                true,
		"tmessage":              true,
		"qmessage":              true,
		"qqsync":                true,
		"floatbottle":           true,
		"lbsapp":                true,
		"shakeapp":              true,
		"medianote":             true,
		"qqfriend":              true,
		"readerapp":             true,
		"blogapp":               true,
		"facebookapp":           true,
		"masssendapp":           true,
		"meishiapp":             true,
		"feedsapp":              true,
		"voip":                  true,
		"blogappweixin":         true,
		"weixin":                true,
		"brandsessionholder":    true,
		"weixinreminder":        true,
		"officialaccounts":      true,
		"wxitil":                true,
		"userexperience_alarm":  true,
		"notification_messages": true,
	}
	dta := new(response.Webwxsync)
	json.Unmarshal(m.MsgRec, dta)
	if dta.AddMsgCount > 0{
		for _,ms := range dta.AddMsgList{
			if SpecialContact[ms.FromUserName]{
				continue
			}
			if ms.FromUserName == m.Wxbot.User.UserName{
				continue
			}
			if isEmpty(ms){
				continue
			}
			//fmt.Println(m.Wxbot)
		 	m.Response(ms,m.Wxbot)
			//m.Wxbot.SendText(ms.FromUserName,to )
		}
	}
}

func isEmpty(a interface{}) bool {
	v := reflect.ValueOf(a)
	if v.Kind() == reflect.Ptr {
		v=v.Elem()
	}
	return v.Interface() == reflect.Zero(v.Type()).Interface()
}