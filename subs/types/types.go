package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	
	"github.com/rhizome-chain/tendermint-daemon/daemon/worker"
)

// LogHandler handler for ethereum log
type LogHandler interface {
	Name() string
	HandleLog(helper *worker.Helper, log types.Log) error
}

// EthSubsJobInfo ..
type EthSubsJobInfo struct {
	Handler           string   `json:"handler"`
	CAs               []string `json:"cas"`
	From              uint64   `json:"from"`
	contractAddresses []common.Address
}

func (info *EthSubsJobInfo) GetContractAddresses() []common.Address {
	if info.contractAddresses == nil {
		if info.CAs != nil {
			addrs := []common.Address{}
			
			for _, ca := range info.CAs {
				addr := common.HexToAddress(ca)
				addrs = append(addrs, addr)
			}
			
			info.contractAddresses = addrs
		}
	}
	return info.contractAddresses
}

// BlockCheckPoint ..
type BlockCheckPoint struct {
	BlockNumber uint64
	Index       uint
}
