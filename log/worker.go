package log

import (
	"encoding/json"
	"errors"
	"fmt"
	erc20 "github.com/rhizome-chain/ethereum/subs/erc20"
	erc721 "github.com/rhizome-chain/ethereum/subs/erc721"
	"github.com/rhizome-chain/tendermint-daemon/daemon/worker"
	"github.com/rhizome-chain/tendermint-daemon/types"
)

// EthLogWorker implements worker.Worker
type EthLogWorker struct {
	id          string
	jobInfo     *EthLogJobInfo
	helper      *worker.Helper
	sourceProxy worker.Proxy
	started     bool
}

var _ worker.Worker = (*EthLogWorker)(nil)

// ID ..
func (worker *EthLogWorker) ID() string {
	return worker.id
}

// Start ..
func (worker *EthLogWorker) Start() (err error) {
	worker.started = true
	
	var parser func(value []byte) string
	
	if worker.jobInfo.DataType == "erc20" {
		parser = func(value []byte)string {
			var event erc20.Erc20Event
			types.BasicCdc.UnmarshalBinaryBare(value,&event)
			var log string
			btz, err := json.Marshal(event)
			if err != nil{
				log = err.Error()
			} else {
				log = string(btz)
			}
			return log
		}
	} else if worker.jobInfo.DataType == "721" {
		parser = func(value []byte)string {
			var event erc721.Erc721Event
			types.BasicCdc.UnmarshalBinaryBare(value,&event)
			var log string
			btz, err := json.Marshal(event)
			if err != nil{
				log = err.Error()
			} else {
				log = string(btz)
			}
			return log
		}
	} else {
		err = errors.New("unknown Data Type " + worker.jobInfo.DataType)
		worker.helper.Error("[EthLog] Unknown Data Type", err)
		return err
	}
	
	worker.helper.Info("Start EthLog " + worker.id)
	
	worker.sourceProxy.GetDataList("in", func(jobID string, topic string, rowID string, value []byte) bool {
		fmt.Println("[EthLog] ", jobID, rowID, parser(value))
		return true
	})
	
	
	worker.started = false
	
	return err
}

// Stop ..
func (worker *EthLogWorker) Stop() error {
	worker.started = false
	return nil
}

// IsStarted ..
func (worker *EthLogWorker) IsStarted() bool {
	return worker.started
}
