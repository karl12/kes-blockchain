package main

import (
	"google.golang.org/grpc"
	"log"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc/reflection"
	blockchain "kes-blockchain/blockchain"
	pb "kes-blockchain/pb/node/v1"
)

type server struct {
	pb.UnimplementedNodeServer
}

func (s *server) DiscoverPeers(ctx context.Context, in *pb.Empty) (*pb.PeersReply, error) {
	response := []string{"Bob", "John", "Peter"}
	return &pb.PeersReply{Message: response}, nil
}

func (s *server) LastBlock(ctx context.Context, in *pb.Empty) (*pb.BlockReply, error) {
	return &pb.BlockReply{Block: blockToMessage(blockchain.GetLastBlock())}, nil
}

func (s *server) MineBlock(ctx context.Context, in *pb.Empty) (*pb.BlockReply, error) {
	blockchain.Mine(nil)
	return &pb.BlockReply{Block: blockToMessage(blockchain.GetLastBlock())}, nil
}

func blockToMessage(block *blockchain.Block) *pb.Block {
	if &block == nil {
		return nil
	}

	transactions := transactionsToMessage(block)

	message := &pb.Block{
		Hash:         block.Hash,
		PreviousHash: block.PreviousHash,
		Timestamp:    block.Timestamp,
		Nonce:        block.Nonce,
		Data:         transactions,
	}

	if block.Previous != nil {
		message.Previous = blockToMessage(block.Previous)
	}

	return message
}

func transactionsToMessage(block *blockchain.Block) []*pb.Transaction {
	var transactions []*pb.Transaction

	for _, t := range block.Data {
		transactions = append(transactions, &pb.Transaction{
			Sender:   t.Sender,
			Receiver: t.Receiver,
			Amount:   t.Amount,
		})
	}
	return transactions
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
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
