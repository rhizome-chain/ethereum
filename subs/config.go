package subs

import (
	"github.com/spf13/cobra"
	
	"github.com/rhizome-chain/tendermint-daemon/types"
)

const (
	flagEthUrl = "eth_network_url"
)

type EthConfig struct {
	NetworkURL string `mapstructure:"eth_network_url"`
}

var _ types.ModuleConfig = (*EthConfig)(nil)

func AddEthFlags(cmd *cobra.Command) {
	cmd.Flags().String(flagEthUrl, "", "Ethereum network url : wss://mainnet.infura.io/v3/PROJECT-ID")
}

func (e EthConfig) GetTemplate() string {
	return templateText
}

var templateText = `#Ethereum subscriber config
# This is a TOML config file.
# For more information, see https://github.com/toml-lang/toml

eth_network_url = "{{ .NetworkURL }}"
`
