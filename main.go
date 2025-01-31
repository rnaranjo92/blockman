package main

import (
	"blockman/wallet"
	"fmt"
)

func main() {
	// myBCAddress := "my_blockchain_address"
	// bc := b.NewBlockChain(myBCAddress)
	// bc.Print()

	// bc.AddTransaction("A", "B", 1.0)
	// bc.Mining()
	// bc.Print()

	// bc.AddTransaction("C", "D", 2.0)
	// bc.AddTransaction("X", "Y", 3.0)
	// bc.Mining()
	// bc.Print()

	// fmt.Printf("my_blockchain_address %.1f\n", bc.CalculateTotalAmount("my_blockchain_address"))
	// fmt.Printf("C %.1f\n", bc.CalculateTotalAmount("C"))
	// fmt.Printf("A %.1f\n", bc.CalculateTotalAmount("A"))
	// fmt.Printf("D %.1f\n", bc.CalculateTotalAmount("D"))

	w := wallet.NewWallet()
	fmt.Println(w.GetPrivateKeyStr())
	fmt.Println(w.GetPublicKeyStr())
	fmt.Println(w.GetBlockChainAddress())
}
