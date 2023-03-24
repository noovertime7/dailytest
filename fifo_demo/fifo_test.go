package fifo_demo

import (
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

func TestFIFO_requeueOnPop(t *testing.T) {
	f := NewFIFO(testFifoObjectKeyFunc)
	var testData = []testFifoObject{
		{name: "test", val: 10},
	}
	for _, test := range testData {
		f.Add(test)
		time.Sleep(3 * time.Second)
	}
	for i := 0; i < len(testData); i++ {
		t.Logf("keys : %v", f.ListKeys())
		data := Pop(f)
		t.Log(data)
	}
	item, ok, err := f.Get(testFifoObject{name: "test", val: 10})
	t.Log(item, ok, err)
}
