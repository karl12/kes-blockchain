package util

func CreateNetworkAddressSlice(peerMap map[string]bool) []string {
	var addresses []string
	for address, _ := range peerMap {
		addresses = append(addresses, address)
	}
	return addresses
}
