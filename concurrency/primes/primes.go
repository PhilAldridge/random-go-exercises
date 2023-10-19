package main

import (
	"fmt"
	"sync"
	"time"
)

func isPrime(number int) bool {
	for i := 2; i < number/2; i++ {
		if number%i == 0 {
			return false
		}
	}
	return true
}

func maxInts(x int, y int) int {
	if x > y {
		return x
	}
	return y
}

func main() {
	started := time.Now()
	var mu sync.Mutex
	primes := 0
	largestPrime := 0

	var wg sync.WaitGroup
	for i := 2; i < 1000000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if isPrime(i) {
				mu.Lock()
				primes++
				largestPrime = maxInts(i, largestPrime)
				mu.Unlock()
			}
		}(i)

	}
	wg.Wait()
	fmt.Printf("%d primes found in %v (with goroutines)\nThe largest prime is %d\n", primes, time.Since(started), largestPrime)

	started = time.Now()
	largestPrime = 0
	primes = 2

	for i := 6; i < 1000000; i+=6 {
		if isPrime(i-1) {
			primes++
			largestPrime = i-1
		}
		if isPrime(i+1) {
			primes++
			largestPrime = i+1
		}
	}
	fmt.Printf("\n%d primes found in %v (without using goroutines)\nThe largest prime is %d\n", primes, time.Since(started), largestPrime)
}
