package mutex

import (
	//"fmt"
	"sync"
	"sync/atomic"
)

var mu sync.Mutex
var balance int

var ops atomic.Uint64

func deposit(amount int) {
	balance += amount
}

func Deposit(amount int, wg *sync.WaitGroup) {
	ops.Add(1)
	mu.Lock()
	defer mu.Unlock()
	defer wg.Done()
	deposit(amount)
	//fmt.Printf(" - %d depositted successfully. New balance is %d\n", amount, balance)
}

func Withdraw(amount int, wg *sync.WaitGroup) bool {
	ops.Add(1)
	mu.Lock()
	defer mu.Unlock()
	defer wg.Done()
	deposit(-amount)
	
	if balance < 0 {
		deposit(amount)
		//fmt.Printf(" - Withdrawal of %d was unsuccessfull. Balance remains %d\n", amount, balance)
		return false
	}
	//fmt.Printf(" - %d withdrawn successfully. New balance is %d\n", amount, balance)
	return true
}

func Balance(wg *sync.WaitGroup) int {
	ops.Add(1)
	mu.Lock()
	defer mu.Unlock()
	defer wg.Done()
	//fmt.Printf(" - Balance is %d\n", balance)
	return balance
}

func GetOps() int {
	return int(ops.Load())
}
