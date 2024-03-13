package main

import (
	"client/eth"
	"client/helpers"
	"client/server"
	"fmt"
	"os"
	"strings"
	"sync"
)

func main() {
	var wallets sync.Map
	config := helpers.ParseConfigFile("config.json")
	fmt.Println(config, wallets)
	loadedWallets, walletsFilePath := eth.LoadWallets(&wallets)
	memory := server.Memory{
		Wallets:         &wallets,
		ApiKey:          config.ApiKey,
		Port:            config.Port,
		ServerAddr:      config.ServerAddr,
		UserId:          config.UserId,
		WalletsFilePath: walletsFilePath,
	}
	if len(loadedWallets) == 0 {
		fmt.Println("no wallets in './w' folder")
		return
	}
	walletsFileName := fmt.Sprintf("%d.txt", len(loadedWallets))

	os.WriteFile(walletsFileName, []byte(strings.Join(loadedWallets, "\n")), 0644)

	server.ServerStart(memory.Port, 1000, 1000, &memory, memory.ApiKey)
}
