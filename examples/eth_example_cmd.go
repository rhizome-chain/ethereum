package main

import (
	cmd "github.com/rhizome-chain/tendermint-daemon/cmd/commands"
	"github.com/rhizome-chain/tendermint-daemon/daemon"
	"github.com/rhizome-chain/tendermint-daemon/daemon/worker"
	
	"github.com/rhizome-chain/ethereum/subs"
	erc20 "github.com/rhizome-chain/ethereum/subs/erc20"
	erc721 "github.com/rhizome-chain/ethereum/subs/erc721"
)

func main() {
	// project secret : e986d9e29d5243fd8f6e478198ba464b
	manager := subs.NewEthSubsManager("wss://mainnet.infura.io/ws/v3/a7f6d7ea8be04689a9b0394b7378451b")
	// Register default built-in handlers
	manager.RegisterLogHandler(&erc20.ERC20LogHandler{})
	manager.RegisterLogHandler(&erc721.ERC721LogHandler{})
	
	daemonProvider := &daemon.BaseProvider{}
	daemonProvider.Factories = []worker.Factory{manager}
	cmd.DoCmd(daemonProvider)
}


