package node

import (
	"github.com/jasonlvhit/gocron"
	bc "kes-blockchain/blockchain"
	"log"
)

type Client interface {
	Sync()
}

var (
	address    string
	peers      = map[string]bool{"localhost:50051": true}
	blockchain = bc.New()
)

func Start(assignedAddress string) {
	var client Client
	client = GRPCNode{}

	address = assignedAddress
	scheduler := gocron.NewScheduler()
	scheduler.Every(5).Seconds().Do(client.Sync)
	scheduler.Start()
}

func GetPeers() map[string]bool {
	return peers
}

func GetBlockchain() bc.Blockchain {
	return blockchain
}

func MineNext() bc.Blockchain {
	blockchain.MineNext(nil)
	return blockchain
}

func checkAndReplaceBlockchainIfLonger(foreignBlock bc.Blockchain) {
	if &foreignBlock != nil && blockchain.ShouldReplace(foreignBlock) {
		blockchain = foreignBlock
		log.Printf("Block replaced")
	}
}

func mergePeers(newPeers []string) {
	for _, newAddress := range newPeers {
		peers[newAddress] = true
	}
}
