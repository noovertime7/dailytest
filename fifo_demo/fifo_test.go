package fifo_demo

import (
	"fmt"
	"testing"
	"time"
)

type testFifoObject struct {
	name string
	val  interface{}
}

func testFifoObjectKeyFunc(obj interface{}) (string, error) {
	return obj.(testFifoObject).name, nil
}

func mkFifoObj(name string, val interface{}) testFifoObject {
	return testFifoObject{name: name, val: val}
}

func TestFIFO_requeueOnPop(t *testing.T) {
	f := NewFIFO(testFifoObjectKeyFunc)
	var testData = []testFifoObject{
		{name: "test", val: 10},
		{name: "test2", val: 10},
		{name: "test3", val: 10},
		{name: "tes4", val: 10},
	}
	for _, test := range testData {
		f.Add(test)
		time.Sleep(3 * time.Second)
	}
	for i := 0; i < len(testData); i++ {
		data := Pop(f)
		fmt.Println(data)
	}
	item, ok, err := f.Get(testFifoObject{name: "test", val: 10})
	fmt.Println(item, ok, err)
}
