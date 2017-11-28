package httpClient

import (
	"net/http"
	"bytes"
	"net/http/cookiejar"

)



func Get(u string, cookies []*http.Cookie,headers map[string]string) (*http.Response,error) {
	//bd := bytes.NewReader(body)
	req,_ := http.NewRequest("GET",u,nil)
	jar,_ := cookiejar.New(nil)
	jar.SetCookies(req.URL,cookies)
	hc := http.Client{}
	hc.Jar = jar
	default_headers := map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36",
		"Accept": "application/json, text/plain, */*",
		"Accept-Language": "zh-CN,zh;q=0.8",
		"Connection": "keep-alive",
		"Host": req.URL.Host,
	}

	for k, v := range default_headers {
		req.Header.Add(k, v)
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	return hc.Do(req)
}

func Post(u string ,body []byte, cookies []*http.Cookie,headers map[string]string)(*http.Response,error){
	bd := bytes.NewReader(body)
	req,_ := http.NewRequest("POST",u,bd)
	jar,_ := cookiejar.New(nil)
	jar.SetCookies(req.URL,cookies)
	h :=http.Client{}
	h.Jar = jar
	default_headers := map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36",
		"Accept": "application/json, text/plain, */*",
		"Accept-Language": "zh-CN,zh;q=0.8",
		"Connection": "keep-alive",
		"Host": req.URL.Host,
	}

	for k, v := range default_headers {
		req.Header.Add(k, v)
	}

	return h.Do(req)
}