package main

import (
	"file_websocket_demo/runtime"
	"fmt"
	"time"
)

func main() {
	// 指定文件夹和文件名
	folderPath := "/Users/admin/code/kubemanage/log/workflow/7164494158366175232"
	fileName := "delayTime"

	builder := runtime.FileBuilder{}
	w := builder.FolderName(folderPath).FileName(fileName).Build()

	for {
		time.Sleep(1 * time.Second)
		fmt.Println(w.Text())
	}

}
