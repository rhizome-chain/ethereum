package module

import (
	"fmt"
	erc20 "github.com/rhizome-chain/ethereum/subs/erc20"
	erc721 "github.com/rhizome-chain/ethereum/subs/erc721"
	"path/filepath"
	
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/config"
	
	// "errors"
	"github.com/rhizome-chain/ethereum/subs"
	"github.com/rhizome-chain/tendermint-daemon/daemon"
	"github.com/rhizome-chain/tendermint-daemon/daemon/worker"
	"github.com/rhizome-chain/tendermint-daemon/types"
)

type EthModule struct {
	modcfg *subs.EthConfig
	manager *subs.EthSubsManager
}

var _ daemon.Module = (*EthModule)(nil)

func (e *EthModule) GetDefaultConfig() types.ModuleConfig {
	config := &subs.EthConfig{}
	return config
}

func (e *EthModule) GetConfig() types.ModuleConfig {
	return e.modcfg
}

func (e *EthModule) Factories() (facs []worker.Factory) {
	return []worker.Factory{e.manager}
}

func (e *EthModule) AddFlags(cmd *cobra.Command) {
	subs.AddEthFlags(cmd)
}

func (e *EthModule) InitFile(config *config.Config) {
	confFilePath := filepath.Join(config.RootDir, "config", "ethereum.toml")
	ethConfig := &subs.EthConfig{}
	err := viper.Unmarshal(ethConfig)
	if err != nil {
		panic("Unmarshal EthConfig" + err.Error())
	}
	
	types.WriteModuleConfigFile(confFilePath, ethConfig)
	fmt.Println("[EthModule] Write EthConfig file:", confFilePath)
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
	
	manager.RegisterLogHandler(&erc20.ERC20LogHandler{})
	manager.RegisterLogHandler(&erc721.ERC721LogHandler{})
	
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
