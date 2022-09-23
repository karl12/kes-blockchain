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
}

type Blockchain struct {
	Blocks []Block
}

const MaxUint64 = ^uint64(0)
const DIFFICULTY = 4

func New() Blockchain {
	return Blockchain{
		Blocks: []Block{createGenesis()},
	}
}

func (bc Blockchain) MineNext(data []Transaction) Blockchain {
	for {
		startTimestamp := time.Now().Unix()
		for nonce := uint64(0); nonce <= MaxUint64; nonce++ {
			if startTimestamp != time.Now().Unix() {
				break
			}

			potentialHash := generateHash(bc.Blocks[len(bc.Blocks)-1].Hash, dataString(data), startTimestamp, nonce)
			if validateHash(potentialHash) {
				bc.Blocks = append(bc.Blocks, Block{potentialHash, bc.Blocks[len(bc.Blocks)-1].Hash, data, startTimestamp, nonce})
				return bc
			}
		}

		if startTimestamp == time.Now().Unix() {
			fmt.Println("Too quick")
		}
	}
}

func (bc Blockchain) ShouldReplace(comp Blockchain) bool {
	if len(comp.Blocks) > len(bc.Blocks) && comp.validateChain() {
		return true
	}
	return false
}

func (bc Blockchain) validateChain() bool {
	for i := 1; i < len(bc.Blocks); i++ {
		expectedHash := generateHash(bc.Blocks[i].PreviousHash, dataString(bc.Blocks[i].Data), bc.Blocks[i].Timestamp, bc.Blocks[i].Nonce)
		if bc.Blocks[i].Hash != expectedHash || bc.Blocks[i].PreviousHash != bc.Blocks[i-1].Hash || !validateHash(expectedHash) {
			return false
		}
	}
	return true
}

func validateHash(potentialHash string) bool {
	prefix := ""
	for i := 0; i < DIFFICULTY; i++ {
		prefix += "0"
	}

	return strings.HasPrefix(potentialHash, prefix)
}

func generateHash(previousHash string, data string, timestamp int64, nonce uint64) string {
	toHash := previousHash + data + strconv.FormatInt(timestamp, 10) + strconv.FormatUint(nonce, 10)
	hash := sha256.Sum256([]byte(toHash))
	return fmt.Sprintf("%x", hash)
}

func createGenesis() Block {
	inceptionTime := int64(0)
	var data []Transaction
	data = append(data, Transaction{Sender: "", Receiver: "karl", Amount: 100000000})
	previousHash := ""
	nonce := uint64(0)
	hash := generateHash(previousHash, "", inceptionTime, nonce)

	block := Block{hash, previousHash, data, inceptionTime, nonce}
	return block
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
