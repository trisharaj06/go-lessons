package main

import "fmt"

var val float32

func main() {
    var a = "initial"
    fmt.Println(a)

    var b, c int = 1, 2
    fmt.Println(b, c)

    var d = true
    fmt.Println(d)

    var e int
    fmt.Println(e)

    f := "apple"
    fmt.Println(f)

	val = 3.14
	fmt.Println(val)
}