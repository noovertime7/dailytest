package main

import (
	"encoding/json"
	"fmt"
	"github.com/noovertime7/dailytest/jiqima_demo/ciper"
	"time"
)

type License struct {
	Name         string    `json:"name"`
	Capacity     int64     `json:"capacity"`
	AgentNum     int64     `json:"agent_num"`
	IssuanceTime time.Time `json:"issuance_time"`
	ExpireTime   time.Time `json:"expire_time"`
	Code         string    `json:"code"`
}

func main() {
	license := License{
		Name:         "lincess",
		Capacity:     100,
		AgentNum:     100,
		IssuanceTime: time.Now(),
		ExpireTime:   time.Now().Add(time.Hour * 24 * 365),
		Code:         "VPbgqKUsykSIfuv8O2Y5PiKSuEm4ihXN2Mo52fAdRqjkfvbkZq7tp4s7ayR75+tP",
	}

	l, _ := json.Marshal(license)

	data, _ := ciper.Encrypt(l)

	fmt.Println(data)

	decodeData, _ := ciper.Decrypt(data)

	fmt.Println(string(decodeData))
}
