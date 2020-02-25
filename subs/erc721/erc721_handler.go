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

const erc721Abi = `[{"constant":true,"inputs":[],"name":"name","outputs":[{"name":"_name","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"_tokenId","type":"uint256"}],"name":"getApproved","outputs":[{"name":"_approved","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_to","type":"address"},{"name":"_tokenId","type":"uint256"}],"name":"approve","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"implementsERC721","outputs":[{"name":"_implementsERC721","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"totalSupply","outputs":[{"name":"_totalSupply","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_from","type":"address"},{"name":"_to","type":"address"},{"name":"_tokenId","type":"uint256"}],"name":"transferFrom","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"_owner","type":"address"},{"name":"_index","type":"uint256"}],"name":"tokenOfOwnerByIndex","outputs":[{"name":"_tokenId","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"_tokenId","type":"uint256"}],"name":"ownerOf","outputs":[{"name":"_owner","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"_tokenId","type":"uint256"}],"name":"tokenMetadata","outputs":[{"name":"_infoUrl","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"_owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"_balance","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_owner","type":"address"},{"name":"_tokenId","type":"uint256"},{"name":"_approvedAddress","type":"address"},{"name":"_metadata","type":"string"}],"name":"mint","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"_symbol","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_to","type":"address"},{"name":"_tokenId","type":"uint256"}],"name":"transfer","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"numTokensTotal","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"_owner","type":"address"}],"name":"getOwnerTokens","outputs":[{"name":"_tokenIds","type":"uint256[]"}],"payable":false,"stateMutability":"view","type":"function"},{"anonymous":false,"inputs":[{"indexed":true,"name":"_to","type":"address"},{"indexed":true,"name":"_tokenId","type":"uint256"}],"name":"Mint","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"_from","type":"address"},{"indexed":true,"name":"_to","type":"address"},{"indexed":false,"name":"_tokenId","type":"uint256"}],"name":"Transfer","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"_owner","type":"address"},{"indexed":true,"name":"_approved","type":"address"},{"indexed":false,"name":"_tokenId","type":"uint256"}],"name":"Approval","type":"event"}]`

var (
	erc721TransferSig        = []byte("Transfer(address,address,uint256)")
	erc721ApprovalSig        = []byte("Approval(address,address,uint256)")
	erc721ApprovalAllSig     = []byte("ApprovalForAll(address,address,bool)")
	erc721MintSig            = []byte("Mint(address,uint256)")
	erc721TransferSigHash    = crypto.Keccak256Hash(erc721TransferSig).Hex()
	erc721ApprovalSigHash    = crypto.Keccak256Hash(erc721ApprovalSig).Hex()
	erc721ApprovalAllSigHash = crypto.Keccak256Hash(erc721ApprovalAllSig).Hex()
	erc721MintSigHash        = crypto.Keccak256Hash(erc721MintSig).Hex()
)

// ERC721LogHandler implements LogHandler
type ERC721LogHandler struct {
	erc721Abi *abi.ABI
}

var _ ethtypes.LogHandler = (*ERC721LogHandler)(nil)

type Erc721Event struct {
	Address     string   `json:"addr"`
	From        string   `json:"from"`
	To          string   `json:"to"`
	Type        string   `json:"type"`
	Tokens      *big.Int `json:"Tokens"`
	BlockNumber uint64   `json:"blockNumber"`
	TxIndex     uint     `json:"txIndex"`
}

func init() {
	tdtypes.BasicCdc.RegisterConcrete(Erc721Event{}, "Erc721Event", nil)
}

// Name : erc20
func (handler *ERC721LogHandler) Name() string { return "erc721" }

func NewERC721LogHandler() *ERC721LogHandler {
	handler := &ERC721LogHandler{}
	abi, _ := abi.JSON(strings.NewReader(erc721Abi))
	handler.erc721Abi = &abi
	return handler
}

// HandleLog ..
func (handler *ERC721LogHandler) HandleLog(helper *worker.Helper, elog types.Log) error {
	logHash := elog.Topics[0].Hex()
	
	address := elog.Address.Hex()
	event := Erc721Event{Address: address,
		BlockNumber: elog.BlockNumber, TxIndex: elog.TxIndex}
	
	if len(elog.Topics) > 2 {
		fromAddr := common.HexToAddress(elog.Topics[1].Hex()).Hex()
		toAddr := common.HexToAddress(elog.Topics[2].Hex()).Hex()
		event.From = fromAddr
		event.To = toAddr
	}
	var err error
	switch logHash {
	case erc721TransferSigHash:
		event.Type = "Transfer"
		err = handler.erc721Abi.Unpack(&event, "Transfer", elog.Data)
		
		if err != nil {
			log.Println("[ERROR-ERC721] Unpack Transfer event data ", err)
		}
		
		break
	case erc721ApprovalSigHash:
		event.Type = "Approval"
		err = handler.erc721Abi.Unpack(&event, "Approval", elog.Data)
		if err != nil {
			log.Println("[ERROR-ERC721] Unpack Approval event data ", err)
		}
		break
	case erc721ApprovalAllSigHash:
		event.Type = "ApprovalForAll"
		err = handler.erc721Abi.Unpack(&event, "ApprovalForAll", elog.Data)
		if err != nil {
			log.Println("[ERROR-ERC721] Unpack ApprovalForAll event data ", err)
		}
		break
	
	case erc721MintSigHash:
		event.Type = "Mint"
		err = handler.erc721Abi.Unpack(&event, "Mint", elog.Data)
		if err != nil {
			log.Println("[ERROR-ERC721] Unpack Mint event data ", err)
		}
		break
		
	}
	
	if err == nil {
		// rowID := fmt.Sprintf("%d-%d", elog.BlockNumber, elog.TxIndex)
		
		fmt.Println(" - ", handler.Name(), event)
		
		// err = helper.PutData(rowID, event)
	}
	
	return err
	
	// address := log.Address.Hex()
	// logHash := log.Topics[0].Hex()
	
	// if len(log.Topics) > 2 {
	// 	fromAddr := log.Topics[1].Hex()
	// 	toAddr := log.Topics[2].Hex()
	// 	fmt.Println("LOG ", handler.Name(), "- addr:", address, ", loghash:", logHash, ",fromAddr=", fromAddr, ",toAddr=", toAddr)
	// } else {
	// 	fmt.Println("LOG ", handler.Name(), "- addr:", address, ", loghash:", logHash)
	// }
	
	// return nil
}
