package subs

import "github.com/rhizome-chain/tendermint-daemon/types"

type EthConfig struct {
	NetworkURL string
}

var _ types.ModuleConfig = (*EthConfig)(nil)



func (e EthConfig) GetTemplate() string {
	return ``
}
