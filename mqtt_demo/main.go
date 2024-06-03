package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"time"
)

const (
	topic = "/dyyq/65001"
	//topic = "/v1/opc/dyyq/lqsw"
)

func Mqtt() (mqtt.Client, error) {
	opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://dyyqmqtt.15net.top:60010"))
	opts.SetClientID("123")
	opts.SetUsername("admin")
	opts.SetPassword("Dyyq@2024")

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("连接mqtt失败:%v", token.Error())
	}
	return client, nil
}
func main() {
	client, err := Mqtt()
	if err != nil {
		log.Fatal(err)
	}
	// 订阅主题
	if token := client.Subscribe(topic, 0, onMessageReceived); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	log.Println("start")
	time.Sleep(10 * time.Hour)
}

func onMessageReceived(client mqtt.Client, message mqtt.Message) {
	fmt.Println("message = ", string(message.Payload()))
}
