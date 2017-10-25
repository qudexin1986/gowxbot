package response

import (
	"strconv"
	"strings"
)

type Msg struct {

	MsgId string
	FromUserName string
	ToUserName string
	MsgType int64
	Content string
	Status int64
	ImgStatus  int
	CreateTime  int64
	VoiceLength  int
	PlayLength int64
	FileName string
	FileSize string
	MediaId string
	Url string
	AppMsgType int
	StatusNotifyCode int
	StatusNotifyUserName string
	RecommendInfo RecommendInfo
	ForwardFlag int
	AppInfo AppInfo
	HasProductId int
	Ticket string
	ImgHeight int64
	ImgWidth int64
	SubMsgType int
	NewMsgId  int64
	OriContent string
}

type AppInfo struct {
	AppID string
	Type int
}

type RecommendInfo struct {
	UserName string
	NickName string
	QQNum int64
	Province string
	City string
	Content string
	Signature string
	Alias string
	Scene int
	VerifyFlag int
	AttrStatus int
	Sex int
	Ticket string
	OpCode int
}

type Webwxsync struct {
	BaseResponse BaseResponse
	AddMsgCount int64
	AddMsgList []Msg
	ModContactCount int64
	ModContactList []interface{}
	DelContactCount int64
	DelContactList []interface{}
	ModChatRoomMemberCount int64
	ModChatRoomMemberList []interface{}
	Profile Profile
	ContinueFlag int64
	SyncKey SyncKey
	SKey string
	SyncCheckKey SyncKey
}

type BaseResponse struct {
	Ret int64
	ErrMsg string
}


type Profile struct {

	BitFlag  int64
	UserName UserName
	NickName NickName
	BindUin int64
	BindEmail BindEmail
	BindMobile BindMobile
	Status int64
	Sex int
	PersonalCard int64
	Alias string
	HeadImgUpdateFlag int64
	HeadImgUrl string
	Signature string
}

type UserName struct {
	Buff string
}

type NickName struct {
	Buff string
}
type BindEmail struct {
	Buff string
}
type BindMobile struct {
	Buff string
}


type SyncKey struct {
	Count int
	List []KV
}

func (s *SyncKey)Encode() string{
	var kvList []string
	for _,v := range s.List {
		kvList = append(kvList,strconv.Itoa(v.Key)+"_"+strconv.Itoa(v.Val))
	}
	return strings.Join(kvList,"|")
}

type KV struct {
	Key int
	Val int
}

type Member struct {

}

type Contact struct {
	Uin int64
	UserName string
	NickName string
	HeadImgUrl string
	ContactFlag int64
	RemarkName string
	HideInputBarFlag int64
	Sex int
	Signature string
	VerifyFlag int64
	OwnerUin int64
	PYInitial string
	PYQuanPin string
	RemarkPYInitial string
	RemarkPYQuanPin string
	StarFriend int64
	AppAccountFlag int64
	Statues int64
	AttrStatus int64
	Province string
	City string
	Alias string
	SnsFlag int64
	UniFriend int64
	DisplayName string
	ChatRoomId int64
	KeyWord string
	EncryChatRoomId string
	IsOwner int
	MemberCount int64
	MemberList []Member
}

type User struct {
	Uin int64
	UserName string
	NickName string
	HeadImgUrl string
	RemarkName string
	PYInitial string
	PYQuanPin string
	RemarkPYInitial string
	RemarkPYQuanPin string
	HideInputBarFlag int64
	StarFriend int64
	Sex int
	Signature string
	AppAccountFlag int64
	VerifyFlag int64
	ContactFlag int64
	WebWxPluginSwitch int64
	HeadImgFlag int64
	SnsFlag int64
}

type MPSubscribeMsg struct {

}

type WebInit struct {
	BaseResponse BaseResponse
	Count int64
	ContactList []Contact
	SyncKey SyncKey

	User User
	ChatSet string
	SKey string
	ClientVersion int64
	SystemTime int64
	GrayScale int64
	InviteStartCount int64
	MPSubscribeMsgCount int64
	MPSubscribeMsgList []interface{}
	ClickReportInterval int64
}

