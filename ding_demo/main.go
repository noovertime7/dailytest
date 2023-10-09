package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Message struct {
	Type     string   `json:"msgtype"`
	Markdown Markdown `json:"markdown"`
}

type Markdown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

var ACCESS_TOKEN = "ccfd8cf271684da5c50e7ff73c69b6ebb71d1cacccd4ba11e830296208e4b590"

func main() {
	msg := Message{
		Type:     "markdown",
		Markdown: Markdown{Title: "测试", Text: ">测试1\n\n>测试2\n\n>测试3\n\n[发送消息](dtmd://dingtalkclient/sendMessage?content=认领告警，告警ID[ccfd8cf271684da5c50e7ff73c69b6ebb71d1cacccd4ba11e830296208e4b590])"},
	}
	body, _ := json.Marshal(msg)
	client := &http.Client{}
	req, _ := http.NewRequest("POST", fmt.Sprintf("https://oapi.dingtalk.com/robot/send?access_token=%s", ACCESS_TOKEN), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	resp, _ := client.Do(req)
	if resp.StatusCode == 200 {
		fmt.Println("发送成功")
	}
}
