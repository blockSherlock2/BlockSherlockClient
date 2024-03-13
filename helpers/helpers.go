package helpers

import (
	"bufio"
	"bytes"
	"client/models"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
func StrToPK(s string) (privateKey *ecdsa.PrivateKey, address common.Address) {
	privateKey, err := crypto.HexToECDSA(s)
	if err != nil {
		return
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return
	}
	address = crypto.PubkeyToAddress(*publicKeyECDSA)

	return
}
func ReadFromTerminal(startText string) string {
	fmt.Print("\n" + startText + ": ")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(text, "\n", ""), "\r", ""), " ", "")
}

type WalletsPayload struct {
	Wallets []string `json:"wallets"`
}
type ServerResponse struct {
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	ErrorType string      `json:"err_type"`
	IsError   bool        `json:"isErr"`
}

type AddWalletsResponse struct {
	Message string   `json:"message"`
	Wallets []string `json:"wallets"`
}

func SendAddWalletsRequest(wallets []string, userId, apiKey, serverAddr string) (string, bool) {
	client := &http.Client{Timeout: time.Second * time.Duration(5)}
	var msg string
	payload := WalletsPayload{Wallets: wallets}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		msg = "Error creating payload"
		return msg, false
	}
	serverStr := fmt.Sprintf("http://%s/public/addWallets?userId=%s&apiKey=%s", serverAddr, userId, apiKey)
	req, err := http.NewRequest("POST", serverStr, bytes.NewBuffer(payloadBytes))
	if err != nil {
		msg = "Error creating request"
		return msg, false
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		msg = "Error sending request"
		return msg, false

	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		msg = "Error reading response body"
		fmt.Println(msg, err)
		return msg, false
	}
	var s ServerResponse
	json.Unmarshal(body, &s)
	if s.IsError {
		return s.Message, false
	}
	return s.Message, true
}
func JSON(reader io.Reader) (models.Config, error) {
	dec := json.NewDecoder(reader)

	var cfg models.Config
	if err := dec.Decode(&cfg); err != nil {
		return models.Config{}, err
	}
	return cfg, nil
}

func ParseConfigFile(configPath string) *models.Config {
	byteData, err := os.ReadFile(configPath)
	CheckErr(err)
	parsedCfg, err := JSON(strings.NewReader(string(byteData)))
	CheckErr(err)

	return &parsedCfg
}
func AppendFile(fileName string, data string) {
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	if _, err := f.Write([]byte(data)); err != nil {
		panic(err)
	}
	if err := f.Close(); err != nil {
		panic(err)
	}
}

func AddWalletsToFile(wallets []string, walletsFilePath string) {
	walletsStr := strings.Join(wallets, "\n")
	AppendFile(walletsStr, walletsFilePath)
}
