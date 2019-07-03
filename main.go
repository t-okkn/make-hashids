package main

import (
	"fmt"
)

func main() {
	input := "ほげほげ"
	res := Str2Uints(input)
	fmt.Println(res)

	 h, salt := CreateHashids(res)
	fmt.Println(h)
	fmt.Println(salt)
}

