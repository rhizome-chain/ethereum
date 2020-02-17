package main

import (
	"github.com/rhizome-chain/ethereum/module"
	cmd "github.com/rhizome-chain/tendermint-daemon/cmd/commands"
	"github.com/rhizome-chain/tendermint-daemon/daemon"
)

func main() {
	daemonProvider := &daemon.BaseProvider{}
	
	daemonProvider.AddModule(&module.EthModule{})
	
	cmd.DoCmd(daemonProvider)
}
