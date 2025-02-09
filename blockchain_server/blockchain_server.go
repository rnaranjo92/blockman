package main

import (
	"blockman/block"
	"blockman/wallet"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
)

var cache map[string]*block.BlockChain = make(map[string]*block.BlockChain)

type BlockChainServer struct {
	port uint16
}

func NewBlockChainServer(port uint16) *BlockChainServer {
	return &BlockChainServer{port}
}

func (bcs *BlockChainServer) Port() uint16 {
	return bcs.port
}

func (bcs *BlockChainServer) GetBlockchain() *block.BlockChain {
	bc, ok := cache["blockchain"]

	if !ok {
		minersWalet := wallet.NewWallet()
		bc = block.NewBlockChain(minersWalet.BlockChainAddress, bcs.Port())
		cache["blockchain"] = bc
		log.Printf("private_key %v", minersWalet.GetPrivateKeyStr())
		log.Printf("public_key %v", minersWalet.GetPublicKeyStr())
		log.Printf("blockchain_address %v", minersWalet.GetBlockChainAddress())
	}

	return bc
}

func (bcs *BlockChainServer) GetChain(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")
		bc := bcs.GetBlockchain()
		m, _ := bc.MarshalJSON()
		io.WriteString(w, string(m[:]))
	default:
		log.Printf("ERROR: Invalid HTTP Method")
	}
}

func (bcs *BlockChainServer) Transactions(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		// Work on this later
	case http.MethodPost:
		decoder := json.NewDecoder(req.Body)
		var t block.TransactionRequest
		err := decoder.Decode(&t)
		if err != nil {
			log.Println("ERROR: %v", err)
			io.WriteString()
		}
	default:

	}
}

func (bcs *BlockChainServer) Run() {
	http.HandleFunc("/", bcs.GetChain)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(bcs.Port())), nil))
}
