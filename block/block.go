package block

import (
	t "blockman/transaction"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

type Block struct {
	Nonce        int              `json:"nonce"`
	PreviousHash [32]byte         `json:"previous_hash"`
	Timestamp    int64            `json:"timestamp"`
	Transactions []*t.Transaction `json:"transactions"`
}

type BlockChain struct {
	transactionPool []*t.Transaction
	chain           []*Block
}

func NewBlock(nonce int, previousHash [32]byte, transactions []*t.Transaction) *Block {
	return &Block{
		Timestamp:    time.Now().UnixNano(),
		Nonce:        nonce,
		PreviousHash: previousHash,
		Transactions: transactions,
	}
}

func NewBlockChain() *BlockChain {
	b := &Block{}
	bc := &BlockChain{}
	bc.Create(0, b.Hash())
	return bc
}

func (b *Block) Print() {
	fmt.Printf("timestamp     %d\n", b.Timestamp)
	fmt.Printf("nonce     %d\n", b.Nonce)
	fmt.Printf("previous_hash     %x\n", b.PreviousHash)
	for _, transaction := range b.Transactions {
		transaction.Print()
	}
}

func (b *Block) Hash() [32]byte {
	m, err := json.Marshal(b)
	if err != nil {
		log.Fatal("Error while running Hash function : ", err)
		return [32]byte{}
	}
	fmt.Printf(string(m))
	return sha256.Sum256(m)
}

func (bc *BlockChain) Print() {
	for i, block := range bc.chain {
		fmt.Printf("%s Chain %d %s\n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 25))
}

func (bc *BlockChain) Create(nonce int, previousHash [32]byte) *Block {
	b := NewBlock(nonce, previousHash, bc.transactionPool)
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*t.Transaction{}
	return b
}

func (bc *BlockChain) Last() *Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *BlockChain) AddTransaction(sender string, recipient string, value float32) {
	t := t.NewTransaction(sender, recipient, value)
	bc.transactionPool = append(bc.transactionPool, t)
}
