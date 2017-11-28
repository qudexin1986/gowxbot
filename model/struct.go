package model


type Msg struct {
	//gorm.Model
	Id int
	Data string
	CreateTime string
	Status int
	Nickname string
	UpdateTime string
	Uin	int
}

func (m Msg) TableName() string{
	return "msg"
}


type User struct {
	Id int
	Name string
	Password string
	Status int
	Email string
	CreateTime string
	UpdateTime string
}


func (m User) TableName() string{
	return "user"
}


type Wxbot struct {
	Id int
	Uid int
	Uin string
	Config string
	CreateTime string
	UpdateTime string
	Status int
	Uuid string
}


func (w Wxbot) TableName() string{
	return "wxbot"
}



type TextMessage struct {
	Type         int
	Content      string
	FromUserName string
	ToUserName   string
	LocalID      int
	ClientMsgId  int
}

