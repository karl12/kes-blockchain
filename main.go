package main

import (
	blockchain "blockchainservice/blockchain"
	"fmt"
)

func main() {
	fmt.Println(blockchain.GetBlockChainLength())
	fmt.Println(blockchain.GetLastBlock().Hash)
	for i := 0; i < 10; i++ {
		blockchain.Mine(nil)
		fmt.Println(blockchain.GetBlockChainLength())
		fmt.Println(blockchain.GetLastBlock().Hash)
	}
	fmt.Println(blockchain.ValidateChain(blockchain.GetFirstBlock()))
	blockchain.CalculateAccountTotals(blockchain.GetFirstBlock())

}
