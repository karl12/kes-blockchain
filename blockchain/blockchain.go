package blockchainservice

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Block struct {
	Hash         string
	PreviousHash string
	Data         []Transaction
	Timestamp    int64
	Nonce        uint64
	Next         *Block
	Previous     *Block
}

func dataString(data []Transaction) string {
	ret := ""
	for _, s := range data {
		ret += s.Sender
		ret += "->"
		ret += s.Receiver
		ret += ":"
		ret += strconv.FormatUint(s.Amount, 10)
		ret += "|"
	}
	return ret
}

type blockchain struct {
	firstBlock Block
	lastBlock  Block
	length     uint64
	difficulty int
}

const MaxUint = ^uint64(0)

var (
	bc = blockchain{createGenesis(), createGenesis(), 1, 4}
)

func createGenesis() Block {
	inceptionTime := int64(0)
	var data []Transaction
	data = append(data, Transaction{Sender: "", Receiver: "karl", Amount: 100000000})
	previousHash := ""
	nonce := uint64(0)
	hash := GenerateHash(previousHash, "", inceptionTime, nonce)

	block := Block{hash, previousHash, data, inceptionTime, nonce, nil, nil}
	return block
}

func GenerateHash(previousHash string, data string, timestamp int64, nonce uint64) string {
	toHash := previousHash + data + strconv.FormatInt(timestamp, 10) + strconv.FormatUint(nonce, 10)
	hash := sha256.Sum256([]byte(toHash))
	return fmt.Sprintf("%x", hash)
}

func Mine(data []Transaction) {
	for {
		startTimestamp := time.Now().Unix()
		for nonce := uint64(0); nonce <= MaxUint; nonce++ {
			if startTimestamp != time.Now().Unix() {
				break
			}

			potentialHash := GenerateHash(GetLastBlock().Hash, dataString(data), startTimestamp, nonce)
			if validateHash(potentialHash) {
				nextBlock := Block{potentialHash, GetLastBlock().Hash, data, startTimestamp, nonce, nil, &bc.lastBlock}
				bc.lastBlock.Next = &nextBlock
				bc.lastBlock = nextBlock
				bc.length++
				return
			}
		}

		if startTimestamp == time.Now().Unix() {
			fmt.Println("Too quick")
		}
	}
}

func GetBlockChainLength() uint64 {
	return bc.length
}

func GetLastBlock() Block {
	return bc.lastBlock
}
func GetFirstBlock() Block {
	return bc.firstBlock
}

func validateHash(potentialHash string) bool {
	prefix := ""
	for i := 0; i < bc.difficulty; i++ {
		prefix += "0"
	}

	return strings.HasPrefix(potentialHash, prefix)
}

func ValidateChain(firstBlock Block) bool {
	previousHash := createGenesis().Hash
	block := firstBlock.Next
	for ; block != nil; block = block.Next {
		expectedHash := GenerateHash(block.PreviousHash, dataString(block.Data), block.Timestamp, block.Nonce)
		if previousHash != block.PreviousHash || expectedHash != block.Hash || validateHash(expectedHash) {
			return false
		}
	}
	return true
}
