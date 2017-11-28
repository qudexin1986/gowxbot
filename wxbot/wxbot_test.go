package wxbot

import (
	"testing"
	"fmt"
	"time"
	"strconv"
	//"github.com/gin-gonic/gin/binding/example"
)

func TestJsLogin(t *testing.T){
	td := new(Wxbot)
	fmt.Println(td.JsLogin())
}

func TestQrcode(t *testing.T){
	td := new(Wxbot)
	fmt.Println(td.GetQrCode())
}

func TestMile(t *testing.T) {
	fmt.Println(strconv.FormatInt(time.Now().UnixNano()/1000000,10))
	fmt.Println(time.Now().UnixNano()/1000000)
}

func TestFmt(t *testing.T) {
	type Student struct {
		Name string
	}

	var s *Student = new(Student)

	s.Name = "jack"

	td := *s

	fmt.Println("t=", td, "s=", s)

	s.Name = "rose"

	fmt.Println("t=", td, "s=", s)
}