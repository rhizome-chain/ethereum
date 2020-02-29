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

func GetUnmarshal(dataType string) (unmarshal func(value []byte) interface{}, err error) {
	if dataType == "erc20" {
		unmarshal = func(value []byte) interface{} {
			var event erc20.Erc20Event
			types.BasicCdc.UnmarshalBinaryBare(value, &event)
			return event
		}
	} else if dataType == "erc721" {
		unmarshal = func(value []byte) interface{} {
			var event erc721.Erc721Event
			types.BasicCdc.UnmarshalBinaryBare(value, &event)
			return event
		}
	} else {
		err = errors.New("unknown Data Type for unmarshal" + dataType)
		return nil, err
	}
	
	return unmarshal, err
}

func GetJsonStringer(dataType string) (parser func(value []byte) string, err error) {
	unmarshal, err := GetUnmarshal(dataType)
	
	if err == nil {
		parser = func(value []byte) string {
			var event = unmarshal(value)
			var log string
			btz, err := json.Marshal(event)
			if err != nil {
				log = fmt.Sprintf("{\"err\"=\"%s\"}", err.Error())
			} else {
				log = string(btz)
			}
			return log
		}
	}
	
	return parser, err
}

func GetSimpleStringer(dataType string) (parser func(value []byte) string, err error) {
	unmarshal, err := GetUnmarshal(dataType)
	
	if err == nil {
		parser = func(value []byte) string {
			var event = unmarshal(value)
			return fmt.Sprintf("%v", event)
		}
	}
	
	return parser, err
}

func GetNoneStringer(dataType string) (parser func(value []byte) string, err error) {
	return func(value []byte) string {
		return ""
	}, err
}

// Start ..
func (worker *EthLogWorker) Start() (err error) {
	worker.wait = make(chan bool)
	
	var stringer func(value []byte) string
	
	if worker.jobInfo.LogType == "json" {
		stringer, err = GetJsonStringer(worker.jobInfo.DataType)
	} else if worker.jobInfo.LogType == "simple" {
		stringer, err = GetSimpleStringer(worker.jobInfo.DataType)
	} else {
		stringer, err = GetNoneStringer(worker.jobInfo.DataType)
	}
	
	if err != nil {
		worker.helper.Error("[EthLog] get stringer ", err)
		return err
	}
	
	worker.started = true
	
	worker.helper.Info("Start EthLog " + worker.id)
	
	var lastRow string
	worker.helper.GetCheckpoint(&lastRow)
	
	cancel, err := worker.sourceProxy.CollectAndSubscribe("in", lastRow, func(jobID string, topic string, rowID string, value []byte) bool {
		fmt.Println("[EthLog]", jobID, rowID, stringer(value))
		
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
