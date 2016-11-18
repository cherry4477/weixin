package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
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

func main() {
	go updatatoken()
	http.HandleFunc("/", sayhelloName) //设置访问的路由
	http.HandleFunc("/interface", follow)
	err := http.ListenAndServe(":9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// RemoteCallWithBody send http
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

// GetResponseData parse Response
func GetResponseData(r *http.Response) ([]byte, error) {
	if r != nil {
		defer r.Body.Close()
	}

	return ioutil.ReadAll(r.Body)

}

// GetRequestData parse Request
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

	if hex.EncodeToString(bs) == signature {
		return true
	}
	return false

}
