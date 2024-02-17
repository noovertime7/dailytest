package runtime

import (
	"bufio"
	"fmt"
	"os"
)

type FileBuilder struct {
	folderName string
	fileName   string
	fw         *file
}

func (f *FileBuilder) FileName(name string) *FileBuilder {
	f.fileName = name
	return f
}

func (f *FileBuilder) FolderName(name string) *FileBuilder {
	f.folderName = name
	return f
}

func (f *FileBuilder) Build() File {
	filePath := fmt.Sprintf("%s/%s", f.folderName, f.fileName)
	// 打开文件，如果文件不存在则创建
	ff, _ := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	writer := bufio.NewWriter(ff)
	scanner := bufio.NewScanner(ff)
	return &file{w: writer, file: ff, r: scanner}
}
