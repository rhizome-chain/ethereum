module github.com/rhizome-chain/ethereum

go 1.13

require (
	github.com/ethereum/go-ethereum v1.9.9
	github.com/rhizome-chain/tendermint-daemon v0.0.2
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.6.2
	github.com/tendermint/tendermint v0.33.0
)

replace github.com/rhizome-chain/tendermint-daemon => ../tendermint-daemon
