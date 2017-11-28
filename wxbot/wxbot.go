package wxbot

import (
	"wxbot2/response"
	"net/http"
	"io/ioutil"
	"regexp"
	"time"
	"strconv"
	"fmt"
	"wxbot2/httpClient"
	"net/url"
	"encoding/xml"
	"encoding/json"
	"wxbot2/model"
)


type Data struct {
	Uuid string
	PreUri string
	XmlConfig *XmlConfig
	DeviceID string
	Cookies []*http.Cookie
	Plug 		string
}
type XmlConfig struct {
	XMLName     xml.Name `xml:"error"`
	Ret         int      `xml:"ret"`
	Message     string   `xml:"message"`
	Skey        string   `xml:"skey"`
	Wxsid       string   `xml:"wxsid"`
	Wxuin       string   `xml:"wxuin"`
	PassTicket  string   `xml:"pass_ticket"`
	IsGrayscale int      `xml:"isgrayscale"`

}

type Wxbot struct {
	Data
	InitData *response.WebInit
	UserList []response.User
	LastId int64
	User response.User
	SyncKey response.SyncKey
	Handler func(msg []byte)
	InitOver func()
	SpecialContact map[string]bool
	MsgHander *MsgHander
	Groups []response.Contact

}

func(bot *Wxbot) Run(cache bool){
	if !cache {
		bot.DeviceID = GetRandomStringFromNum(15)
		fmt.Println(bot.GetQrCode())
		for i := 0; i < 10; i++ {
			fmt.Println("belogin")
			status_code, redirect_url, err := bot.Login()
			if err != nil {
				panic( err)
			}
			if status_code == 200 {
				fmt.Println("WebNewLoginPage")
				bot.WebNewLoginPage(redirect_url)

				break
			} else if status_code == 201 {
				fmt.Println("Press login on your phone")
				time.Sleep(3 * time.Second)
			} else if status_code == 400 {
				fmt.Errorf("login timeout, this qr is no longer valid, restart again\n")
				return
			} else if status_code == 408 {
				time.Sleep(25 * time.Second)
			}
		}

	}

	bot.WebWxInit()

	bot.StatusNotify()

	bot.WebWxGetContact()
	//i := 0
	for{
		fmt.Println("SyncCheck")
		status,selecter,err := bot.SyncCheck()
		if err!= nil{
			//i++
			continue
		}
		if status != 0{
			fmt.Println("status :",status,bot.Uuid)
			return
		}
		if selecter == 2{
			fmt.Println("WebWxSync")
			bot.WebWxSync()
			//if i > 0{
			//	fmt.Println("StatusNotify")
			//	bot.StatusNotify()
			//	i = 0
			//
			//}
			continue
		}
		//i++
	}
}


func Init(data *Data) *Wxbot{
	bot := new(Wxbot)
	uuid,_:= bot.JsLogin()
	bot.Uuid = uuid
	if data != nil{
		bot.Data = *data
	}
	bot.MsgHander = new(MsgHander)
	bot.SpecialContact = map[string]bool{
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
	return bot
}



func (w *Wxbot) JsLogin() (string,error){
	w.LastId = getMilSecond()

	urr := "https://login.wx.qq.com/jslogin?appid=wx782c26e4c19acffb&redirect_uri=https%3A%2F%2Fwx.qq.com%2Fcgi-bin%2Fmmwebwx-bin%2Fwebwxnewloginpage&fun=new&lang=zh_CN&_=" +strconv.FormatInt(w.LastId,10)
	w.LastId = w.LastId +1
	resp,err := httpClient.Get(urr,nil,nil)
	defer resp.Body.Close()
	if err != nil{
		return "",err
	}
	body, err :=  ioutil.ReadAll(resp.Body)
	r := regexp.MustCompile(`window\.QRLogin\.code *= *(\d{3}) *; *window\.QRLogin\.uuid *= *"(\S+)"`)
	match := r.FindStringSubmatch(string(body))
	if len(match) == 0 {
		return "", fmt.Errorf("can not find qr_login_uuid from %s", string(body))
	} else if match[1] != "200" {
		return "", fmt.Errorf("qr_login_code is not 200, %s", string(match[1]))
	}
	w.Uuid = match[2]
	return match[2], nil
}




func (w *Wxbot) Login() (int64 ,string ,error){
	unix_time := get_unix_time(13)
	unix_time_int, _ := strconv.ParseInt(unix_time, 10, 64)
	r := ^unix_time_int & 0xFFFFFFFF
	uri := "https://login.wx.qq.com/cgi-bin/mmwebwx-bin/login?loginicon=true&uuid="+w.Uuid+"&tip=0&r="+string(r)+"&_="+strconv.FormatInt(w.LastId,10)
	w.LastId = w.LastId +1
	resp,err := httpClient.Get(uri,nil,nil)
	defer resp.Body.Close()
	if err != nil {
		return 0, "", err
	}
	body,_ := ioutil.ReadAll(resp.Body)
	status_regexp := regexp.MustCompile(`window\.code *= *(\d+)`)
	var status_code int64 = 500
	var redirect_url string = ""
	match := status_regexp.FindStringSubmatch(string(body))
	if len(match) == 0 {
		return status_code, redirect_url, fmt.Errorf("no window.code in %s\n", body)
	}

	status_code, err = strconv.ParseInt(match[1], 10, 64)
	if err != nil {
		return status_code, redirect_url, err
	}

	if status_code == 200 {
		redirect_regexp := regexp.MustCompile(`window.redirect_uri *= *"(\S+)"`)
		match = redirect_regexp.FindStringSubmatch(string(body))
		if len(match) == 0 {
			return status_code, redirect_url, fmt.Errorf("no window.redirect_uri in %s\n", body)
		}
		redirect_url = match[1]
		return status_code, redirect_url, nil
	}
	return status_code, redirect_url, nil
}

func (w *Wxbot) GetQrCode() string{
	if w.Uuid == ""{
		w.JsLogin()
	}
	url := "https://login.weixin.qq.com/qrcode/" + w.Uuid
	//resp ,_ := http.Get(url)
	return url

}

func (w *Wxbot)WebNewLoginPage(uri string)( []*http.Cookie,error){
	u, _ := url.Parse(uri)
	km := u.Query()
	w.PreUri = u.Scheme +"://"+u.Host
	km.Add("fun", "new")
	//&fun=new&version=v2&lang=zh_CN
	//uri = w.PreUri + "/cgi-bin/mmwebwx-bin/webwxnewloginpag?" + km.Encode()
	uri = uri + "&fun=new&version=v2&lang=zh_CN"
	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	xc := new (XmlConfig)
	if err := xml.Unmarshal(body, xc); err != nil {
		return nil, err
	}

	if xc.Ret != 0 {
		return nil, fmt.Errorf("xc.Ret != 0, %s", string(body))
	}
	w.XmlConfig = xc
	w.Cookies = resp.Cookies()
	return w.Cookies, nil
}

func (w *Wxbot)WebWxInit() (*response.WebInit, error) {
	unix_time := get_unix_time(13)
	unix_time_int, _ := strconv.ParseInt(unix_time, 10, 64)
	r := ^unix_time_int & 0xFFFFFFFF
	km := url.Values{}
	km.Add("pass_ticket", w.XmlConfig.PassTicket)
	km.Add("skey", w.XmlConfig.Skey)
	km.Add("r", string(r))

	uri := w.PreUri+"/cgi-bin/mmwebwx-bin/webwxinit?" + km.Encode()

	req := response.BaseRequest{
		w.XmlConfig.Wxuin,
		w.XmlConfig.Wxsid,
		w.XmlConfig.Skey,
		w.DeviceID,
	}
	q := make(map[string]response.BaseRequest)
	q["BaseRequest"] = req

	b, _ := json.Marshal(q)
	resp, err := httpClient.Post(uri,b,w.Cookies,nil)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	initData := new(response.WebInit)
	jerr := json.Unmarshal(body,initData)
	if err != nil{
		return nil,jerr
	}
	w.InitData = initData
	w.User = initData.User
	w.SyncKey = initData.SyncKey
	if w.InitOver != nil{
		go w.InitOver()
	}
	return initData, nil
}

//https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxstatusnotify?lang=zh_CN&pass_ticket=bYlThDVpZxFRt4Xt1PVAOb8SSp9kP7o9P5LQ1vDvM9XMAiTX%252FoUi285aj2NNpRhB
func (w *Wxbot)StatusNotify() string{
	uri := w.PreUri+"/cgi-bin/mmwebwx-bin/webwxstatusnotify?lang=zh_CN&pass_ticket="+w.XmlConfig.PassTicket
	params := make(map[string]interface{})
	BaseRequest := response.BaseRequest{
		w.XmlConfig.Wxuin,
		w.XmlConfig.Wxsid,
		w.XmlConfig.Skey,
		w.DeviceID,
	}
	params["BaseRequest"] = BaseRequest
	params["Code"] = 3
	params["ClientMsgId"] = time.Now().UnixNano()/1000000
	params["FromUserName"] = w.User.UserName
	params["ToUserName"] = w.User.UserName
	b,_ := json.Marshal(params)

	resp ,_ := httpClient.Post(uri,b,w.Cookies,nil)
	defer resp.Body.Close()
	body,_ := ioutil.ReadAll(resp.Body)
	v := make(map[string]string)
	json.Unmarshal(body,&v)
	return v["MsgID"]
}

func (w *Wxbot)WebWxGetContact() []response.User {
	///cgi-bin/mmwebwx-bin/webwxgetcontact?lang=zh_CN&pass_ticket=bYlThDVpZxFRt4Xt1PVAOb8SSp9kP7o9P5LQ1vDvM9XMAiTX%252FoUi285aj2NNpRhB&r=1510663154398&seq=0&skey=@crypt_c8432a31_bd49746d86b43e194795d2927dc77fb9
	t := getMilSecond()
	uri := w.PreUri + "/cgi-bin/mmwebwx-bin/webwxgetcontact?lang=zh_CN&r=" +strconv.FormatInt(t,10) +"&seq=0&skey="+w.XmlConfig.Skey
	resp,_ := httpClient.Get(uri,w.Cookies,nil)
	body,_ := ioutil.ReadAll(resp.Body)
	respData := new(response.WebWxGetContact)
	json.Unmarshal(body,respData)
	w.UserList = respData.MemberList
	return respData.MemberList
}

//https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxbatchgetcontact?type=ex&r=1511248546451
/*
func (w  *Wxbot)WebWxbatchGetcontact(){
	uri := w.PreUri + "/cgi-bin/mmwebwx-bin/webwxbatchgetcontact?type=ex&r="+strconv.FormatInt(time.Now().Unix()*1000, 10)
	resp,err := httpClient.Get(uri,w.Cookies,nil)
	if err != nil{
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body,_ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body,)
}
*/
func (w *Wxbot)SyncCheck() (int,int ,error) {
//	https://webpush.wx2.qq.com/cgi-bin/mmwebwx-bin/synccheck?r=1510709778549&skey=%40crypt_c07a7e97_4f5d7621518c700184bcf4600ceba154&sid=VmxB3fJdLDDC%2BV0w&uin=515008662&deviceid=e597615953927488&synckey=1_648370450%7C2_648370575%7C3_648370572%7C1000_1510707242&_=1510709599793
 	km := url.Values{}
	km.Add("r", strconv.FormatInt(time.Now().Unix()*1000, 10))
	km.Add("sid", w.XmlConfig.Wxsid)
	km.Add("uin", w.XmlConfig.Wxuin)
	km.Add("skey", w.XmlConfig.Skey)
	km.Add("deviceid", w.DeviceID)
	km.Add("synckey", w.SyncKey.Encode())
	km.Add("_", strconv.FormatInt(time.Now().Unix()*1000, 10))
	uri :=  "https://webpush.wx2.qq.com/cgi-bin/mmwebwx-bin/synccheck?" + km.Encode()
	resp, err := httpClient.Get(uri,w.Cookies,nil)
	//fmt.Println(resp.Cookies())
	if err != nil {
		return 0, 0, err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	strb := string(body)
	fmt.Println(strb)
	reg := regexp.MustCompile("window.synccheck={retcode:\"(\\d+)\",selector:\"(\\d+)\"}")
	sub := reg.FindStringSubmatch(strb)
	retcode, _ := strconv.Atoi(sub[1])
	selector, _ := strconv.Atoi(sub[2])
	return retcode, selector, nil
}

//https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxsync?sid=VmxB3fJdLDDC+V0w&skey=@crypt_c07a7e97_4f5d7621518c700184bcf4600ceba154&lang=zh_CN
func (w *Wxbot)WebWxSync() error {
	km := url.Values{}
	km.Add("skey", w.XmlConfig.Skey)
	km.Add("sid", w.XmlConfig.Wxsid)
	km.Add("lang", "zh_CN")
	//km.Add("pass_ticket", w.XmlConfig.PassTicket)
	uri := w.PreUri + "/cgi-bin/mmwebwx-bin/webwxsync?" + km.Encode()

	params := make(map[string]interface{})
	BaseRequest := response.BaseRequest{
		w.XmlConfig.Wxuin,
		w.XmlConfig.Wxsid,
		w.XmlConfig.Skey,
		w.DeviceID,
	}
	params["BaseRequest"] = BaseRequest
	params["SyncKey"] = w.SyncKey
	params["rr"] = ^int(time.Now().Unix()) + 1

	b, _ := json.Marshal(params)

	resp, err := httpClient.Post(uri, b, w.Cookies, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	dta := new(response.Webwxsync)
	json.Unmarshal(body, dta)
	tmp := make(map[string]*http.Cookie)
	for _, v := range resp.Cookies() {
		tmp[v.Name] = v
	}
	for _ ,v := range w.Cookies{
		va,ok := tmp[v.Name]
		if ok{
			v = va
		}
	}
	w.SyncKey = dta.SyncCheckKey
	if w.MsgHander != nil  {
		w.MsgHander.MsgRec = body
		w.MsgHander.Wxbot = w
		go w.MsgHander.Init()

	}
	//fmt.Println(string(body))
	return nil
}

func (w *Wxbot)SendText(to string,content string) ([]byte,error){
	km := url.Values{}
	km.Add("pass_ticket", w.XmlConfig.PassTicket)
	km.Add("lang", "zh_CN")

	uri := w.PreUri + "/cgi-bin/mmwebwx-bin/webwxsendmsg?" + km.Encode()

	params := make(map[string]interface{})
	BaseRequest := response.BaseRequest{
		w.XmlConfig.Wxuin,
		w.XmlConfig.Wxsid,
		w.XmlConfig.Skey,
		w.DeviceID,
	}

	msg := model.TextMessage{
			Type:         1,
			Content:      content,
			FromUserName: w.User.UserName,
			ToUserName:   to,
			LocalID:      int(time.Now().Unix() * 1e4),
			ClientMsgId:  int(time.Now().Unix() * 1e4),
		}

	params["BaseRequest"] = BaseRequest
	params["Msg"] = msg
	params["Scene"] = 0
	b, _ := json.Marshal(params)

	resp, err := httpClient.Post(uri,b,w.Cookies,nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))
	return body, nil
}
