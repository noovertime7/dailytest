package main

import (
	"bufio"
	"fmt"
	"github.com/grafana/tail"
	"os"
	"time"
)

func main() {
	go write()
	t, err := tail.TailFile("./test.txt", tail.Config{Follow: true, MustExist: true, Poll: true})
	for line := range t.Lines {
		fmt.Println(line.Text)
	}
	if err != nil {
		panic(err)
	}
}

func write() {
	ff, _ := os.OpenFile("./test.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	writer := bufio.NewWriter(ff)

	for i := 0; i <= 100; i++ {
		time.Sleep(1 * time.Second)
		_, err := writer.WriteString(fmt.Sprintf("写入数据:%d\n", i))
		if err != nil {
			panic(err)
		}
		err = writer.Flush()
		if err != nil {
			panic(err)
		}
	}
}
