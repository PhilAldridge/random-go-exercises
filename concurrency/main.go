package main

import (
	"concurrency/mutex"
	"concurrency/stateful"
	"fmt"
	"sync"
	"time"
)

func main() {
	started := time.Now()
	fmt.Print("This bank uses Mutex to avoid race conditions:\n")
	var wg sync.WaitGroup

	for i := 0; i < 10000; i++ {
		wg.Add(3)
		go mutex.Deposit(50, &wg)
		go mutex.Withdraw(60, &wg)
		go mutex.Balance(&wg)
	}
	wg.Wait()
	fmt.Printf("%v taken to complete %d operations\n", time.Since(started), mutex.GetOps())

	started = time.Now()
	fmt.Print("\n\n\nWhilst this bank uses a single stateful goRoutine:\n")

	deposits := make(chan int)
	balanceChecks := make(chan int)
	withdrawals := make(chan int)

	go stateful.BalanceState(deposits, balanceChecks, withdrawals, &wg)

	for i := 0; i < 10000; i++ {
		wg.Add(3)
		go func() {
			deposits <- 50
		}()
		go func() {
			withdrawals <- 60
		}()
		go func() {
			balanceChecks <- 0
		}()
	}
	wg.Wait()
	fmt.Printf("%v taken to complete %d operations\n", time.Since(started), stateful.GetOps())
	close(deposits)
	close(balanceChecks)
	close(withdrawals)
}
