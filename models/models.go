package models

import (
	"crypto/ecdsa"
)

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
}
type Config struct {
	ApiKey     string `json:"api_key"`
	Port       string `json:"port"`
	UserId     string `json:"user_id"`
	ServerAddr string `json:"server_addr"`
}
