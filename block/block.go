package block

import (
	s "blockman/signature"
	t "blockman/transaction"
	"blockman/utils"
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	MINING_DIFFICULTY = 3
	MINING_SENDER     = "THE BLOCKCHAIN"
	MINING_REWARD     = 1.0
	MINING_TIMER_SEC  = 20

	BLOCKCHAIN_PORT_RANGE_START       = 5000
	BLOCKCHAIN_PORT_RANGE_END         = 5003
	NEIGHBOR_IP_RANGE_START           = 0
	NEIGHBOR_IP_RANGE_END             = 1
	BLOCKCHAIN_NEIGHBOR_SYNC_TIME_SEC = 20
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
	port              uint16
	mux               sync.Mutex

	neighbors    []string
	muxNeighbors sync.Mutex
}

func NewBlock(nonce int, previousHash [32]byte, transactions []*t.Transaction) *Block {
	return &Block{
		Timestamp:    time.Now().UnixNano(),
		Nonce:        nonce,
		PreviousHash: previousHash,
		Transactions: transactions,
	}
}

func (bc *BlockChain) GetChain() []*Block {
	return bc.chain
}

func (b *Block) GetPreviousHash() [32]byte {
	return b.PreviousHash
}

func (b *Block) GetNonce() int {
	return b.Nonce
}

func (b *Block) GetTransactions() []*t.Transaction {
	return b.Transactions
}

func NewBlockChain(blockChainAddress string, port uint16) *BlockChain {
	b := &Block{}
	bc := &BlockChain{blockChainAddress: blockChainAddress, port: port}
	bc.Create(0, b.Hash())
	return bc
}

func (bc *BlockChain) TransactionPool() []*t.Transaction {
	return bc.transactionPool
}

func (bc *BlockChain) Run() {
	bc.StartSyncNeighbors()
	bc.ResolveConflicts()
	bc.StartMining()
}

func (bc *BlockChain) SetNeighbors() {
	bc.neighbors = utils.FindNeighbors("127.0.0.1", bc.port, NEIGHBOR_IP_RANGE_START, NEIGHBOR_IP_RANGE_END, BLOCKCHAIN_PORT_RANGE_START, BLOCKCHAIN_PORT_RANGE_END)
	log.Printf("%v", bc.neighbors)
}

func (bc *BlockChain) SyncNeighbors() {
	bc.muxNeighbors.Lock()
	defer bc.muxNeighbors.Unlock()
	bc.SetNeighbors()
}

func (bc *BlockChain) StartSyncNeighbors() {
	bc.SyncNeighbors()
	_ = time.AfterFunc(time.Second*BLOCKCHAIN_NEIGHBOR_SYNC_TIME_SEC, bc.StartSyncNeighbors)
}

func (bc *BlockChain) ClearTransactionPoo() {
	bc.transactionPool = bc.transactionPool[:0]
}

func (bc *BlockChain) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Blocks []*Block `json:"chains"`
	}{
		Blocks: bc.chain,
	})
}

func (bc *BlockChain) UnmarshalJSON(data []byte) error {
	v := struct {
		Blocks *[]*Block `json:"chains"`
	}{
		Blocks: &bc.chain,
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	return nil
}

func (b *Block) Print() {
	fmt.Printf("timestamp     %d\n", b.Timestamp)
	fmt.Printf("nonce     %d\n", b.Nonce)
	fmt.Printf("previous_hash     %x\n", b.PreviousHash)
	for _, transaction := range b.Transactions {
		transaction.Print()
	}
}

func (b *Block) UnmarshalJSON(data []byte) error {
	v := &struct {
		Timestamp    *int64            `json:"timestamp"`
		Nonce        *int              `json:"none"`
		PreviousHash *string           `json:"previous_hash"`
		Transaction  *[]*t.Transaction `json:"transactions"`
	}{
		Timestamp:    &b.Timestamp,
		Nonce:        &b.Nonce,
		PreviousHash: nil, //for now
		Transaction:  &b.Transactions,
	}

	if err := json.Unmarshal(data, v); err != nil {
		return err
	}

	ph, _ := hex.DecodeString(*v.PreviousHash)
	copy(b.PreviousHash[:], ph[:32])
	return nil
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
	for _, n := range bc.neighbors {
		endpoint := fmt.Sprintf("http://%s/transactions", n)
		client := &http.Client{}
		req, _ := http.NewRequest("DELETE", endpoint, nil)
		resp, _ := client.Do(req)
		log.Printf("%v", resp)
	}
	return b
}

func (bc *BlockChain) Last() *Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *BlockChain) CreateTransaction(senderPublicKey *ecdsa.PublicKey, s *s.Signature, sender string, recipient string, value float32) bool {
	isTransacted := bc.AddTransaction(senderPublicKey, s, sender, recipient, value)

	if isTransacted {
		for _, n := range bc.neighbors {
			publicKeyStr := fmt.Sprintf("%064x%064x", senderPublicKey.X.Bytes(), senderPublicKey.Y.Bytes())
			signatureStr := s.String()
			bt := &TransactionRequest{
				&sender, &recipient, &publicKeyStr, &value, &signatureStr}
			m, _ := json.Marshal(bt)
			buf := bytes.NewBuffer(m)
			endpoint := fmt.Sprintf("http://%s/transactions", n)
			client := &http.Client{}
			req, _ := http.NewRequest("PUT", endpoint, buf)
			resp, _ := client.Do(req)
			log.Printf("%v", resp)
		}
	}

	return isTransacted
}

func (bc *BlockChain) AddTransaction(senderPublicKey *ecdsa.PublicKey, s *s.Signature, sender string, recipient string, value float32) bool {
	t := NewTransaction(sender, recipient, value)

	if sender == MINING_SENDER {
		bc.transactionPool = append(bc.transactionPool, t)
		return true
	}

	if bc.VerifyTransactionSignature(senderPublicKey, s, t) {
		if bc.CalculateTotalAmount(sender) < value {
			log.Println("ERROR: Not enough balance in a wallet")
			return false
		}
		bc.transactionPool = append(bc.transactionPool, t)
		return true
	} else {
		log.Println("ERROR : Verify Transaction")
	}
	return false
}

func (bc *BlockChain) VerifyTransactionSignature(senderPublicKey *ecdsa.PublicKey, s *s.Signature, t *t.Transaction) bool {
	m, _ := t.MarshalJSON()
	h := sha256.Sum256([]byte(m))
	return ecdsa.Verify(senderPublicKey, h[:], s.R, s.S)
}

func (bc *BlockChain) CopyTransactionPool() []*t.Transaction {
	transactions := []*t.Transaction{}
	for _, transaction := range bc.transactionPool {
		transactions = append(transactions,
			t.NewTransaction(nil, nil, transaction.SenderBlockchainAddress, transaction.RecipientBlockchainAddress, transaction.Value),
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
	bc.mux.Lock()
	defer bc.mux.Unlock()

	// if len(bc.transactionPool) == 0 {
	// 	return false
	// }

	bc.AddTransaction(nil, nil, MINING_SENDER, bc.blockChainAddress, MINING_REWARD)
	nonce := bc.ProofOfWork()
	previousHash := bc.Last().Hash()
	bc.Create(nonce, previousHash)
	log.Println("action=mining, status success")

	for _, n := range bc.neighbors {
		endpoint := fmt.Sprintf("http://%s/concensus", n)
		client := &http.Client{}
		req, _ := http.NewRequest("PUT", endpoint, nil)
		resp, _ := client.Do(req)
		log.Printf("%v", resp)
	}

	return true
}

func (bc *BlockChain) StartMining() {
	bc.Mining()
	_ = time.AfterFunc(time.Second*MINING_TIMER_SEC, bc.StartMining)
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

func (bc *BlockChain) ValidChain(chain []*Block) bool {
	preBlock := chain[0]
	currentIdx := 1
	for currentIdx < len(chain) {
		b := chain[currentIdx]
		if b.PreviousHash != preBlock.Hash() {
			return false
		}

		if !bc.ValidProof(b.Nonce, b.GetPreviousHash(), b.GetTransactions(), MINING_DIFFICULTY) {
			return false
		}

		preBlock = b
		currentIdx++
	}

	return true
}

func (bc *BlockChain) ResolveConflicts() bool {
	var longestChain []*Block = nil
	maxLength := len(bc.chain)

	for _, n := range bc.neighbors {
		endpoint := fmt.Sprintf("http://%s/chain", n)
		resp, _ := http.Get(endpoint)
		if resp.StatusCode == 200 {
			var bcResp BlockChain
			decoder := json.NewDecoder(resp.Body)
			_ = decoder.Decode(&bcResp)

			chain := bcResp.GetChain()
			if len(chain) > maxLength && bc.ValidChain(chain) {
				maxLength = len(chain)
				longestChain = chain
			}
		}
	}
	if longestChain != nil {
		bc.chain = longestChain
		log.Printf("Resolve Conflicts Replaced")
		return true
	}
	log.Printf("Resolve Conflicts Not Replaced")
	return false
}

func NewTransaction(sender string, recipient string, value float32) *t.Transaction {
	return &t.Transaction{SenderBlockchainAddress: sender, RecipientBlockchainAddress: recipient, Value: value}
}

type TransactionRequest struct {
	SenderBlockchainAddress    *string  `json:"sender_blockchain_address"`
	RecipientBlockchainAddress *string  `json:"recipient_blockchain_address"`
	SenderPublicKey            *string  `json:"sender_public_key"`
	Value                      *float32 `json:"value"`
	Signature                  *string  `json:"signature"`
}

func (tr *TransactionRequest) Validate() bool {
	if tr.SenderPublicKey == nil ||
		tr.RecipientBlockchainAddress == nil ||
		tr.SenderBlockchainAddress == nil ||
		tr.Value == nil ||
		tr.Signature == nil {
		return false
	}
	return true
}

type AmountResponse struct {
	Amount float32 `json:"amount"`
}

func (ar *AmountResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Amount float32 `json"amount"`
	}{
		Amount: ar.Amount,
	})
}
