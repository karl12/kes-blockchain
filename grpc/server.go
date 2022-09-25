package grpc

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	pb "kes-blockchain/gen/node/v1"
	"kes-blockchain/node"
	nw "kes-blockchain/util"
	"log"
	"net"
)

func Start(address string) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterNodeServer(s, &server{})

	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type server struct {
	pb.UnimplementedNodeServer
}

func (s *server) JoinNetwork(ctx context.Context, in *pb.DiscoverRequest) (*pb.PeersReply, error) {
	addresses := nw.CreateNetworkAddressSlice(node.GetPeers())
	node.GetPeers()[in.Address] = true
	return &pb.PeersReply{Address: addresses}, nil
}

func (s *server) NodeList(ctx context.Context, in *pb.Empty) (*pb.PeersReply, error) {
	addresses := nw.CreateNetworkAddressSlice(node.GetPeers())
	return &pb.PeersReply{Address: addresses}, nil
}

func (s *server) MineBlock(ctx context.Context, in *pb.Empty) (*pb.BlockReply, error) {
	blockchain := node.MineNext()
	toMessage := nw.BlockToMessage(blockchain.Blocks)
	return &pb.BlockReply{Block: toMessage[len(toMessage)-1]}, nil
}

func (s *server) Blockchain(ctx context.Context, in *pb.Empty) (*pb.BlockchainReply, error) {
	return &pb.BlockchainReply{
		Blockchain: &pb.Blockchain{
			Blocks: nw.BlockToMessage(node.GetBlockchain().Blocks),
		},
	}, nil
}
