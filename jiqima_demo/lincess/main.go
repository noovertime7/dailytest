package main

import (
	"encoding/json"
	"fmt"
	"github.com/noovertime7/dailytest/jiqima_demo/ciper"
	"time"
)

type License struct {
	Name           string
	IsAuth         bool      `json:"isAuth"`
	AgentNum       int64     `json:"agentNum"`
	BackupCapacity int64     `json:"backupCapacity"`
	IssuanceDate   time.Time `json:"issuanceDate"`
	ExpireDate     time.Time `json:"expireDate"`
	Code           string    `json:"code"`
	UsedAgentNum   int64     `json:"usedAgentNum"`
	UsedCapacity   int64     `json:"usedCapacity"`
}

func main() {
	license := License{
		Name:           "license_test2",
		BackupCapacity: 3000000000000,
		AgentNum:       30,
		IssuanceDate:   time.Now(),
		Code:           "bc9f93a4ffbd056bef99ba811b2c0d5a",
	}

	l, _ := json.Marshal(license)

	data, _ := ciper.Encrypt(l)

	fmt.Println(data)

	decodeData, _ := ciper.Decrypt(data)

	fmt.Println(string(decodeData))
}
