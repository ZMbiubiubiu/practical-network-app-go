package main

import "fmt"

func main() {
	// for range 会copy一份range后的expression，对于遍历array来说，可能会有坑
	// The reason is that participating in the for range loop is a copy of the range expression.
	// That is, in the above example, it is the copy of a that is actually participating in the loop, not the real a

	// for array
	var a = [5]int{1, 2, 3, 4, 5}
	var b [5]int

	fmt.Println("original a =", a) // [1 2 3 4 5]

	for i, v := range a {
		if i == 0 {
			a[1] = 12
			a[2] = 13
		}
		b[i] = v
	}

	fmt.Println("after for range loop, r =", b) // [1 2 3 4 5]
	fmt.Println("after for range loop, a =", a) // [1 12 13 4 5]

	// for slice
	var c = []int{1, 2, 3, 4, 5}
	var d []int = make([]int, len(c))

	fmt.Println("original c =", c) // [1 2 3 4 5]

	for i, v := range c {
		if i == 0 {
			c[1] = 12
			c[2] = 13
		}
		d[i] = v
	}

	fmt.Println("after for range loop, d =", d) // [1 2 3 4 5]
	fmt.Println("after for range loop, c =", c) // [1 12 13 4 5]
}
