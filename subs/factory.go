package subs

import (
	"errors"
	
	"encoding/json"
	
	ethtypes "github.com/rhizome-chain/ethereum/subs/types"
	"github.com/rhizome-chain/tendermint-daemon/daemon/worker"
)

const (
	FactoryName = "eth_subs"
)
// EthSubsManager implements worker.Factory, name eth_subs
type EthSubsManager struct {
	networkURL string
	handlers   map[string]ethtypes.LogHandler
}

var _ worker.Factory = (*EthSubsManager)(nil)

// NewEthSubsManager ..
func NewEthSubsManager(networkURL string) *EthSubsManager {
	manager := EthSubsManager{networkURL: networkURL, handlers: make(map[string]ethtypes.LogHandler)}
	return &manager
}

// RegisterLogHandler ..
func (manager *EthSubsManager) RegisterLogHandler(handler ethtypes.LogHandler) {
	manager.handlers[handler.Name()] = handler
}

// Name implements worker.Factory.Space
func (manager *EthSubsManager) Space() string {
	return "ethereum"
}

// Name implements worker.Factory.Name
func (manager *EthSubsManager) Name() string {
	return FactoryName
}

// NewWorker implements worker.Factory.NewWorker
func (manager *EthSubsManager) NewWorker(helper *worker.Helper) (wroker worker.Worker, err error) {
	jobInfo := new(ethtypes.EthSubsJobInfo)
	json.Unmarshal(helper.Job().Data, jobInfo)
	
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
