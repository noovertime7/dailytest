package main

import (
	"fmt"
	"github.com/noovertime7/dailytest/scheme_demo/product"
	"github.com/noovertime7/dailytest/scheme_demo/scheme"
)

var sh = scheme.NewScheme()

var localSchemeBuilder = scheme.SchemeBuilder{
	product.AddToScheme,
}

var AddToScheme = localSchemeBuilder.AddScheme

func main() {
	err := AddToScheme(sh)
	if err != nil {
		return
	}
	var gvk = scheme.GroupVersionKind{Group: "food", Version: "v1"}
	res, err := sh.GetName(gvk)
	Must(err)
	fmt.Println(res)
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}
