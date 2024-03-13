package eth

import (
	"client/helpers"
	"client/models"
	"fmt"
	"os"
	"path"
	"strings"
	"sync"
)

func LoadWallets(memory *sync.Map) ([]string, string) {
	walletsFolder := "w"
	walletFiles, err := os.ReadDir(walletsFolder)
	helpers.CheckErr(err)
	walletsFilePath := path.Join(walletsFolder, walletFiles[0].Name())
	walletsByte, err := os.ReadFile(walletsFilePath)
	helpers.CheckErr(err)
	walletsList := strings.Split(strings.ReplaceAll(string(walletsByte), "\r", ""), "\n")
	loadedWallets := []string{}
	for _, walletString := range walletsList {
		keyChain := strings.Split(walletString, " ")
		if len(keyChain) != 2 {
			continue
		}
		privateKeyStr := keyChain[1]
		privateKey, address := helpers.StrToPK(strings.Replace(privateKeyStr, "0x", "", 1))
		if privateKey == nil {
			fmt.Println("invalid privateKey", privateKey)
			continue
		}
		loadedWallets = append(loadedWallets, address.String())
		memory.Store(address, models.Wallet{PrivateKey: privateKey})
	}
	return loadedWallets, walletsFilePath
}
