package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
)

type Block struct {
	Data     string
	Nounce   string
	Hash     string
	PrevHash string
}

const (
	genesisHash = "000000000000000000000000000000000000"
	nounceRule  = "00000"
)

// A Store that maps Blocks against its hash
var blockStore = make(map[string]*Block)

// The block Chain last Block Hash
var blockChain string

// A function to Generate hash
func generateHash(newB *Block) (string, string) {
	copyB := make([]string, 3)
	nounce := 0
	hash := ""
	// generate a nounce and a hash that matches the nounce rule
	for true {
		// we use array to maintain the order
		// [0] -> Data
		// [1] -> Nounce
		// [2] -> PrevHash
		copyB[0] = newB.Data
		copyB[1] = fmt.Sprintf("%d", nounce)
		copyB[2] = newB.PrevHash
		dump, _ := json.Marshal(copyB)
		h := sha256.New()
		h.Write(dump)
		hash = fmt.Sprintf("%x", h.Sum(nil))
		if strings.HasPrefix(hash, nounceRule) {
			break
		}
		nounce++
	}
	return hash, fmt.Sprintf("%d", nounce)
}

// A function to Verify Hash
func verifyHash(block *Block) (bool, string) {
	// we use array to maintain the order
	// [0] -> Data
	// [1] -> Nounce
	// [2] -> PrevHash
	copyB := make([]string, 3)
	copyB[0] = block.Data
	copyB[1] = block.Nounce
	copyB[2] = block.PrevHash
	dump, _ := json.Marshal(copyB)
	h := sha256.New()
	h.Write(dump)
	hash := fmt.Sprintf("%x", h.Sum(nil))
	if block.Hash == hash {
		return true, hash
	}
	return false, hash
}

// A function that adds a block to the blockchain
func addToBlockChain(newB *Block) {
	newB.Hash, newB.Nounce = generateHash(newB)
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
		// Check if previousHash of the next block matchs
		if previousHash != "" && previousHash != block.Hash {
			fmt.Printf("Invalid block hash in next block, expected hash: %s, but found: %s\n", previousHash, block.Hash)
			return false
		}

		// Check if the hash match nounce value
		if !strings.HasPrefix(block.Hash, nounceRule) {
			fmt.Printf("Invalid hash as nounce doesn't match, block hash: %s, but nounce rule: %s\n", block.Hash, nounceRule)
			return false
		}

		// Check if the block hash is valid
		if verified, hash := verifyHash(block); !verified {
			fmt.Printf("Invalid hash for block, block hash: %s, but generated hash: %s\n", block.Hash, hash)
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
