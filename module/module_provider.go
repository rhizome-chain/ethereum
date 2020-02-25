package module

import (
	"fmt"
	common2 "github.com/rhizome-chain/ethereum/common"
	"path/filepath"
	
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	
	cfg "github.com/tendermint/tendermint/config"
	
	"github.com/rhizome-chain/tendermint-daemon/daemon"
	"github.com/rhizome-chain/tendermint-daemon/daemon/common"
	"github.com/rhizome-chain/tendermint-daemon/types"
)

type EthModuleProvider struct {
}

func (e EthModuleProvider) NewModule(tmCfg *cfg.Config, config common.DaemonConfig) daemon.Module {
	return &EthModule{}
}

func (e *EthModuleProvider) GetDefaultConfig() types.ModuleConfig {
	config := &common2.EthConfig{}
	return config
}

func (e *EthModuleProvider) AddFlags(cmd *cobra.Command) {
	common2.AddEthFlags(cmd)
}

func (e *EthModuleProvider) InitFile(config *cfg.Config) {
	confFilePath := filepath.Join(config.RootDir, "config", "ethereum.toml")
	ethConfig := &common2.EthConfig{}
	err := viper.Unmarshal(ethConfig)
	if err != nil {
		panic("Unmarshal EthConfig" + err.Error())
	}
	
	types.WriteModuleConfigFile(confFilePath, ethConfig)
	fmt.Println("[EthModule] Write EthConfig file:", confFilePath)
}

var _ daemon.ModuleProvider = (*EthModuleProvider)(nil)
