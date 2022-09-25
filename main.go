package main

import (
	"kes-blockchain/grpc"
	"kes-blockchain/node"
	"os"
)

func main() {
	address := os.Args[1]

	go node.Start(address)
	grpc.Start(address)

}
