package runtime

import (
	"bufio"
	"os"
)

type File interface {
	WriteString(content string) (nn int, err error)
	Flush() error
	Close() error
	Text() (string, error)
}

type file struct {
	file *os.File
	w    *bufio.Writer
	r    *bufio.Scanner
}

func (f *file) WriteString(content string) (nn int, err error) {
	return f.w.WriteString(content)
}

func (f *file) Flush() error {
	return f.w.Flush()
}

func (f *file) Close() error {
	return f.file.Close()
}

func (f *file) Text() (string, error) {
	if f.r.Scan() {
		return f.r.Text(), nil
	}

	return f.r.Text(), f.r.Err()
}
