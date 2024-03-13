package server

import (
	"client/helpers"
	"client/models"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Memory struct {
	Wallets         *sync.Map
	ApiKey          string
	Port            string
	UserId          string
	ServerAddr      string
	WalletsFilePath string
}
type SignatureReqData struct {
	Addr    common.Address
	Tx      *types.Transaction
	ChainId *big.Int
}

type Response struct {
	Tx *types.Transaction
}

func (mem *Memory) signHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	var signatureReqData SignatureReqData
	fmt.Println(r.Body, "body")
	err := json.NewDecoder(r.Body).Decode(&signatureReqData)
	helpers.CheckErr(err)
	targetAccount, ok := mem.Wallets.Load(signatureReqData.Addr)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	account := targetAccount.(models.Wallet)
	fmt.Println(signatureReqData.Addr, signatureReqData.Tx)
	fmt.Println(account, signatureReqData.Tx.GasPrice())
	var signedTx *types.Transaction
	switch signatureReqData.Tx.Type() {
	case 0:
		signedTx, err = types.SignTx(signatureReqData.Tx, types.NewEIP155Signer(signatureReqData.ChainId), account.PrivateKey)
	case 2:
		signedTx, err = types.SignTx(signatureReqData.Tx, types.NewLondonSigner(signatureReqData.ChainId), account.PrivateKey)
	}
	helpers.CheckErr(err)
	responseBody := Response{Tx: signedTx}
	responceByte, err := json.Marshal(responseBody)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(responceByte)
}
func (mem *Memory) FindAddress(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	var data map[string]string
	fmt.Println(r.Body, "body")
	err := json.NewDecoder(r.Body).Decode(&data)
	helpers.CheckErr(err)
	targetAccount, ok := mem.Wallets.Load(common.HexToAddress(data["address"]))
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	fmt.Println(targetAccount)
	w.WriteHeader(http.StatusOK)

}
func ValidateWallets(allWallets sync.Map, wallets []string) ([]string, []string) {
	var validWalletsForSend []string
	var validWalletsForWrite []string
	var uniqWallets sync.Map

	for _, addrPk := range wallets {

		addrPkSplit := strings.Split(addrPk, " ")
		if len(addrPkSplit) != 2 {
			continue
		}
		_, address := helpers.StrToPK(addrPkSplit[1])

		_, inLocal := uniqWallets.Load(address)
		_, inGlobal := allWallets.Load(address)
		if inLocal || inGlobal {
			continue
		}
		validWalletsForSend = append(validWalletsForSend, address.String())
		uniqWallets.Store(address, address)
	}

	return validWalletsForSend, validWalletsForWrite
}

type ReqAddWallets struct {
	Wallets []string `json:"wallets"`
}

func (mem *Memory) AddWalletHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	var req ReqAddWallets
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return
	}
	validWalletsForSend, validWalletsForWrite := ValidateWallets(*mem.Wallets, req.Wallets)

	msg, ok := helpers.SendAddWalletsRequest(validWalletsForSend, mem.UserId, mem.ApiKey, mem.ServerAddr)
	fmt.Println(msg)
	if !ok {
		return
	}
	helpers.AddWalletsToFile(validWalletsForWrite, mem.WalletsFilePath)

}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.WriteHeader(http.StatusOK)
}
