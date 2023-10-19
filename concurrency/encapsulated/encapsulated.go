package main

import (
	"fmt"
	"time"
)

type operation struct {
	action   string
	amount   int
	response chan string
}

var requests chan operation = make(chan operation, 100)
var done chan struct{} = make(chan struct{})
var balance int
var started time.Time

func Deposit(amount int) string {
	answer := make(chan string)
	op := operation{
		action:   "deposit",
		amount:   amount,
		response: answer,
	}
	requests <- op
	return <-answer
}

func Withdraw(amount int) string {
	answer := make(chan string)
	op := operation{
		action:   "withdraw",
		amount:   amount,
		response: answer,
	}
	requests <- op
	return <-answer
}
func GetBalance() string {
	answer := make(chan string)
	op := operation{
		action:   "balance",
		amount:   0,
		response: answer,
	}
	requests <- op
	return <-answer
}

func Start() {
	started = time.Now()
	go monitorRequests()
}

func Stop() {
	shutdown := operation{action: "shutdown", amount: 0, response: nil}
	requests <- shutdown
	<-done
}

func monitorRequests() {
	time.Sleep(time.Second)
	for op := range requests {
		switch op.action {
		case "deposit":
			balance += op.amount
			newresponse := fmt.Sprintf("Depositted %d successfully. Balance is now %d", op.amount, balance)
			op.response <- newresponse
		case "withdraw":
			if op.amount > balance {
				newresponse := fmt.Sprintf("Insufficient funds to withdraw %d. Balance remains at %d", op.amount, balance)
				op.response <- newresponse
			} else {
				balance -= op.amount
				newresponse := fmt.Sprintf("Withdrew %d successfully. Balance is now %d", op.amount, balance)
				op.response <- newresponse
			}
		case "balance":
			newresponse := fmt.Sprintf("Current balance is %d", balance)
			op.response <- newresponse
		case "shutdown":
			fmt.Println("Shutting down")
			close(requests)
		}


	}
	fmt.Printf("All requests processed in %v", time.Since(started))
	close(done)
}

func main() {
	Start()
	defer Stop()
	for i := 0; i < 100; i++ {
		go func() {
			response := Deposit(50)
			fmt.Println(response)
		}()
		go func() {
			response := Withdraw(60)
			fmt.Println(response)
		}()
		go func() {
			response := GetBalance()
			fmt.Println(response)
		}()
	}
	time.Sleep(time.Second)
}
