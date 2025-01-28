package main

import (
	b "blockman/block"
)

func main() {
	bc := b.NewBlockChain()
	bc.Print()

	bc.AddTransaction("A", "B", 1.0)
	previousHash := bc.Last().Hash()
	bc.Create(5, previousHash)
	bc.Print()

	bc.AddTransaction("C", "D", 2.0)
	bc.AddTransaction("C", "D", 3.0)
	previousHash = bc.Last().Hash()
	bc.Create(2, previousHash)
	bc.Print()
}
