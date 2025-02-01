package main

import (
	"blockman/block"
	transaction "blockman/transaction"
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

	walletM := wallet.NewWallet()
	walletA := wallet.NewWallet()
	walletB := wallet.NewWallet()

	t := transaction.NewTransaction(walletA.PrivateKey, walletA.PublicKey, walletA.BlockChainAddress, walletB.BlockChainAddress, 1.0)

	blockchain := block.NewBlockChain(walletM.BlockChainAddress)
	isAdded := blockchain.AddTransaction(walletA.PublicKey, t.GenerateSignature(), walletA.BlockChainAddress, walletB.BlockChainAddress, 1.0)
	fmt.Println("Added? ", isAdded)

	blockchain.Mining()
	blockchain.Print()

	fmt.Printf("A %.1f\n", blockchain.CalculateTotalAmount(walletA.BlockChainAddress))
	fmt.Printf("B %.1f\n", blockchain.CalculateTotalAmount(walletB.BlockChainAddress))
	fmt.Printf("M %.1f\n", blockchain.CalculateTotalAmount(walletM.BlockChainAddress))
}
