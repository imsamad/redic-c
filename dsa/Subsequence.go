package main

import (
	"fmt"
)

func PrintSub(str string, current int, sub string) {
	if len(str) == current {
		if len(sub) == 0 {
			return
		}
		fmt.Println(sub)
		return
	}

	PrintSub(str, current+1, sub)
	PrintSub(str, current+1, fmt.Sprintf("%s%s", sub, string(str[current])))
}
