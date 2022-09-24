package main

import (
	"github.com/jasonlvhit/gocron"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	bc "kes-blockchain/blockchain"
	pb "kes-blockchain/gen/node/v1"
	"log"
	"net"
	"os"
)

var (
	blockchain bc.Blockchain
	peers      = map[string]bool{"localhost:50051": true}
	ownAddress string
)

type server struct {
	pb.UnimplementedNodeServer
}

func (s *server) JoinNetwork(ctx context.Context, in *pb.DiscoverRequest) (*pb.PeersReply, error) {
	addresses := createNetworkAddressSlice()
	peers[in.Address] = true
	return &pb.PeersReply{Address: addresses}, nil
}

func (s *server) NodeList(ctx context.Context, in *pb.Empty) (*pb.PeersReply, error) {
	addresses := createNetworkAddressSlice()
	return &pb.PeersReply{Address: addresses}, nil
}

func (s *server) MineBlock(ctx context.Context, in *pb.Empty) (*pb.BlockReply, error) {
	blockchain = blockchain.MineNext(nil)
	tooMany := BlockToMessage(blockchain.Blocks)
	return &pb.BlockReply{Block: tooMany[len(tooMany)-1]}, nil
}

func (s *server) Blockchain(ctx context.Context, in *pb.Empty) (*pb.BlockchainReply, error) {
	return &pb.BlockchainReply{
		Blockchain: &pb.Blockchain{
			Blocks: BlockToMessage(blockchain.Blocks),
		},
	}, nil
}

func createNetworkAddressSlice() []string {
	var addresses []string
	for address, _ := range peers {
		addresses = append(addresses, address)
	}
	return addresses
}

func sync() {
	addresses := createNetworkAddressSlice()
	syncFunctions := []func(client pb.NodeClient) error{
		getPeersFromTarget,
		syncBlockchain,
	}

	for _, address := range addresses {
		if address != ownAddress {
			execute[[]string](address, syncFunctions)
		}
	}
	log.Printf("Peers in network: %d", len(peers))
}

func syncBlockchain(client pb.NodeClient) error {
	ret, err := client.Blockchain(
		context.Background(),
		&pb.Empty{})

	if &err != nil {
		return err
	}

	foreignBlock := messageToBlockchain(ret.Blockchain)

	if &foreignBlock != nil && blockchain.ShouldReplace(foreignBlock) {
		blockchain = foreignBlock
		log.Printf("Block replaced")
	}

	return nil
}

func getPeersFromTarget(client pb.NodeClient) error {
	ret, err := client.JoinNetwork(
		context.Background(),
		&pb.DiscoverRequest{Address: ownAddress})

	for _, newAddress := range ret.Address {
		peers[newAddress] = true
	}
	return err
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

func main() {
	blockchain = bc.New()
	ownAddress = os.Args[1]
	lis, err := net.Listen("tcp", ownAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterNodeServer(s, &server{})

	scheduler := gocron.NewScheduler()
	scheduler.Every(5).Seconds().Do(sync)
	scheduler.Start()

	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
