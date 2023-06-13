package main

import (
	"fmt"
	"regexp"
)

func main() {
	checkIsDeleteCluster := func(in, clusterName string) bool {
		re := regexp.MustCompile("^[^/]*")
		match := re.FindString(in)
		fmt.Println(match)
		return match == clusterName
	}
	fmt.Println(checkIsDeleteCluster("cloud/default", "cloud"))
}
