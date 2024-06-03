package main

import (
	"fmt"

	"github.com/goburrow/modbus"
)

func main() {
	handler := modbus.NewTCPClientHandler("127.0.0.1:65001")
	// Connect manually so that multiple requests are handled in one session
	err := handler.Connect()
	defer handler.Close()
	client := modbus.NewClient(handler)
	commond := []byte{0x01, 0x03, 0x01, 0xF4, 0x00, 0x06, 0x85, 0xC6}
	_, err = client.WriteMultipleRegisters(0, 3, commond)
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	results, err := client.ReadHoldingRegisters(0, 3)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	fmt.Printf("results %v\n", results)
}
