package main

import (
	blockchain "blockchainservice/blockchain"
	"crypto/sha256"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	bc = blockchain.Blockchain{[]blockchain.Block{}, 4}
)

const MaxUint = ^uint64(0)

func createGenesis() {
	inceptionTime := int64(0)
	data := ""
	previousHash := ""
	nonce := uint64(0)
	hash := hash(previousHash, data, inceptionTime, nonce)

	block := Block{hash, previousHash, data, inceptionTime, nonce}

	bc.blocks = append(bc.blocks, block)
}

func mine(previousHash string, data string) Block {
	for {
		for nonce := uint64(0); nonce <= MaxUint; nonce++ {
			timestamp := time.Now().Unix()
			potentialHash := hash(previousHash, data, timestamp, nonce)
			if validateHash(potentialHash) {
				return Block{potentialHash, previousHash, data, timestamp, nonce}
			}
		}
	}
}

func validateHash(potentialHash string) bool {
	prefix := ""
	for i := 0; i < bc.difficulty; i++ {
		prefix += "0"
	}

	return strings.HasPrefix(potentialHash, prefix)
}

func hash(previousHash string, data string, timestamp int64, nonce uint64) string {
	toHash := previousHash + data + strconv.FormatInt(timestamp, 10) + strconv.FormatUint(nonce, 10)
	hash := sha256.Sum256([]byte(toHash))
	return fmt.Sprintf("%x", hash)
}

func blockChainLength() int {
	return len(bc.blocks)
}

func main() {
	createGenesis()
	print(bc.blocks[0].hash)

	for i := 0; i < 10; i++ {
		bc.blocks = append(bc.blocks, mine(bc.blocks[blockChainLength()-1].previousHash, "test"))
		print()
	}
}
