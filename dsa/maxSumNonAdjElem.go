package main

func MaxSumNonAdjElem(arr []int, crtIdx int, sum int, res *[]int) {
	// fmt.Print(sum, "")
	if crtIdx == len(arr) {
		*res = append(*res, sum)
		return
	}

	MaxSumNonAdjElem(arr, crtIdx+1, sum+arr[crtIdx], res)
	MaxSumNonAdjElem(arr, crtIdx+2, sum, res)
}
