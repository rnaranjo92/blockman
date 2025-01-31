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

const (
	MINING_DIFFICULTY = 3
	MINING_SENDER     = "THE BLOCKCHAIN"
	MINING_REWARD     = 1.0
)

type Block struct {
	Nonce        int              `json:"nonce"`
	PreviousHash [32]byte         `json:"previous_hash"`
	Timestamp    int64            `json:"timestamp"`
	Transactions []*t.Transaction `json:"transactions"`
}

type BlockChain struct {
	transactionPool   []*t.Transaction
	chain             []*Block
	blockChainAddress string
}

func NewBlock(nonce int, previousHash [32]byte, transactions []*t.Transaction) *Block {
	return &Block{
		Timestamp:    time.Now().UnixNano(),
		Nonce:        nonce,
		PreviousHash: previousHash,
		Transactions: transactions,
	}
}

func NewBlockChain(blockChainAddress string) *BlockChain {
	b := &Block{}
	bc := &BlockChain{blockChainAddress: blockChainAddress}
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

func (bc *BlockChain) CopyTransactionPool() []*t.Transaction {
	transactions := []*t.Transaction{}
	for _, transaction := range bc.transactionPool {
		transactions = append(transactions,
			t.NewTransaction(transaction.SenderBlockchainAddress, transaction.RecipientBlockchainAddress, transaction.Value),
		)
	}
	return transactions
}

func (bc *BlockChain) ValidProof(nonce int, previousHash [32]byte, transactions []*t.Transaction, difficulty int) bool {
	zeros := strings.Repeat("0", difficulty)
	guessBlock := Block{Timestamp: 0, Nonce: nonce, Transactions: transactions}
	guessHash := fmt.Sprintf("%x", guessBlock.Hash())

	return guessHash[:difficulty] == zeros
}

func (bc *BlockChain) ProofOfWork() int {
	transactions := bc.CopyTransactionPool()
	previousHash := bc.Last().Hash()
	nonce := 0
	for !bc.ValidProof(nonce, previousHash, transactions, MINING_DIFFICULTY) {
		nonce++
	}

	return nonce
}

func (bc *BlockChain) Mining() bool {
	bc.AddTransaction(MINING_SENDER, bc.blockChainAddress, MINING_REWARD)
	nonce := bc.ProofOfWork()
	previousHash := bc.Last().Hash()
	bc.Create(nonce, previousHash)
	log.Println("action=mining, status success")
	return true
}

func (bc *BlockChain) CalculateTotalAmount(blockchainAddress string) float32 {
	var totalAmount float32 = 0.0

	for _, b := range bc.chain {
		for _, t := range b.Transactions {
			if blockchainAddress == t.RecipientBlockchainAddress {
				totalAmount += t.Value
			}

			if blockchainAddress == t.SenderBlockchainAddress {
				totalAmount -= t.Value
			}
		}
	}

	return totalAmount
}
