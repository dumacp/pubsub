package main

import (
	"fmt"
	"math"
)

func IsPrime(value int64) bool {
	for i := int64(2); i <= int64(math.Floor(float64(value)/2)); i++ {
		if value%i == 0 {
			return false
		}
	}
	return value > 1
}

func IsPrimeSqrt(value int64) bool {
	for i := int64(2); i <= int64(math.Floor(math.Sqrt(float64(value)))); i++ {
		if value%i == 0 {
			return false
		}
	}
	return value > 1
}

func SieveOfEratosthenes(value int) {
	f := make([]bool, value)
	for i := 2; i <= int(math.Sqrt(float64(value))); i++ {
		if f[i] == false {
			for j := i * i; j < value; j += i {
				f[j] = true
			}
		}
	}
	for i := 2; i < value; i++ {
		if f[i] == false {
			fmt.Printf("%v ", i)
		}
	}
	fmt.Println("")
}

func main() {
	/**
	for i := 1; i <= 100; i++ {
		if IsPrime(i) {
			fmt.Printf("%v ", i)
		}
	}
	fmt.Println("")
	/**/
	chs := make([]chan int64, 2)

	chs[0] = make(chan int64)
	chs[1] = make(chan int64)
	go func() {
		for v := range chs[0] {
			if IsPrimeSqrt(v) {
				fmt.Printf("%v ", v)
			}
		}
	}()
	go func() {
		for v := range chs[1] {
			if IsPrimeSqrt(v) {
				fmt.Printf("%v ", v)
			}
		}
	}()

	for i := int64(0); i < int64(100); i++ {

		select {
		case chs[0] <- i:
		case chs[1] <- i:
		}
	}
	close(chs[0])
	close(chs[1])
	fmt.Println("")
	/**
	SieveOfEratosthenes(100)
	/**/
}
