package bank

import (
	"fmt"
)

var deposits = make(chan int)
var balances = make(chan int)
var withdraws = make(chan withdraw)

type withdraw struct {
	amount int
	notice chan bool
}

func Deposit(amount int) {
	deposits <- amount
}

func Balance() int {
	return <-balances
}

func teller() {
	var balance int //confined to the monitor routine
	for {
		select {
		case withdraw := <- withdraws:
			t := balance-withdraw.amount
			if t >= 0 {
				balance = t
				withdraw.notice <- true
			}else {
				withdraw.notice <- false
			}
		case amount := <-deposits:
			balance += amount
		case balances <- balance:
		}
	}
}

func init() {
	go teller()
}

func Withdraw(amount int) bool {
	notifier := make(chan bool)
	withdraws <- withdraw{ amount: amount, notice: notifier}
	for {
		select {
		case ok := <-notifier:
			if !ok {
				fmt.Println("Insufficient funds")
			}
			return ok
		}
	}
}