package module

import (
	"fmt"
	"path/filepath"
	
	erc20 "github.com/rhizome-chain/ethereum/subs/erc20"
	erc721 "github.com/rhizome-chain/ethereum/subs/erc721"
	
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/config"
	
	// "errors"
	"github.com/rhizome-chain/ethereum/subs"
	"github.com/rhizome-chain/tendermint-daemon/daemon"
	"github.com/rhizome-chain/tendermint-daemon/daemon/worker"
	"github.com/rhizome-chain/tendermint-daemon/types"
)

const Name = "eth"

type EthModule struct {
	modcfg  *subs.EthConfig
	manager *subs.EthSubsManager
}

var _ daemon.Module = (*EthModule)(nil)

func (e *EthModule) GetFactory(name string) worker.Factory {
	if name == subs.FactoryName {
		return e.manager
	}
	return nil
}

func (e *EthModule) Name() string {
	return Name
}

func (e *EthModule) GetConfig() types.ModuleConfig {
	return e.modcfg
}

func (e *EthModule) Factories() (facs []worker.Factory) {
	return []worker.Factory{e.manager}
}

func (e *EthModule) Init(config *config.Config) {
	e.modcfg = loadFile(config)
	ethUrl := e.modcfg.NetworkURL
	
	if len(ethUrl) > 0 {
		e.modcfg.NetworkURL = ethUrl
	} else {
		ethUrl = e.modcfg.NetworkURL
	}
	
	if len(ethUrl) == 0 {
		panic("Ethereum network url is not set.")
	}
	
	manager := subs.NewEthSubsManager(ethUrl)
	
	manager.RegisterLogHandler(erc20.NewERC20LogHandler())
	manager.RegisterLogHandler(erc721.NewERC721LogHandler())
	
	e.manager = manager
}

func loadFile(config *config.Config) *subs.EthConfig {
	confFilePath := filepath.Join(config.RootDir, "config", "ethereum.toml")
	ethConfig := &subs.EthConfig{}
	types.LoadModuleConfigFile(confFilePath, ethConfig)
	fmt.Println("[EthModule] Load EthConfig file:", confFilePath, ethConfig)
	return ethConfig
}

func (e *EthModule) BeforeDaemonStarting(cmd *cobra.Command, dm *daemon.Daemon) {
	if e.modcfg == nil {
		panic("EthConfig is net set.")
	}
}

func (e *EthModule) AfterDaemonStarted(dm *daemon.Daemon) {
	dm.GetContext().Info("EthModule starts")
}
