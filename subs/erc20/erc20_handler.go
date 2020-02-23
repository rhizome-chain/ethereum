package ethereum

import (
	"fmt"
	"log"
	"math/big"
	"strings"
	
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	
	"github.com/rhizome-chain/tendermint-daemon/daemon/worker"
	tdtypes "github.com/rhizome-chain/tendermint-daemon/types"
	
	ethtypes "github.com/rhizome-chain/ethereum/subs/types"
)

// ERC20LogHandler implements LogHandler
type ERC20LogHandler struct {
	erc20Abi *abi.ABI
}

type erc20Event struct {
	Address     string   `json:"addr"`
	From        string   `json:"from"`
	To          string   `json:"to"`
	Type        string   `json:"type"`
	Tokens      *big.Int `json:"Tokens"`
	BlockNumber uint64   `json:"blockNumber"`
	TxIndex     uint     `json:"txIndex"`
}

func init() {
	tdtypes.BasicCdc.RegisterConcrete(erc20Event{}, "eth/erc20Event", nil)
}

var _ ethtypes.LogHandler = (*ERC20LogHandler)(nil)

func NewERC20LogHandler() *ERC20LogHandler {
	handler := &ERC20LogHandler{}
	abi, _ := abi.JSON(strings.NewReader(erc20Abi))
	handler.erc20Abi = &abi
	return handler
}

// Name : erc20
func (handler *ERC20LogHandler) Name() string { return "erc20" }

const erc20Abi = `[{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"tokenOwner","type":"address"},{"indexed":true,"internalType":"address","name":"spender","type":"address"},{"indexed":false,"internalType":"uint256","name":"tokens","type":"uint256"}],"name":"Approval","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"from","type":"address"},{"indexed":true,"internalType":"address","name":"to","type":"address"},{"indexed":false,"internalType":"uint256","name":"tokens","type":"uint256"}],"name":"Transfer","type":"event"}]`

var (
	erc20TransferSig     = []byte("Transfer(address,address,uint256)")
	erc20ApprovalSig     = []byte("Approval(address,address,uint256)")
	erc20TransferSigHash = crypto.Keccak256Hash(erc20TransferSig).Hex()
	erc20ApprovalSigHash = crypto.Keccak256Hash(erc20ApprovalSig).Hex()
)

// HandleLog ..
func (handler *ERC20LogHandler) HandleLog(helper *worker.Helper, elog types.Log) error {
	logHash := elog.Topics[0].Hex()
	
	address := elog.Address.Hex()
	fromAddr := common.HexToAddress(elog.Topics[1].Hex()).Hex()
	toAddr := common.HexToAddress(elog.Topics[2].Hex()).Hex()
	
	event := erc20Event{Address: address, From: fromAddr, To: toAddr,
		BlockNumber: elog.BlockNumber, TxIndex: elog.TxIndex}
	
	var err error
	switch logHash {
	case erc20TransferSigHash:
		event.Type = "Transfer"
		err = handler.erc20Abi.Unpack(&event, "Transfer", elog.Data)
		if err != nil {
			log.Println("[ERROR-ERC20] Unpack Transfer event data ", err)
		}
		break
	case erc20ApprovalSigHash:
		event.Type = "Approval"
		err = handler.erc20Abi.Unpack(&event, "Approval", elog.Data)
		if err != nil {
			log.Println("[ERROR-ERC20] Unpack Approval event data ", err)
		}
		break
	}
	
	if err == nil {
		rowID := fmt.Sprintf("%d-%d", elog.BlockNumber, elog.TxIndex)
		
		// fmt.Println(" - ", helper.ID(), rowID)
		
		err = helper.PutObject("in", rowID, event)
		if err != nil {
			helper.Error("Put Eth Event", "jobID", helper.ID(), "rowID", rowID, "err", err)
		}
	}
	
	return err
}
