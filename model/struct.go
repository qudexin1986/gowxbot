package model

import (
	"time"
)

type Msg struct {
	Id         int       `xorm:"not null pk autoincr INT(11)"`
	Data       string    `xorm:"TEXT"`
	CreateTime time.Time `xorm:"TIMESTAMP"`
	Status     int       `xorm:"INT(11)"`
	Nickname   string    `xorm:"TEXT"`
	UpdateTime time.Time `xorm:"DATETIME"`
}
