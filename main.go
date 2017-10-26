package main

import (
    "log"
    "wxbot2/wx"
    //"github.com/tonnerre/golang-pretty"
    "net/http"
    "encoding/json"
    _ "net/http/pprof"
    //"fmt"
)

var bot *wx.Wxbot

func main() {
    bot = wx.NewWxbot("hello", 50)
    err := bot.Login()
    if err != nil {
        log.Fatal(err)
    }
    http.HandleFunc("/send",Send)
    http.HandleFunc("/getFriends",GetFriends)
    http.ListenAndServe("127.0.0.1:9122",nil)
}

func Send(w http.ResponseWriter,r *http.Request){
    bot.SendMsg("没事就逛逛","test")
}

func GetFriends(w http.ResponseWriter,r *http.Request)  {
	//w.Header("")
	//fmt.Println(bot)
    b,_ := json.Marshal(bot.GetFriends())
    w.Write(b)
}


