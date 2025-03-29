package main

import "fmt"

func PascalTriangle(n int) {
	for i := 0; i <= n; i++ {
		x := 1
		fmt.Print(x, " ")
		for j := 1; j <= i; j++ {
			x = (x * (i - j + 1)) / j
			fmt.Print(x, " ")
		}
		fmt.Println()
	}
}
