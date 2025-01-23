package main

import (
	b "blockman/types"
)

func main() {
	bc := b.NewBlockChain()
	bc.Create(5, "hash 1")
	bc.Create(2, "hash 2")
	bc.Print()
}
