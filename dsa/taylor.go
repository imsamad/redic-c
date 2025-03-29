package main

func Taylor(n int, x float64) float64 {
	var a float64
	a = 1
	b := 1

	var handler func(n int) float64

	handler = func(n int) float64 {
		if n == 0 {
			return 1
		}

		r := handler(n-1)
		a = a * x
		b = b * n

		return r + a/float64(b)
	}

	return handler(n)
}

func TaylorHornet(n int, x float64) float64 {
	var S float64
	S = 1

	var handle func(n int) float64


	handle = func(n int) float64 {
		if n == 0 {
			return S
		}
		S = 1 + (x / float64(n) * S)
		return handle(n-1)
	}

	return handle(n)
}
