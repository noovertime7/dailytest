package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/noovertime7/dailytest/jiqima_demo/ciper"
	"net"
	"strings"
)

// getMACAddresses retrieves the MAC addresses of the machine
func getMACAddresses() ([]string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var macs []string
	for _, interf := range interfaces {
		if interf.Flags&net.FlagUp != 0 && interf.Flags&net.FlagLoopback == 0 {
			addr := interf.HardwareAddr.String()
			if addr != "" {
				macs = append(macs, addr)
			}
		}
	}
	return macs, nil
}

// generateMachineID combines MAC addresses and CPU serial to generate a machine ID
func generateMachineID() (string, error) {
	macAddresses, err := getMACAddresses()
	if err != nil {
		return "", err
	}

	hardwareInfo := strings.Join(append(macAddresses, ""), "")
	hash := md5.Sum([]byte(hardwareInfo))
	return hex.EncodeToString(hash[:]), nil
}

func main() {
	machineID, err := generateMachineID()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(ciper.Encrypt([]byte(machineID)))
}
