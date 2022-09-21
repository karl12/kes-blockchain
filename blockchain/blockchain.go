package blockchainservice/blockchain

import "fmt"

type Block struct {
	hash         string
	previousHash string
	data         string
	timestamp    int64
	nonce        uint64
}

type Blockchain struct {
	blocks     []Block
	difficulty int
}

func init() {
	fmt.Println("Block package initialized")
}
