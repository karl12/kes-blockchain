package main

import (
	bc "kes-blockchain/blockchain"
	pb "kes-blockchain/gen/node/v1"
)

func BlockToMessage(blocks []bc.Block) []*pb.Block {
	if &blocks == nil {
		return nil
	}

	var ret []*pb.Block

	for _, block := range blocks {
		transactions := TransactionsToMessage(&block)

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

func TransactionsToMessage(block *bc.Block) []*pb.Transaction {
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

func MessageToBlock(blocks []*pb.Block) []bc.Block {
	if &blocks == nil {
		return nil
	}

	var ret []bc.Block

	for _, block := range blocks {

		transactions := MessageToTransaction(block)

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

func MessageToTransaction(block *pb.Block) []bc.Transaction {
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

func messageToBlockchain(b *pb.Blockchain) bc.Blockchain {
	return bc.Blockchain{
		Blocks: MessageToBlock(b.Blocks),
	}
}
