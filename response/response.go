package response

import (
	"strconv"
	"strings"

)

type Msg struct {
	MsgId string
	FromUserName string
	ToUserName string
	MsgType int
	Content string
	Status int
	ImgStatus  int
	CreateTime  int
	VoiceLength  int
	PlayLength int
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
	ImgHeight int
	ImgWidth int
	SubMsgType int
	NewMsgId  int
	OriContent string
}

type AppInfo struct {
	AppID string
	Type int
}

type RecommendInfo struct {
	UserName string
	NickName string
	QQNum int
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
	AddMsgCount int
	AddMsgList []Msg
	ModContactCount int
	ModContactList []interface{}
	DelContactCount int
	DelContactList []interface{}
	ModChatRoomMemberCount int
	ModChatRoomMemberList []interface{}
	Profile Profile
	ContinueFlag int
	SyncKey SyncKey
	SKey string
	SyncCheckKey SyncKey
}

type BaseResponse struct {
	Ret int
	ErrMsg string
}
type BaseRequest struct {
	Uin     string
	Sid     string
	Skey    string
	DeviceID    string
}

type Profile struct {

	BitFlag  int
	UserName UserName
	NickName NickName
	BindUin int
	BindEmail BindEmail
	BindMobile BindMobile
	Status int
	Sex int
	PersonalCard int
	Alias string
	HeadImgUpdateFlag int
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
	Uin int
	UserName string
	NickName string
	HeadImgUrl string
	ContactFlag int
	MemberCount int
	MemberList []string
	RemarkName string
	HideInputBarFlag int
	Sex int
	Signature string
	VerifyFlag int
	OwnerUin int
	PYInitial string
	PYQuanPin string
	RemarkPYInitial string
	RemarkPYQuanPin string
	StarFriend int
	AppAccountFlag int
	Statues int
	AttrStatus int
	Province string
	City string
	Alias string
	SnsFlag int
	UniFriend int
	DisplayName string
	ChatRoomId int
	KeyWord string
	EncryChatRoomId string
	IsOwner int
}

type User struct {
	Uin int
	UserName string
	NickName string
	HeadImgUrl string
	RemarkName string
	PYInitial string
	PYQuanPin string
	RemarkPYInitial string
	RemarkPYQuanPin string
	HideInputBarFlag int
	StarFriend int
	Sex int
	Signature string
	AppAccountFlag int
	VerifyFlag int
	ContactFlag int
	WebWxPluginSwitch int
	HeadImgFlag int
	SnsFlag int
}



type MPSubscribeMsg struct {

}

type WebInit struct {
	BaseResponse BaseResponse
	Count int
	ContactList []Contact
	SyncKey SyncKey

	User User
	ChatSet string
	SKey string
	ClientVersion int
	SystemTime int
	GrayScale int
	InviteStartCount int
	MPSubscribeMsgCount int
	MPSubscribeMsgList []interface{}
	ClickReportInterval int
}

