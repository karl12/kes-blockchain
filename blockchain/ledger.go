package blockchainservice

type Transaction struct {
	Sender   string
	Receiver string
	Amount   uint64
}

func CalculateAccountTotals(blocks []Block) map[string]uint64 {
	ret := make(map[string]uint64)
	for i := 0; i < len(blocks); i++ {
		block := blocks[i]
		for _, transaction := range block.Data {

			if transaction.Sender != "" {
				ret[transaction.Sender] -= transaction.Amount
			}

			ret[transaction.Receiver] += transaction.Amount
		}
	}

	return ret
}
