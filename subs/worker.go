package subs

import (
	"context"
	"fmt"
	"math/big"
	
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	
	"github.com/rhizome-chain/tendermint-daemon/daemon/worker"
)

// EthSubscriber implements worker.Worker
type EthSubscriber struct {
	id         string
	client     *ethclient.Client
	networkURL string
	jobInfo    *EthSubsJobInfo
	helper     *worker.Helper
	started    bool
	handler    LogHandler
}

var _ worker.Worker = (*EthSubscriber)(nil)

// LogHandler handler for ethereum log
type LogHandler interface {
	Name() string
	HandleLog(helper *worker.Helper, log types.Log) error
}

// EthSubsJobInfo ..
type EthSubsJobInfo struct {
	Handler           string   `json:"handler"`
	CAs               []string `json:"cas"`
	contractAddresses []common.Address
	From              uint64 `json:"from"`
}

// BlockCheckPoint ..
type BlockCheckPoint struct {
	BlockNumber uint64
	Index       uint
}

// ID ..
func (subscriber *EthSubscriber) ID() string {
	return subscriber.id
}

// Start ..
func (subscriber *EthSubscriber) Start() error {
	if subscriber.client != nil {
		subscriber.client.Close()
	}
	client, err := ethclient.Dial(subscriber.networkURL)
	if err != nil {
		subscriber.helper.Error("[ERROR] Cannot Connect to ", "network_url", subscriber.networkURL, "err", err)
		return err
	}
	subscriber.client = client
	
	subscriber.helper.Info("[Debug] ETH Subs :", "CAs", subscriber.jobInfo.CAs, ", from:", subscriber.jobInfo.From)
	checkPoint := &BlockCheckPoint{}
	subscriber.helper.GetCheckpoint(checkPoint)
	
	go func() {
		if checkPoint.BlockNumber > 0 {
			subscriber.collect(checkPoint)
		}
		
		subscriber.subscribe(checkPoint)
		subscriber.helper.Info("[WARN] ETH Subs Ends. ", "job_id", subscriber.ID())
	}()
	return nil
}

func (subscriber *EthSubscriber) handleLog(elog types.Log, checkPoint *BlockCheckPoint) {
	err := subscriber.handler.HandleLog(subscriber.helper, elog)
	if err != nil {
		subscriber.helper.Error("[FATAL-ETH-LogHandler] ", "job_id", subscriber.ID(), "err", err)
	}
	checkPoint.BlockNumber = elog.BlockNumber
	checkPoint.Index = elog.Index
	subscriber.helper.PutCheckpoint(checkPoint)
}

func (subscriber *EthSubscriber) subscribe(checkPoint *BlockCheckPoint) {
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(checkPoint.BlockNumber)),
		Addresses: subscriber.jobInfo.contractAddresses,
	}
	
	logs := make(chan types.Log)
	sub, err := subscriber.client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		subscriber.helper.Error("[ERROR] SubscribeFilterLogs ", "job_id", subscriber.ID(), "err", err)
		return
	}
	
	defer sub.Unsubscribe()
	
	subscriber.started = true
	
	subscriber.helper.Info(fmt.Sprintf("[EthSubscriber %s] starts subscribing. ",subscriber.ID()))
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
			
			// fmt.Printf("Sub Log Block Number: %d:%d  Addr: %s\n", vLog.BlockNumber, vLog.Index, vLog.Address.Hex())
			
			subscriber.handleLog(vLog, checkPoint)
		}
	}
	
	subscriber.started = false
}

func (subscriber *EthSubscriber) collect(checkPoint *BlockCheckPoint) {
	remained := subscriber.collectStep(checkPoint, 10, 0)
	
	for remained > 0 {
		remained = subscriber.collectStep(checkPoint, 10, 1)
	}
}

func (subscriber *EthSubscriber) collectStep(checkPoint *BlockCheckPoint, step uint64, offset uint64) (remained int64) {
	if subscriber.client == nil {
		return 0
	}
	
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(checkPoint.BlockNumber + offset)),
		ToBlock:   big.NewInt(int64(checkPoint.BlockNumber + step)),
		Addresses: subscriber.jobInfo.contractAddresses,
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
			for _, vLog := range logs {
				if vLog.BlockNumber == checkPoint.BlockNumber && vLog.Index <= checkPoint.Index {
					fmt.Println("------ Skip Log : Block - ", vLog.BlockNumber, ", Index - ", vLog.Index, "<=", checkPoint.Index)
					continue
				}
				subscriber.handleLog(vLog, checkPoint)
				beenHandled = true
			}
		} else {
			for _, vLog := range logs {
				subscriber.handleLog(vLog, checkPoint)
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
		panic("Cannot get Eth HeaderByNumber")
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
