package wx

import (
    "fmt"
    "log"
    "time"
    "regexp"
    "bytes"
    "strconv"
    "strings"
    "math/rand"
    "encoding/xml"
    "encoding/json"
    "io/ioutil"
    "net/url"
    "net/http"
    "net/http/cookiejar"
    "oa.com/melow_dog/qrterm"
    //"github.com/tonnerre/golang-pretty"
)

type BaseResponse struct {
    Ret int
    ErrMsg string
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

type Wxbot struct {
    name string
    session *http.Client

    wx_host string
    lang string
    login_prefix string
    file_prefix string
    webpush_prefix string


    api map[string]string
    login_info map[string]string

    special_user []string
    friend map[string]string
    group map[string]string
    mp map[string]string
    offical_user []string
}


func NewWxbot(name string, timeout int) *Wxbot {
    var bot Wxbot
    bot.name = name
    wx_jar, err := cookiejar.New(nil)
    if err != nil {
        log.Fatal(err)
    }

    bot.session = &http.Client{Jar: wx_jar, Timeout: time.Duration(timeout) *time.Second} 
    bot.lang = "zh_cn"
    bot.wx_host = "wx2.qq.com"
    bot.login_prefix = "login." + bot.wx_host
    bot.file_prefix = "file." + bot.wx_host

    bot.offical_user = []string{}
    bot.friend = make(map[string]string)
    bot.group = make(map[string]string)
    bot.mp = make(map[string]string)
    bot.login_info = make(map[string]string)
    bot.special_user = []string{
        "newsapp", "fmessage", "filehelper", "weibo", "qqmail",
        "fmessage", "tmessage", "qmessage", "qqsync", "floatbottle",
        "lbsapp", "shakeapp", "medianote", "qqfriend", "readerapp",
        "blogapp", "facebookapp", "masssendapp", "meishiapp",
        "feedsapp", "voip", "blogappweixin", "weixin", "brandsessionholder",
        "weixinreminder", "wxid_novlwrv3lqwv11", "gh_22b87fa7cb3c",
        "officialaccounts", "notification_messages", "wxid_novlwrv3lqwv11",
        "gh_22b87fa7cb3c", "wxitil", "userexperience_alarm", "notification_messages",
    }

    bot.api = make(map[string]string)
    bot.api["js_login"]  =  "https://" + bot.login_prefix + "/jslogin?appid=wx782c26e4c19acffb&fun=new&lang=" + bot.lang
    bot.api["check_login"]  =  "https://" + bot.login_prefix + "/cgi-bin/mmwebwx-bin/login"
    bot.api["web_init"]  =  "https://" + bot.wx_host + "/cgi-bin/mmwebwx-bin/webwxinit"
    bot.api["status_notify"]  =  "https://" + bot.wx_host + "/cgi-bin/mmwebwx-bin/webwxstatusnotify"
    bot.api["send_msg"]  =  "https://" + bot.wx_host + "/cgi-bin/mmwebwx-bin/webwxsendmsg"
    bot.api["get_contact"]  =  "https://" + bot.wx_host + "/cgi-bin/mmwebwx-bin/webwxgetcontact"
    return &bot
}

func (bot *Wxbot) CheckBaseResponse(resp_body string) error {
    var objmap map[string]*json.RawMessage
    err := json.Unmarshal([]byte(resp_body), &objmap)
    if err != nil {
        return err
    }

    var b BaseResponse
    err = json.Unmarshal(*objmap["BaseResponse"], &b)
    if err != nil {
        return err
    }
    if b.Ret != 0 {
        return fmt.Errorf("error code is %d\n", b.Ret) 
    }
    return nil
}

func (bot *Wxbot) Login() error {
    qr_login_uuid, err := bot.GetQrLoginUuid()
    if err != nil {
        return err
    }
    bot.login_info["qr_login_uuid"] = qr_login_uuid
    bot.DrawQrOnTty()

    var status_code int64
    var redirect_url string
    for i:=0; i<300; i++ {
        status_code, redirect_url, err = bot.CheckLogin()
        if err != nil {
            return err
        }
        if status_code == 200 {
            break
        } else if status_code  == 201 {
            fmt.Println("Press login on your phone")
        } else if status_code == 408 {
            return fmt.Errorf("login timeout, this qr is no longer valid, restart again\n")
        }
        time.Sleep(time.Second)
    }

    err = bot.GetLoginInfo(redirect_url)
    if err != nil {
        return err
    }
    err = bot.WebInit()
    if err != nil {
        return err
    }
    err = bot.StatusNotify()
    if err != nil {
        return err
    }
    err = bot.GetContact()
    if err != nil {
        return err
    }

    for i:=0; i<20; i ++ {
        err = bot.SendMsg("麦春明", fmt.Sprintf("麦仑%d号小管家叫你下班了...", i))
        if err != nil {
            return err
        }
        time.Sleep(time.Duration(100) * time.Microsecond)
    }

    return nil
}

func (bot *Wxbot) NewRequest(method string, url_str string, headers map[string]string, body string) (*http.Request, error) {
    url, err := url.Parse(url_str)
    host := url.Host
    default_headers := map[string]string{
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36",
        "Accept": "application/json, text/plain, */*",
        "Accept-Language": "zh-CN,zh;q=0.8",
        "Connection": "keep-alive",
        "Host": host,
    }

    b := bytes.NewBufferString(body)
    req, err := http.NewRequest(method, url_str, b)
    if err != nil {
        return req, err
    }

    for k, v := range default_headers {
        req.Header.Add(k, v)
    }

    for k,v := range headers {
        req.Header.Add(k, v)
    }

    return req, nil
}

func (bot *Wxbot) AddParams(req *http.Request, params map[string]string) {
    q := req.URL.Query() 
    for k, v := range params {
        q.Add(k, v)
    }
    req.URL.RawQuery = q.Encode()
}

func (bot *Wxbot) Do(req *http.Request) (string, error) {

    resp, err := bot.session.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return string(body), err
    }

    return string(body), err
}

func (bot *Wxbot) GetQrLoginUuid() (string, error) {
    req, err := bot.NewRequest("GET", bot.api["js_login"], nil, "")
    if err != nil {
        return "", err
    }
    body, err := bot.Do(req)
    r := regexp.MustCompile(`window\.QRLogin\.code *= *(\d{3}) *; *window\.QRLogin\.uuid *= *"(\S+)"`)
    match := r.FindStringSubmatch(string(body))
    if len(match) == 0 { 
         return "", fmt.Errorf("can not find qr_login_uuid from %s", string(body))
    } else if match[1] != "200" {
        return "", fmt.Errorf("qr_login_code is not 200, %s", string(match[1]))
    }   
    return match[2], nil 
}

func (bot *Wxbot) DrawQrOnTty() {
    scan_url := "https://login.weixin.qq.com/l/" + bot.login_info["qr_login_uuid"]
    qrterm.Draw(scan_url)
}

func get_unix_time(n uint8) string {
    unix_time := time.Now().UnixNano()
    return strconv.Itoa(int(unix_time))[:n]
}


func (bot *Wxbot) CheckLogin() (int64, string, error) {
    req, err := bot.NewRequest("GET", bot.api["check_login"], nil, "")
    unix_time := get_unix_time(13)
    unix_time_int, _ := strconv.ParseInt(unix_time, 10, 64)
    r := ^unix_time_int & 0xFFFFFFFF

    params := map[string]string{
        "loginicon": "false",
        "uuid": bot.login_info["qr_login_uuid"],
        "tip": "1",
        "r": string(r),
        "_": unix_time,
    }
    bot.AddParams(req, params)
    body, err := bot.Do(req) 
    if err != nil {
        return 0, "", err
    }

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

func (bot *Wxbot) GetLoginInfo(redirect_url string) error {
    req, err := bot.NewRequest("GET", redirect_url, nil, "")
    if err != nil {
        return err
    }
    params := map[string]string{
        "version": "v2",
        "fun": "new",
    }  

    bot.AddParams(req, params)
    //pretty.Println(req)
    body, err := bot.Do(req)
    if err != nil {
        return err
    }


    type msg struct {
        Error xml.Name `xml:"error"`
        Ret string `xml:"ret"`
        Message string `xml:"message"`
        Skey string `xml:"skey"`
        Wxsid string `xml:"wxsid"`
        Wxuin string `xml:"wxuin"`
        Pass_ticket string `xml:"pass_ticket"`
        Isgrayscale string `xml:"isgrayscale"`
    }
    v := msg{}
    err = xml.Unmarshal([]byte(body), &v)
    if err != nil {
        return err
    }

    status_code, err := strconv.ParseInt(v.Ret, 10, 64)
    if err != nil {
        return err
    }
    if status_code != 0 {
        return fmt.Errorf("fail to get logininfo from %s\n", redirect_url)
    }

    bot.login_info["skey"] = v.Skey
    bot.login_info["wx_sid"] = v.Wxsid
    bot.login_info["wx_uin"] = v.Wxuin
    bot.login_info["pass_ticket"] = v.Pass_ticket
    return nil
}

func randSeq(n int) string {
     var chars = []rune("0123456789")
     
     b := make([]rune, n)
     for i := range b {
         b[i] = chars[rand.Intn(len(chars))]
     }
     return string(b)
}

func (bot *Wxbot) getBaseRequest() string {
    return fmt.Sprintf(`"BaseRequest":{"Uin":"%s","Sid":"%s","Skey":"%s","DeviceID":"%s"}`, 
        bot.login_info["wx_uin"],
        bot.login_info["wx_sid"],
        bot.login_info["skey"],
        bot.login_info["device_id"],
    )
}

func (bot *Wxbot) WebInit() error {
    header := map[string]string{"Content-Type":"application/json;charset=UTF-8"} 
    bot.login_info["device_id"] = "e" + randSeq(15)
    data := fmt.Sprintf("{%s}", bot.getBaseRequest())
    params := map[string]string{
        "lang": bot.lang,
        "pass_ticket": bot.login_info["pass_ticket"],
    }
    req, err := bot.NewRequest("POST", bot.api["web_init"], header, data)
    if err != nil {
        return err
    }
    bot.AddParams(req, params)

    body, err := bot.Do(req)
    if err != nil {
        return err
    }

    type user struct {
        Uin int
        UserName string
        ickName string
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
 
    type sync_key struct {
        Count int
        List [] struct { Key int
                        Val int}
    }

    type huge_resp struct {
        BaseResponse  struct { Ret int
                               ErrMsg string 
                             }
        Count int
        ContactList []map[string]interface{}
        SyncKey sync_key
        User user
        ChatSet string
        Skey string
        ClientVersion int
        SystemTime int
        GrayScale int
        InviteStartCount int
        MPSubscribeMsgCount int
        MPSubscribeMsgList []map[string]interface{}
        ClickReportInterval int
    }
    var s huge_resp
    json.Unmarshal([]byte(body), &s)
    if s.BaseResponse.Ret != 0 {
        return fmt.Errorf("ret value not zero from %s", bot.api["web_init"])
    }

    synckey := ""
    for i:=0; i<s.SyncKey.Count; i++ {
        k := s.SyncKey.List[i].Key
        v := s.SyncKey.List[i].Val
        if i == 3 {
            synckey += strconv.Itoa(k) + "_" + strconv.Itoa(v) 
        } else {
            synckey += strconv.Itoa(k) + "_" + strconv.Itoa(v) + "|"
        }
    }
    bot.login_info["sync_key"] = synckey
    bot.login_info["user_name"] = s.User.UserName
    return nil
}

func (bot *Wxbot) StatusNotify() error {
    unix_time := get_unix_time(13)
    header := map[string]string{"Content-Type":"application/json;charset=UTF-8"} 
    data := fmt.Sprintf(`{%s,"Code":3,"FromUserName":"%s","ToUserName":"%s", "ClientMsgId":%s}`, 
            bot.getBaseRequest(), bot.login_info["user_name"], bot.login_info["user_name"], unix_time)

    fmt.Println(data)
    params := map[string]string{
        "lang": bot.lang,
        "pass_ticket": bot.login_info["pass_ticket"],
    }

    req, err := bot.NewRequest("POST", bot.api["status_notify"], header, data)
    if err != nil {
        return err
    }
    bot.AddParams(req, params)
    //pretty.Println(req)
    body, err := bot.Do(req)
    if err != nil {
        return err
    }

    err = bot.CheckBaseResponse(body)
    if err != nil {
        return err
    }
    return nil
}

func (bot *Wxbot) SendMsg(to_user string, msg string) error {
    unix_time := get_unix_time(17)
    header := map[string]string{"Content-Type":"application/json;charset=UTF-8"} 
    data := fmt.Sprintf(
        `{%s,"Msg":{"Type":1,"Content":"%s","FromUserName":"%s","ToUserName":"%s","LocalID":"%s","ClientMsgId":"%s"},"Scene":0}`, 
          bot.getBaseRequest(), msg, bot.login_info["user_name"], bot.friend[to_user], unix_time, unix_time)
    params := map[string]string{
        "lang": bot.lang,
        "pass_ticket": bot.login_info["pass_ticket"],
    }

    req, err := bot.NewRequest("POST", bot.api["send_msg"], header, data)
    if err != nil {
        return err
    }
    bot.AddParams(req, params)
    body, err := bot.Do(req)
    if err != nil {
        return err
    }

    err = bot.CheckBaseResponse(body)
    if err != nil {
        return err
    }
    
    return nil
}

func (bot *Wxbot) GetContact() error {
    unix_time := get_unix_time(13)

    params := map[string]string{
        "lang": bot.lang,
        "pass_ticket": bot.login_info["pass_ticket"],
        "seq": "0",
        "skey": bot.login_info["skey"],
        "r": unix_time,
    }

    req, err := bot.NewRequest("POST", bot.api["get_contact"], nil, "")
    if err != nil {
        return err
    }
    bot.AddParams(req, params)
    body, err := bot.Do(req)
    if err != nil {
        return err
    }

    err = bot.CheckBaseResponse(body)
    if err != nil {
        return err
    }

    type contact_resp struct {
        BP BaseResponse
        MemberCount int
        MemberList []Contact
        Seq int
    }

    var cr contact_resp
    json.Unmarshal([]byte(body), &cr)

    for _, c := range cr.MemberList {
        if c.VerifyFlag & 8 != 0 {
           bot.mp[c.NickName] = c.UserName 
        } else if stringInSlice(c.UserName, bot.special_user) {
            _ = append(bot.offical_user, c.NickName)
        } else if strings.HasPrefix(c.UserName, "@@") {
           bot.group[c.NickName] = c.UserName
        } else if c.UserName == bot.login_info["user_name"] {
            //do nothing
        } else {
            bot.friend[c.NickName] = c.UserName
        }
    }

    return nil
}

func stringInSlice(s string, list []string) bool {
    for _, b := range list {
        if b == s {
            return true
        }
    }
    return false
}

