package main

import (
	"fmt"
	"os/exec"
)

func main() {
	// 定义要执行的 Shell 脚本
	script := `#!/bin/bash 
ls

`

	// 使用 sh -x 命令执行脚本
	cmd := exec.Command("sh", "-x", "-c", script)

	// 执行命令并等待执行结果
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("执行脚本时出错:%v,[%v]", err, string(output))
		return
	}

	fmt.Println("执行成功:\n", string(output))

}
