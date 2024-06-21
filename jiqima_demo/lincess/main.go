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
		BackupCapacity: 200000000000,
		AgentNum:       30,
		IssuanceDate:   time.Now(),
		Code:           "8db63be3c8b5edb937045e6eebb7a011",
	}

	l, _ := json.Marshal(license)

	data, _ := ciper.Encrypt(l)

	fmt.Println(data)

	decodeData, _ := ciper.Decrypt(data)

	fmt.Println(string(decodeData))
}
