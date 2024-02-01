package exampletasks

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/RichardKnop/machinery/v2/log"
)

// Add ...
func Add(args ...int64) (int64, error) {
	sum := int64(0)
	for _, arg := range args {
		sum += arg
	}
	return sum, nil
}

// Multiply ...
func Multiply(args ...int64) (int64, error) {
	sum := int64(1)
	for _, arg := range args {
		sum *= arg
	}
	return sum, nil
}

// SumInts ...
func SumInts(numbers []int64) (int64, error) {
	var sum int64
	for _, num := range numbers {
		sum += num
	}
	return sum, nil
}

// SumFloats ...
func SumFloats(numbers []float64) (float64, error) {
	var sum float64
	for _, num := range numbers {
		sum += num
	}
	return sum, nil
}

// Concat ...
func Concat(strs []string) (string, error) {
	var res string
	for _, s := range strs {
		res += s
	}
	return res, nil
}

type Demo struct {
	In   []string
	Name string
}

// Split ...
func Split(str []string, name string) (string, error) {
	fmt.Println("被调用")
	fmt.Println("接收到参数", str)
	data, _ := json.Marshal(Demo{
		In:   str,
		Name: name,
	})

	return string(data), nil
}

// PanicTask ...
func PanicTask() (string, error) {
	panic(errors.New("oops"))
}

// LongRunningTask ...
func LongRunningTask(str string) (string, error) {
	log.INFO.Print("Long running task started")
	fmt.Println("str args = ", str)
	for i := 0; i < 10; i++ {
		log.INFO.Print(10 - i)
		time.Sleep(1 * time.Second)
	}
	log.INFO.Print("Long running task finished")
	return str, nil
}
