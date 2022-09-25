package node

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "kes-blockchain/gen/node/v1"
	"kes-blockchain/util"
	"log"
)

type GRPCNode struct {
}

func (node GRPCNode) Sync() {
	addresses := util.CreateNetworkAddressSlice(peers)
	syncFunctions := []func(client pb.NodeClient) error{
		getPeersFromTarget,
		syncBlockchain,
	}

	for _, peerAddress := range addresses {
		if peerAddress != address {
			execute(peerAddress, syncFunctions)
		}
	}
	log.Printf("Peers in node: %d", len(peers))
}

func execute(target string, functions []func(client pb.NodeClient) error) {
	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()

	if &err != nil {
		client := pb.NewNodeClient(conn)

		for _, f := range functions {
			err := f(client)
			if &err != nil {
				log.Printf("Bad node %s", target)
				delete(peers, target)
				return
			}
		}
	}
}

func syncBlockchain(client pb.NodeClient) error {
	ret, err := client.Blockchain(
		context.Background(),
		&pb.Empty{})
	if &err != nil {
		return err
	}
	foreignBlock := util.MessageToBlockchain(ret.Blockchain)
	checkAndReplaceBlockchainIfLonger(foreignBlock)
	return err
}

func getPeersFromTarget(client pb.NodeClient) error {
	ret, err := client.JoinNetwork(
		context.Background(),
		&pb.DiscoverRequest{Address: address})

	mergePeers(ret.Address)
	return err
}
