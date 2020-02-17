module github.com/rhizome-chain/ethereum

go 1.13

require (
	github.com/ethereum/go-ethereum v1.9.9
	github.com/rhizome-chain/tendermint-daemon v0.0.1
	github.com/spf13/cobra v0.0.5
)

replace github.com/rhizome-chain/tendermint-daemon => ../tendermint-daemon
