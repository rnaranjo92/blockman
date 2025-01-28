package transaction

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Transaction struct {
	SenderBlockchainAddress    string  `json:"sender_blockchain_address"`
	RecipientBlockchainAddress string  `json:"recipient_blockchain_address"`
	Value                      float32 `json:"value"`
}

func NewTransaction(sender string, recipient string, value float32) *Transaction {
	return &Transaction{sender, recipient, value}
}

func (t *Transaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 40))
	fmt.Printf(" sender_blockchain_address  %s\n", t.SenderBlockchainAddress)
	fmt.Printf(" recipient_blockchain_address  %s\n", t.RecipientBlockchainAddress)
	fmt.Printf(" value  %.1f\n", t.Value)
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(*t)
}
