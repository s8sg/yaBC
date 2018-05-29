package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

type Block struct {
	Data     string
	Hash     string
	PrevHash string
}

const (
	genesisHash = "000000000000000000000000000000000000"
)

// A Store that maps Blocks against its hash
var blockStore = make(map[string]*Block)

// The block Chain last Block Hash
var blockChain string

// A function to Generate hash
func generateHash(newB *Block) string {
	// we use array to maintain the order
	copyB := make([]string, 2)
	copyB[0] = newB.Data
	copyB[1] = newB.PrevHash
	dump, _ := json.Marshal(copyB)
	h := sha256.New()
	h.Write(dump)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// A function that adds a block to the blockchain
func addToBlockChain(newB *Block) {
	newB.Hash = generateHash(newB)
	// add to blockStore
	blockStore[newB.Hash] = newB
	// change the last block hash
	blockChain = newB.Hash
}

// A function to verify a blockchain
func verifyBlockChain(debug bool) bool {
	block := blockStore[blockChain]
	previousHash := ""
	for true {
		if debug {
			fmt.Printf("Block %s { \n Data: %s,\n PrevHash: %s \n}\n", block.Hash, block.Data, block.PrevHash)
		}
		hash := generateHash(block)
		// Check if current hash is right
		if hash != block.Hash {
			fmt.Printf("Invalid block hash for block, expected hash %s, but found: %s\n", hash, block.Hash)
			return false
		}
		// Check if previousHash of the next block matchs
		if previousHash != "" && previousHash != hash {
			fmt.Printf("Invalid block hash in next block, expected hash %s, but found: %s\n", previousHash, hash)
			return false
		}

		// Get the previous hash and break if genesis
		previousHash = block.PrevHash
		if previousHash == genesisHash {
			break
		}

		// Get the previous blog
		block = blockStore[previousHash]
	}
	return true
}

// A function that create a new block and adds it to a verified blockchain
func verifyAndAdd(data string) {
	// verify if the current blockchain is valid
	if !verifyBlockChain(false) {
		panic("Block chain is not sanitized")
	}
	// create a new block
	newB := &Block{Data: data, PrevHash: blockChain}
	// add it to block chain
	addToBlockChain(newB)
}

func initBlockChain() {
	// Add genesis block
	newB := &Block{Data: "genesis", PrevHash: genesisHash}
	addToBlockChain(newB)
}

func main() {
	initBlockChain()
	verifyAndAdd("data1")
	verifyAndAdd("data2")
	verifyAndAdd("data3")
	verifyAndAdd("data4")
	verifyBlockChain(true)
}
