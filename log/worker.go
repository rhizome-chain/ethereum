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
	wait        chan bool
}

var _ worker.Worker = (*EthLogWorker)(nil)

// ID ..
func (worker *EthLogWorker) ID() string {
	return worker.id
}

func GetParser(dataType string) (parser func(value []byte) string, err error) {
	if dataType == "erc20" {
		parser = func(value []byte) string {
			var event erc20.Erc20Event
			types.BasicCdc.UnmarshalBinaryBare(value, &event)
			var log string
			btz, err := json.Marshal(event)
			if err != nil {
				log = err.Error()
			} else {
				log = string(btz)
			}
			return log
		}
	} else if dataType == "721" {
		parser = func(value []byte) string {
			var event erc721.Erc721Event
			types.BasicCdc.UnmarshalBinaryBare(value, &event)
			var log string
			btz, err := json.Marshal(event)
			if err != nil {
				log = err.Error()
			} else {
				log = string(btz)
			}
			return log
		}
	} else {
		err = errors.New("unknown Data Type " + dataType)
		return nil, err
	}
	
	return parser, err
}

// Start ..
func (worker *EthLogWorker) Start() (err error) {
	worker.wait = make(chan bool)
	parser, err := GetParser(worker.jobInfo.DataType)
	
	if err != nil {
		worker.helper.Error("[EthLog] get parser ", err)
		return err
	}
	
	worker.started = true
	
	worker.helper.Info("Start EthLog " + worker.id)
	
	var lastRow string
	worker.helper.GetCheckpoint(&lastRow)
	
	cancel, err := worker.sourceProxy.CollectAndSubscribe("in", lastRow, func(jobID string, topic string, rowID string, value []byte) bool {
		fmt.Println("[EthLog]", jobID, rowID, parser(value))
		worker.helper.PutCheckpoint(rowID)
		worker.started = false
		return true
	})
	
	if err != nil {
		worker.started = false
		worker.helper.Error("[EthLog] fail subscribe ", err)
		return err
	}
	
	defer cancel()
	
	<-worker.wait
	
	worker.helper.Info(fmt.Sprintf("[EthLog] %s Subscribe ends. ", worker.id))
	
	return err
}

// Stop ..
func (worker *EthLogWorker) Stop() error {
	worker.started = false
	worker.wait <- false
	worker.helper.Info(fmt.Sprintf("[EthLog] Stoping Log worker  %s . ", worker.id))
	return nil
}

// IsStarted ..
func (worker *EthLogWorker) IsStarted() bool {
	return worker.started
}
