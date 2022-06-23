package main

import "fmt"

func main() {
	a := 1
	var b int = 0
	for i := 0; i < a; i++ {
		b++
	}
	a = 2

	strList := []string{"hello", "world", "world"}
	for _, str := range strList {
		fmt.Println(str)
	}
}
