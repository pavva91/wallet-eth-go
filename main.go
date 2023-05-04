package main

import (
	"address/wallet"
	"fmt"
)

func main() {
	wallet.CurrentBlock()
	// wallet.CreateWallet()
	senderAddress := "0x941907d6C5Ec6f7dBB1B91F67752D4c127eE6f87"
	receiverAddress := "0x1Dabfed3934E7ab52ebeAbdd9153fe0003ee648f"
    fmt.Println("SENDER BALANCE")
	wallet.CheckBalance(senderAddress)
    fmt.Println("RECEIVER BALANCE")
	wallet.CheckBalance(receiverAddress)
    wallet.QueryTransactions()
    wallet.GetTransactionsPerBlock()
	wallet.CreateTransactionWrap()
	// time.Sleep(5 * time.Second)
	// wallet.CheckBalance(senderAddress)
}
