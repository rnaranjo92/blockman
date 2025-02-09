package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	// "golang.org/x/crypto/ripemd160" //deprecated
)

type Wallet struct {
	PrivateKey        *ecdsa.PrivateKey
	PublicKey         *ecdsa.PublicKey
	BlockChainAddress string
}

func NewWallet() *Wallet {
	w := &Wallet{}
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	w.PrivateKey = privateKey
	w.PublicKey = &w.PrivateKey.PublicKey

	h2 := sha256.New()
	h2.Write(w.PublicKey.X.Bytes())
	h2.Write(w.PublicKey.Y.Bytes())
	digest2 := h2.Sum(nil)

	// h3 := ripemd160.New()
	// h3.Write(digest2)
	// digest3 := h3.Sum(nil)

	vd4 := make([]byte, 21)
	vd4[0] = 0x00
	copy(vd4[1:], digest2[:])

	h5 := sha256.New()
	h5.Write(vd4)
	digest5 := h5.Sum(nil)

	h6 := sha256.New()
	h6.Write(digest5)
	digest6 := h6.Sum(nil)

	chsum := digest6[:4]

	dc8 := make([]byte, 25)
	copy(dc8[:21], vd4[:])
	copy(dc8[21:], chsum[:])

	address := base58.Encode(dc8)
	w.BlockChainAddress = address
	return w
}

func (w *Wallet) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		PrivateKey        string `json:"private_key"`
		PublicKey         string `json:"public_key"`
		BlockChainAddress string `json:"blockchain_address"`
	}{
		PrivateKey:        w.GetPrivateKeyStr(),
		PublicKey:         w.GetPublicKeyStr(),
		BlockChainAddress: w.GetBlockChainAddress(),
	})
}

func (w *Wallet) GetPrivateKey() *ecdsa.PrivateKey {
	return w.PrivateKey
}

func (w *Wallet) GetPrivateKeyStr() string {
	return fmt.Sprintf("%x", w.PrivateKey.D.Bytes())
}

func (w *Wallet) GetPublicKey() *ecdsa.PublicKey {
	return w.PublicKey
}

func (w *Wallet) GetPublicKeyStr() string {
	return fmt.Sprintf("%064x%064x", w.PublicKey.X.Bytes(), w.PublicKey.Y.Bytes())
}

func (w *Wallet) GetBlockChainAddress() string {
	return w.BlockChainAddress
}

type TransactionRequest struct {
	SenderPrivateKey           *string `json:"sender_private_key"`
	SenderBlockchainAddress    *string `json:"sender_blockchain_address"`
	SenderPublicKey            *string `json:"sender_public_key"`
	RecipientBlockchainAddress *string `json:"recipient_blockchain_address"`
	Value                      *string `json:"value"`
}

func (tr *TransactionRequest) Validate() bool {
	if tr.SenderPrivateKey == nil ||
		tr.SenderPublicKey == nil ||
		tr.RecipientBlockchainAddress == nil ||
		tr.SenderBlockchainAddress == nil ||
		tr.Value == nil {
		return false
	}
	return true
}
