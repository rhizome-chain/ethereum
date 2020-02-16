package subs

import (
	"errors"
	
	"encoding/json"
	
	"github.com/ethereum/go-ethereum/common"
	"github.com/rhizome-chain/tendermint-daemon/daemon/worker"
)

// EthSubsManager implements worker.Factory, name eth_subs
type EthSubsManager struct {
	networkURL string
	handlers   map[string]LogHandler
}

var _ worker.Factory = (*EthSubsManager)(nil)

// NewEthSubsManager ..
func NewEthSubsManager(networkURL string) *EthSubsManager {
	manager := EthSubsManager{networkURL: networkURL, handlers: make(map[string]LogHandler)}
	return &manager
}

// RegisterLogHandler ..
func (manager *EthSubsManager) RegisterLogHandler(handler LogHandler) {
	manager.handlers[handler.Name()] = handler
}

// Name implements worker.Factory.Space
func (manager *EthSubsManager) Space() string {
	return "ethereum"
}

// Name implements worker.Factory.Name
func (manager *EthSubsManager) Name() string {
	return "eth_subs"
}

// NewWorker implements worker.Factory.NewWorker
func (manager *EthSubsManager) NewWorker(helper *worker.Helper) (wroker worker.Worker, err error) {
	jobInfo := new(EthSubsJobInfo)
	json.Unmarshal(helper.Job().Data, jobInfo)
	
	addrs := []common.Address{}
	
	for _, ca := range jobInfo.CAs {
		addr := common.HexToAddress(ca)
		addrs = append(addrs, addr)
	}
	
	jobInfo.contractAddresses = addrs
	
	handler := manager.handlers[jobInfo.Handler]
	
	if handler == nil {
		err = errors.New("Unknown Log Handler " + jobInfo.Handler)
		helper.Error("[ERROR-EthSubsManager] Unknown Log Handler ", jobInfo.Handler)
		return nil, err
	}
	
	subscriber := EthSubscriber{id: helper.ID(), jobInfo: jobInfo, networkURL: manager.networkURL,
		helper: helper, handler: handler}
	
	return &subscriber, err
}
