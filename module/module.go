package module

import (
	"fmt"
	"path/filepath"
	
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/config"
	
	// "errors"
	"github.com/rhizome-chain/ethereum/subs"
	erc20 "github.com/rhizome-chain/ethereum/subs/erc20"
	erc721 "github.com/rhizome-chain/ethereum/subs/erc721"
	"github.com/rhizome-chain/tendermint-daemon/daemon"
	"github.com/rhizome-chain/tendermint-daemon/daemon/worker"
	"github.com/rhizome-chain/tendermint-daemon/types"
)

const (
	flagEthUrl = "eth_network_url"
)

type EthModule struct {
}

var _ daemon.Module = (*EthModule)(nil)

func (e EthModule) GetDefaultConfig() types.ModuleConfig {
	config := &subs.EthConfig{}
	return config
}

func (e EthModule) Factories() (facs []worker.Factory) {
	return nil
}

func (e EthModule) AddFlags(cmd *cobra.Command) {
	cmd.Flags().String(flagEthUrl, "", "Ethereum network url : wss://mainnet.infura.io/v3/PROJECT-ID")
}

func (e EthModule) InitFile(config *config.Config) {
	confFilePath := filepath.Join(config.RootDir, "config", "ethereum.toml")
	ethConfig := &subs.EthConfig{}
	err := viper.Unmarshal(ethConfig)
	if err != nil {
		panic("Unmarshal EthConfig" + err.Error())
	}
	fmt.Println(ethConfig)
	types.WriteModuleConfigFile(confFilePath,ethConfig)
	fmt.Println("Write EthConfig file:", confFilePath)
}

func (e EthModule) BeforeDaemonStarting(cmd *cobra.Command, dm *daemon.Daemon, moduleConfig types.ModuleConfig) {
	ethConfig := moduleConfig.(*subs.EthConfig)
	ethUrl := ethConfig.NetworkURL
	
	if len(ethUrl) == 0 {
		ethUrl2, err := cmd.Flags().GetString(flagEthUrl)
		if err != nil {
			panic("run flag ethereum network url " + err.Error())
		}
		ethUrl = ethUrl2
		ethConfig.NetworkURL = ethUrl
	}
	
	if len(ethUrl) == 0 {
		panic("Ethereum network url is not set.")
	}
	
	manager := subs.NewEthSubsManager(ethUrl)
	
	manager.RegisterLogHandler(&erc20.ERC20LogHandler{})
	manager.RegisterLogHandler(&erc721.ERC721LogHandler{})
	
	dm.RegisterWorkerFactory(manager)
}

func (e EthModule) AfterDaemonStarted(dm *daemon.Daemon) {
	dm.GetContext().Info("EthModule starts")
}
