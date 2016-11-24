package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
)

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()                   //解析参数，默认是不会解析的
	fmt.Println(r.Form)             //这些信息是输出到服务器端的打印信息
	fmt.Fprintf(w, "Hello weixin!") //这个写入到w的是输出到客户端的
}

func follow(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, r.FormValue("echostr"))
	if r.Method != "POST" {
		return
	}
	r.ParseForm() //解析参数，默认是不会解析的
	//log.Println(r.Form) //这些信息是输出到服务器端的打印信息
	log.Println("From", r.RemoteAddr, r.Method, r.URL.RequestURI(), r.Proto)

	data, err := GetRequestData(r)
	if err != nil {
		return
	}
	//createtime
	var common = struct {
		FromUserName string `xml:"FromUserName"`
		MsgType      string `xml:"MsgType"`
		Event        string `xml:"Event"`
		CreateTime   int64  `xml:"CreateTime"`
	}{}

	err = xml.Unmarshal(data, &common)
	if err != nil {
		return
	}

	if common.MsgType == "event" {
		if common.Event == "subscribe" {

			var send = struct {
				OpenID      string `json:"openId"`
				ProvideTime int64  `json:"provideTime"`
			}{
				OpenID:      common.FromUserName,
				ProvideTime: common.CreateTime,
			}

			data, err = json.Marshal(&send)

			if err != nil {
				return
			}

			// log.Println("star", data)
			_, data, err := RemoteCallWithBody(
				"POST",
				"http://datafoundry.coupon.app.dataos.io/charge/v1/provide/coupons?number=1",
				"",
				"",
				data,
				"application/json; charset=utf-8",
			)

			if err != nil {
				// log.Println("err", err)
				return
			}

			// log.Println("end", data)

			log.Println("pass")
			type two struct {
				IsProvide bool   `json:"isProvide"`
				Code      string `json:"code"`
			}
			type Result struct {
				Code uint        `json:"code"`
				Msg  string      `json:"msg"`
				Data interface{} `json:"data,omitempty"`
			}
			var three = two{}
			var card = Result{

				Data: &three,
			}

			log.Println("data:", string(data))
			err = json.Unmarshal(data, &card)

			if err != nil {
				log.Println("err2:", err)
				return
			}

			log.Println("Code", three.Code)

			log.Println("IsProvide", three.IsProvide)

			if three.IsProvide {
				return
			}

			type one struct {
				Content string `json:"content"`
			}
			var str = three.Code[:4] + "-" + three.Code[4:8] + "-"
			str = str + three.Code[8:12] + "-" + three.Code[12:16]
			var obj = struct {
				Touser  string `json:"touser"`
				Msgtype string `json:"msgtype"`
				Text    one    `json:"text"`
			}{
				Touser:  common.FromUserName,
				Msgtype: "text",
				Text: one{
					Content: "您的充值卡号为" + str + "，有效期截止至2017年02月31日",
				},
			}

			data, err = json.Marshal(&obj)
			if err != nil {
				return
			}

			request, data, err := RemoteCallWithBody(
				"POST",
				"https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token="+gettoken(),
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
		}
	}

	// if checkSignature(r) {
	// 	fmt.Fprint(w, r.FormValue("echostr"))
	// } else {
	// 	fmt.Fprint(w, "hello wixin sb ") //这个写入到w的是输出到客户端的
	// }

}
