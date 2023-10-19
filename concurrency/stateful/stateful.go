package stateful

import (
	//"fmt"
	"sync"
)

var ops int

func BalanceState(deposits chan int, balanceChecks chan int, withdrawals chan int, wg *sync.WaitGroup) {
	balance := 0


	for {
		select {
		case deposit := <-deposits:
			balance += deposit
			ops++
			//fmt.Printf(" - Depositted %d successfully. Balance is now %d\n", deposit, balance)
			wg.Done()
		case <-balanceChecks:
			//fmt.Printf(" - Balance is currently %d\n", balance)
			ops++
			wg.Done()
		case withdrawal := <-withdrawals:
			ops++
			if withdrawal > balance {
				//fmt.Printf(" - Unable to withdraw %d. Current balance is %d\n", withdrawal, balance)
			} else {
				balance -= withdrawal
				//fmt.Printf(" - Withdrawn %d successfully. Balance is now %d\n", withdrawal, balance)
			}
			wg.Done()
		}
	}
}

func GetOps() int {
	return ops
}
