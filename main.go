package main

import (
	"log"
	"net/http"
)

func main() {
	go updatatoken()
	http.HandleFunc("/", sayhelloName) //设置访问的路由
	http.HandleFunc("/interface", follow)
	err := http.ListenAndServe(":9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
