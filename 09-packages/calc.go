package main

import (
	"fmt"
	"package-example/sum"
	"package-example/diff"
)

func main(){

	fmt.Println(sum.Sum(2,3))
	fmt.Println(diff.Diff(5,3))
}