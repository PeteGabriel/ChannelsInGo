package bank

var balance int

func DepositUnsafe(amount int){
	balance = balance + amount
}

func BalanceUnsafe() int {
	return balance
}



