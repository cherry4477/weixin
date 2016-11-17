package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"
)

var (
	token      string
	tokenMutex sync.Mutex
)

func gettoken() string {
	tokenMutex.Lock()
	defer tokenMutex.Unlock()
	return token

}

func updatatoken() {

	f := func() {
		v := url.Values{}
		v.Set("grant_type", "client_credential")
		v.Set("appid", "wxd653a9d6ef5659ab")
		v.Set("secret", "114967dd70de0c89469d94f3ef493d35")
		//url:=url.URL
		r, err := http.Get("https://api.weixin.qq.com/cgi-bin/token?" + v.Encode())
		if err != nil {
			return
		}
		if r != nil {
			defer r.Body.Close()
		}

		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return
		}
		//GetResponseData(r)
		//json.Unmarshal
		var params = struct {
			Access  string `json:"access_token"`
			Expires int64  `json:"expires_in"`
		}{}
		err = json.Unmarshal(data, &params)
		if err != nil {
			return
		}

		log.Println("params", params)
		tokenMutex.Lock()
		token = params.Access
		tokenMutex.Unlock()

	}
	f()
	for range time.Tick(time.Hour) {
		f()
	}
}

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()       //解析参数，默认是不会解析的
	fmt.Println(r.Form) //这些信息是输出到服务器端的打印信息
	// fmt.Println("path", r.URL.Path)
	// fmt.Println("scheme", r.URL.Scheme)
	// fmt.Println(r.Form["url_long"])
	// for k, v := range r.Form {
	// 	// fmt.Println("key:", k)
	// 	// fmt.Println("val:", strings.Join(v, ""))
	// }
	fmt.Fprintf(w, "Hello astaxie!") //这个写入到w的是输出到客户端的
}
func weixinin(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()       //解析参数，默认是不会解析的
	log.Println(r.Form) //这些信息是输出到服务器端的打印信息
	// log.Println("path", r.URL.Path)
	// log.Println("scheme", r.URL.Scheme)
	// log.Println(r.Form["url_long"])
	fmt.Println("----->here")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	fmt.Println("----->body:", string(body))
	for k, v := range r.Form {
		log.Println("key:", k)
		log.Println("val:", strings.Join(v, ""))
	}
	type one struct {
		Content string `json:"content"`
	}
	var obj = struct {
		Touser  string `json:"touser"`
		Msgtype string `json:"msgtype"`
		Text    one    `json:"text"`
	}{
		Touser:  r.FormValue("openid"),
		Msgtype: "text",
		Text: one{
			Content: "Hello World",
		},
	}

	data, err := json.Marshal(&obj)
	if err != nil {
		return
	}
	// h.Write([]byte(tmpStr))

	request, data, err := RemoteCallWithBody(
		"POST",
		"https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token="+token,
		"",
		"",
		data,
		"application/json; charset=utf-8",
	)

	if err != nil {
		return
	}
	log.Println("data", data)
	log.Println("request", request)
	// if checkSignature(r) {
	// 	fmt.Fprint(w, r.FormValue("echostr"))
	// } else {
	// 	fmt.Fprint(w, "hello wixin sb ") //这个写入到w的是输出到客户端的
	// }

}

//https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=ACCESS_TOKEN
func follow(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		return
	}
	data, err := GetRequestData(r)
	if err != nil {
		return
	}
	var params = struct {
		ToUserName   string `xml:"ToUserName"`
		FromUserName string `xml:"FromUserName"`
		CreateTime   int64  `xml:"CreateTime"`
		MsgType      string `xml:"MsgType"`
		Event        string `xml:"Event"`
		EventKey     string `xml:"EventKey"`
		Ticket       string `xml:"Ticket"`
	}{}

	err = xml.Unmarshal(data, &params)
	if err != nil {
		return
	}

	// log.Println(params)

	type one struct {
		Content string `json:"content"`
	}
	var obj = struct {
		Touser  string `json:"touser"`
		Msgtype string `json:"msgtype"`
		Text    one    `json:"text"`
	}{
		Touser:  params.FromUserName,
		Msgtype: "text",
		Text: one{
			Content: "Hello World",
		},
	}

	data, err = json.Marshal(&obj)
	if err != nil {
		return
	}
	// h.Write([]byte(tmpStr))

	request, data, err := RemoteCallWithBody(
		"POST",
		"https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token="+token,
		"",
		"",
		data,
		"application/json; charset=utf-8",
	)

	if err != nil {
		return
	}
	log.Println("data", data)
	log.Println("request", request)
	// {
	//     "touser":"OPENID",
	//     "msgtype":"text",
	//     "text":
	//     {
	//          "content":"Hello World"
	//     }
	// }

	// <xml>
	// <ToUserName><![CDATA[toUser]]></ToUserName>
	// <FromUserName><![CDATA[FromUser]]></FromUserName>
	// <CreateTime>123456789</CreateTime>
	// <MsgType><![CDATA[event]]></MsgType>
	// <Event><![CDATA[SCAN]]></Event>
	// <EventKey><![CDATA[SCENE_VALUE]]></EventKey>
	// <Ticket><![CDATA[TICKET]]></Ticket>
	// </xml>

}
func main() {
	go updatatoken()
	http.HandleFunc("/", sayhelloName) //设置访问的路由
	http.HandleFunc("/interface", weixinin)
	http.HandleFunc("/follow", follow)
	err := http.ListenAndServe(":9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// return RemoteCallWithBody(method, url, token, user, jsonBody, "application/json; charset=utf-8")
func RemoteCallWithBody(method, url string, token, user string, body []byte, contentType string) (*http.Response, []byte, error) {

	var request *http.Request
	var err error
	if len(body) == 0 {
		request, err = http.NewRequest(method, url, nil)
	} else {
		request, err = http.NewRequest(method, url, bytes.NewReader(body))
	}
	if err != nil {
		return nil, nil, err
	}
	if contentType != "" {
		request.Header.Set("Content-Type", contentType)
	}
	if token != "" {
		request.Header.Set("Authorization", token)
	}
	if user != "" {
		request.Header.Set("User", user)
	}
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	response, err := client.Do(request)
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return nil, nil, err
	}

	bytes, err := ioutil.ReadAll(response.Body)
	return response, bytes, err
}

func GetResponseData(r *http.Response) ([]byte, error) {
	if r != nil {
		defer r.Body.Close()
	}

	return ioutil.ReadAll(r.Body)

}
func GetRequestData(r *http.Request) ([]byte, error) {
	if r.Body == nil {
		return nil, nil
	}

	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func checkSignature(r *http.Request) bool {
	signature := r.FormValue("signature")
	timestamp := r.FormValue("timestamp")
	nonce := r.FormValue("nonce")
	token := "winxin"
	tmpArr := sort.StringSlice{token, timestamp, nonce}
	sort.Sort(tmpArr)
	tmpStr := strings.Join(tmpArr, "")
	//产生一个散列值得方式是 sha1.New()，sha1.Write(bytes)，然后 sha1.Sum([]byte{})。这里我们从一个新的散列开始。
	h := sha1.New()
	//写入要处理的字节。如果是一个字符串，需要使用[]byte(s) 来强制转换成字节数组。
	h.Write([]byte(tmpStr))
	//这个用来得到最终的散列值的字符切片。Sum 的参数可以用来都现有的字符切片追加额外的字节切片：一般不需要要。
	bs := h.Sum(nil)
	//SHA1 值经常以 16 进制输出，例如在 git commit 中。使用%x 来将散列结果格式化为 16 进制字符串。

	// fmt.Println(hex.EncodeToString(bs))
	// fmt.Printf("%s\n", signature)

	// fmt.Println(hex.EncodeToString(bs) == signature)

	if hex.EncodeToString(bs) == signature {
		return true
	} else {
		return false
	}
}
