package log

import (
	"github.com/rhizome-chain/ethereum/common"
	
	"encoding/json"
	
	"github.com/rhizome-chain/tendermint-daemon/daemon/worker"
)

const (
	FactoryName = "eth_log"
)
// EthLogFactory implements worker.Factory, name eth_subs
type EthLogFactory struct {
	config common.EthConfig
}

var _ worker.Factory = (*EthLogFactory)(nil)

// NewEthLogFactory ..
func NewEthLogFactory(config common.EthConfig) *EthLogFactory {
	factory := EthLogFactory{config:config}
	return &factory
}

// Space implements worker.Factory.Space
func (factory *EthLogFactory) Space() string {
	return "ethereum"
}

// Name implements worker.Factory.Name
func (factory *EthLogFactory) Name() string {
	return FactoryName
}

// NewWorker implements worker.Factory.NewWorker
func (factory *EthLogFactory) NewWorker(helper *worker.Helper) (wroker worker.Worker, err error) {
	jobInfo := new(EthLogJobInfo)
	err = json.Unmarshal(helper.Job().Data, jobInfo)
	
	if err != nil {
		return nil, err
	}
	
	proxy, err := helper.NewWorkerProxy(jobInfo.SourceJobID)
	
	if err != nil {
		return nil, err
	}
	
	logWorker := EthLogWorker{id: helper.ID(), jobInfo: jobInfo, helper: helper, sourceProxy:proxy}
	
	return &logWorker, err
}
