package subs

import (
	"context"
	"fmt"
	"github.com/rhizome-chain/tendermint-daemon/tm/client"
	"math/big"
	
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	
	ethtypes "github.com/rhizome-chain/ethereum/subs/types"
	"github.com/rhizome-chain/tendermint-daemon/daemon/worker"
)

// EthSubscriber implements worker.Worker
type EthSubscriber struct {
	id         string
	client     *ethclient.Client
	networkURL string
	jobInfo    *ethtypes.EthSubsJobInfo
	helper     *worker.Helper
	started    bool
	handler    ethtypes.LogHandler
}

var _ worker.Worker = (*EthSubscriber)(nil)

// ID ..
func (subscriber *EthSubscriber) ID() string {
	return subscriber.id
}

// Start ..
func (subscriber *EthSubscriber) Start() error {
	if subscriber.client != nil {
		subscriber.client.Close()
	}
	
	subscriber.started = true
	
	client, err := ethclient.Dial(subscriber.networkURL)
	if err != nil {
		subscriber.helper.Error("[ERROR] Cannot Connect to ", "network_url", subscriber.networkURL, "err", err)
		subscriber.started = false
		return err
	}
	subscriber.client = client
	
	subscriber.helper.Info("[Debug] ETH Subs :", "CAs", subscriber.jobInfo.CAs, ", from:", subscriber.jobInfo.From)
	checkPoint := &ethtypes.BlockCheckPoint{}
	subscriber.helper.GetCheckpoint(checkPoint)
	
	if subscriber.jobInfo.From > 0 && checkPoint.BlockNumber == 0 {
		checkPoint.BlockNumber = subscriber.jobInfo.From
	}
	
	if checkPoint.BlockNumber > 0 {
		subscriber.collect(checkPoint)
	}
	
	subscriber.subscribe(checkPoint)
	subscriber.helper.Info("[WARN] ETH Subs Ends. ", "job_id", subscriber.ID())
	
	subscriber.started = false
	
	return nil
}

func (subscriber *EthSubscriber) handleLog(elog types.Log, checkPoint *ethtypes.BlockCheckPoint) bool {
	if !subscriber.IsStarted() {
		return false
	}
	err := subscriber.handler.HandleLog(subscriber.helper, elog)
	
	if err != nil {
		if client.IsErrTxInCache(err) {
			subscriber.helper.Debug("[ETH-LogHandler] Skip TxInCache",
				"job_id=", subscriber.ID(), "BlockNumber=", elog.BlockNumber, "Index=", elog.Index)
		} else {
			subscriber.helper.Error("[ERROR-ETH-LogHandler] ", "job_id", subscriber.ID(), "err", err)
		}
		
		// panic(fmt.Sprintf("handleLog %s", subscriber.ID()))
	} else {
		checkPoint.BlockNumber = elog.BlockNumber
		checkPoint.Index = elog.Index
		subscriber.helper.PutCheckpoint(checkPoint)
	}
	
	return true
}

func (subscriber *EthSubscriber) subscribe(checkPoint *ethtypes.BlockCheckPoint) {
	if !subscriber.IsStarted() {
		return
	}
	
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(checkPoint.BlockNumber)),
		Addresses: subscriber.jobInfo.GetContractAddresses(),
	}
	
	logs := make(chan types.Log)
	sub, err := subscriber.client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		subscriber.helper.Error("[ERROR] SubscribeFilterLogs ", "job_id", subscriber.ID(), "err", err)
		return
	}
	
	defer sub.Unsubscribe()
	
	subscriber.helper.Info(fmt.Sprintf("[EthSubscriber %s] starts subscribing. ", subscriber.ID()))
	for subscriber.started {
		select {
		case err := <-sub.Err():
			if !subscriber.started {
				break
			}
			subscriber.helper.Error("[ERROR] Eth Sub ", "job_id", subscriber.ID(), "err", err)
		case vLog := <-logs:
			if !subscriber.started {
				subscriber.helper.Info("[WARN] Eth Subscriber Stops .. ", "job_id", subscriber.ID())
				break
			}
			
			if !subscriber.handleLog(vLog, checkPoint) {
				break
			}
		}
	}
}

func (subscriber *EthSubscriber) collect(checkPoint *ethtypes.BlockCheckPoint) {
	remained := subscriber.collectStep(checkPoint, 10, 0)
	
	for remained > 0 {
		remained = subscriber.collectStep(checkPoint, 10, 1)
	}
}

func (subscriber *EthSubscriber) collectStep(checkPoint *ethtypes.BlockCheckPoint, step uint64, offset uint64) (remained int64) {
	if !subscriber.IsStarted() {
		return 0
	}
	
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(checkPoint.BlockNumber + offset)),
		ToBlock:   big.NewInt(int64(checkPoint.BlockNumber + step)),
		Addresses: subscriber.jobInfo.GetContractAddresses(),
	}
	
	logs, err := subscriber.client.FilterLogs(context.Background(), query)
	if err != nil {
		subscriber.helper.Error("[ERROR-ETH-Subs]", err)
		return
	}
	
	subscriber.helper.Info(fmt.Sprintf("[EthSubscriber %s] collect old TX form %d to %d",
		subscriber.ID(), query.FromBlock, query.ToBlock))
	
	beenHandled := false
	if len(logs) > 0 {
		if offset == 0 {
			skippedBlock := uint64(0)
			skippedCount := 0
			for _, vLog := range logs {
				if vLog.BlockNumber == checkPoint.BlockNumber && vLog.Index <= checkPoint.Index {
					//fmt.Println("------ Skip Log : Block - ", vLog.BlockNumber, ", Index - ", vLog.Index, "<=", checkPoint.Index)
					skippedCount++
					skippedBlock = checkPoint.BlockNumber
					continue
				}
				subscriber.helper.Info(fmt.Sprintf("Skip %d logs in block %d", skippedCount, skippedBlock))
				break
			}
			
			for _, vLog := range logs {
				if !subscriber.handleLog(vLog, checkPoint) {
					return 0
				}
				beenHandled = true
			}
		} else {
			for _, vLog := range logs {
				if !subscriber.handleLog(vLog, checkPoint) {
					return 0
				}
				beenHandled = true
			}
		}
	}
	
	if !beenHandled {
		checkPoint.BlockNumber = checkPoint.BlockNumber + step
		checkPoint.Index = 0
	}
	
	header, err := subscriber.client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		// panic("Cannot get Eth HeaderByNumber:" + err.Error())
		subscriber.helper.Error(fmt.Sprintf("Job[%s] cannot get Eth HeaderByNumber", subscriber.helper.ID()), err)
		return 0
	}
	
	var curBlock int64
	
	if header != nil {
		curBlock = header.Number.Int64()
	}
	
	remained = curBlock - int64(checkPoint.BlockNumber)
	
	return remained
}

// Stop ..
func (subscriber *EthSubscriber) Stop() error {
	if subscriber.client != nil {
		subscriber.client.Close()
	}
	subscriber.client = nil
	subscriber.started = false
	return nil
}

// IsStarted ..
func (subscriber *EthSubscriber) IsStarted() bool {
	return subscriber.started
}
