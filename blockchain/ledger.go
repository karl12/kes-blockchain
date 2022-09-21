package blockchainservice

type Transaction struct {
	Sender   string
	Receiver string
	Amount   uint64
}

func CalculateAccountTotals(block Block) map[string]uint64 {
	ret := make(map[string]uint64)
	blockPointer := &block
	for ; blockPointer != nil; blockPointer = blockPointer.Next {
		for _, transaction := range block.Data {

			if transaction.Sender != "" {
				ret[transaction.Sender] -= transaction.Amount
			}

			ret[transaction.Receiver] += transaction.Amount
		}
	}

	return ret
}
