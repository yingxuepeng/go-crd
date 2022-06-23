package main

import "fmt"

type struct1 struct {
	s1var1 string
	s1var2 bool
}

type struct2 struct {
	s2var1 map[int]struct1
}
type struct3 struct {
}

type interface1 interface {
	if1m1(if1m1v1 string, if1m1v2 int) (if1m1r1 string, if1m1r2 string)
	if1m2(if1m2v1 string) struct2
	if1m3(if1m2v1 string) (if1m3r1 struct3)
}

func func1(f1v1 string, f1v2 byte) (f1r1 int, f1r2 int) {
	return 1, 2
}
func main() {
	var1 := 1
	var var2 int = 0
	for i := 0; i < var1; i++ {
		var2++
	}
	var1 = 2

	if var3 := 123; var3 < 1 {
		var4 := 4
		var1 = var4
	}

	for var5 := 0; var5 < var2; var5++ {
	}

	strList := []string{"hello", "world", "world"}
	for _, str := range strList {
		fmt.Println(str)
	}
}
