package main

import (
    "log"
    "oa.com/melow_dog/wx"
    //"github.com/tonnerre/golang-pretty"
)

func main() {
    bot := wx.NewWxbot("hello", 50)
    err := bot.Login()
    if err != nil {
        log.Fatal(err)
    }
}
