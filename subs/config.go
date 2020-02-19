package subs

import "github.com/rhizome-chain/tendermint-daemon/types"

type EthConfig struct {
	NetworkURL string `mapstructure:"eth_network_url"`
}

var _ types.ModuleConfig = (*EthConfig)(nil)



func (e EthConfig) GetTemplate() string {
	return templateText
}

var templateText = `#Ethereum subscriber config
# This is a TOML config file.
# For more information, see https://github.com/toml-lang/toml

eth_network_url = "{{ .NetworkURL }}"
`