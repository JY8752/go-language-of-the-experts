package testdata

import "fmt"

func SimpleFunction(n int) {
	fmt.Println(n)
}

func ComplexFunction(n int) { // want `func ComplexFunction complexity=10`
	if n == 0 {
		println("zero")
	}
	if n == 1 {
		println("One")
	}
	if n == 2 {
		println("two")
	}
	if n == 3 {
		println("three")
	}
	if n == 4 {
		println("four")
	}
	if n == 5 {
		println("five")
	}
	if n == 6 {
		println("six")
	}
	if n == 7 {
		println("seven")
	}
	if n == 8 {
		println("eight")
	}
	if n == 9 {
		println("nine")
	}
	println("not between zero and eight")
}
