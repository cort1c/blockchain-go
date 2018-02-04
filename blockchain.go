package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type (
	Blockchain struct {
		Chain               []Block
		CurrentTransactions []Transaction
	}

	Block struct {
		Index        int
		Timestamp    int64
		Transactions []Transaction
		Proof        int
		PreviousHash string
	}

	Transaction struct {
		Sender    string
		Recipient string
		Amount    int
	}
)

func New() *Blockchain {
	b := &Blockchain{}
	b.NewBlock(100, "1")
	return b
}

func (b *Blockchain) NewBlock(proof int, previousHash string) Block {
	block := Block{
		Index:        len(b.Chain) + 1,
		Timestamp:    time.Now().Unix(),
		Transactions: b.CurrentTransactions,
		Proof:        proof,
		PreviousHash: previousHash,
	}
	log.Printf("block.Index %v\n", block.Index)
	b.CurrentTransactions = []Transaction{}
	b.Chain = append(b.Chain, block)
	return block
}

func (b *Blockchain) NewTransaction(sender, recipient string, amount int) int {
	transaction := Transaction{
		Sender:    sender,
		Recipient: recipient,
		Amount:    amount,
	}
	b.CurrentTransactions = append(b.CurrentTransactions, transaction)
	return b.LastBlock().Index + 1
}

func (b *Blockchain) LastBlock() Block {
	return b.Chain[len(b.Chain)-1]
}

func (b *Blockchain) ProofOfWork(lastProof int) int {
	proof := 0
	for ; !ValidProof(lastProof, proof); proof++ {
	}
	return proof
}

func ValidProof(lastProof, proof int) bool {
	guess := string(lastProof) + string(proof)
	guessHash := encodeToHashString([]byte(guess))
	fmt.Printf("guessHash: %s, lastProof: %v, proof: %v\n", guessHash, lastProof, proof)
	return guessHash[:4] == "0000"
}

func Hash(block Block) (string, error) {
	b, err := json.Marshal(block)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return encodeToHashString(b), nil
}

func encodeToHashString(bytes []byte) string {
	h := sha1.New()
	h.Write(bytes)
	return hex.EncodeToString(h.Sum(nil))
}
