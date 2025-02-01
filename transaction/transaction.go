package transaction

import (
	sig "blockman/signature"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
)

type Transaction struct {
	SenderPrivateKey           *ecdsa.PrivateKey `json:"sender_private_key"`
	SenderPublicKey            *ecdsa.PublicKey  `json:"sender_public_key"`
	SenderBlockchainAddress    string            `json:"sender_blockchain_address"`
	RecipientBlockchainAddress string            `json:"recipient_blockchain_address"`
	Value                      float32           `json:"value"`
}

func (t *Transaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 40))
	fmt.Printf(" sender_blockchain_address  %s\n", t.SenderBlockchainAddress)
	fmt.Printf(" recipient_blockchain_address  %s\n", t.RecipientBlockchainAddress)
	fmt.Printf(" value  %.1f\n", t.Value)
}

func NewTransaction(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey, sender string, recipient string, value float32) *Transaction {
	return &Transaction{privateKey, publicKey, sender, recipient, value}
}

func (t *Transaction) GenerateSignature() *sig.Signature {
	m, _ := t.MarshalJSON()
	h := sha256.Sum256([]byte(m))
	r, s, _ := ecdsa.Sign(rand.Reader, t.SenderPrivateKey, h[:])
	return &sig.Signature{R: r, S: s}
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string  `json:"sender_blockchain_address"`
		Recipient string  `json:"recipient_blockchain_address"`
		Value     float32 `json:"value"`
	}{
		Sender:    t.SenderBlockchainAddress,
		Recipient: t.RecipientBlockchainAddress,
		Value:     t.Value,
	})
}
