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
	addresses := createAddressSlice()
	peers[in.Address] = true
	return &pb.PeersReply{Address: addresses}, nil
}

func (s *server) NodeList(ctx context.Context, in *pb.Empty) (*pb.PeersReply, error) {
	addresses := createAddressSlice()
	return &pb.PeersReply{Address: addresses}, nil
}

func (s *server) MineBlock(ctx context.Context, in *pb.Empty) (*pb.BlockReply, error) {
	blockchain = blockchain.MineNext(nil)
	tooMany := blockToMessage(blockchain.Blocks)
	return &pb.BlockReply{Block: tooMany[len(tooMany)-1]}, nil
}

func (s *server) Blockchain(ctx context.Context, in *pb.Empty) (*pb.BlockchainReply, error) {
	return &pb.BlockchainReply{
		Blockchain: &pb.Blockchain{
			Blocks: blockToMessage(blockchain.Blocks),
		},
	}, nil
}

func blockToMessage(blocks []bc.Block) []*pb.Block {
	if &blocks == nil {
		return nil
	}

	var ret []*pb.Block

	for _, block := range blocks {
		transactions := transactionsToMessage(&block)

		message := pb.Block{
			Hash:         block.Hash,
			PreviousHash: block.PreviousHash,
			Timestamp:    block.Timestamp,
			Nonce:        block.Nonce,
			Data:         transactions,
		}

		ret = append(ret, &message)
	}

	return ret
}

func transactionsToMessage(block *bc.Block) []*pb.Transaction {
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

func messageToBlock(blocks []*pb.Block) []bc.Block {
	if &blocks == nil {
		return nil
	}

	var ret []bc.Block

	for _, block := range blocks {

		transactions := messageToTransaction(block)

		message := bc.Block{
			Hash:         block.Hash,
			PreviousHash: block.PreviousHash,
			Timestamp:    block.Timestamp,
			Nonce:        block.Nonce,
			Data:         transactions,
		}

		ret = append(ret, message)
	}

	return ret
}

func messageToTransaction(block *pb.Block) []bc.Transaction {
	var transactions []bc.Transaction

	for _, t := range block.Data {
		transactions = append(transactions, bc.Transaction{
			Sender:   t.Sender,
			Receiver: t.Receiver,
			Amount:   t.Amount,
		})
	}
	return transactions
}

func createAddressSlice() []string {
	var addresses []string
	for address, _ := range peers {
		addresses = append(addresses, address)
	}
	return addresses
}

func sync() {
	addresses := createAddressSlice()
	for _, address := range addresses {
		if address != ownAddress {
			for _, newAddress := range getPeersFromTarget(address) {
				peers[newAddress] = true
			}
		}
	}

	log.Printf("Peers in network: %d", len(peers))

	for _, address := range addresses {
		foreignBlock := getBlockChain(address)
		if &foreignBlock != nil && blockchain.ShouldReplace(foreignBlock) {
			blockchain = foreignBlock
			log.Printf("Block replaced from %s", address)
		}
	}

}

func getBlockChain(target string) bc.Blockchain {
	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		delete(peers, target)
		return bc.Blockchain{}
	}
	defer conn.Close()
	client := pb.NewNodeClient(conn)

	ret, err := client.Blockchain(
		context.Background(),
		&pb.Empty{})

	if err != nil {
		log.Printf("Bad peer: error %s", err)
		delete(peers, target)
		return bc.Blockchain{}
	}

	return messageToBlockchain(ret.Blockchain)
}

func messageToBlockchain(b *pb.Blockchain) bc.Blockchain {
	return bc.Blockchain{
		Blocks: messageToBlock(b.Blocks),
	}
}

func getPeersFromTarget(target string) []string {
	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		delete(peers, target)
		return []string{}
	}
	defer conn.Close()
	client := pb.NewNodeClient(conn)

	ret, err := client.JoinNetwork(
		context.Background(),
		&pb.DiscoverRequest{Address: ownAddress})

	if err != nil {
		log.Printf("Bad peer: error %s", err)
		delete(peers, target)
		return []string{}
	}

	return ret.Address
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
