package tls_proxy_go

import (
	"encoding/json"
	"os"
)

type ClientConfig struct {
	ListenHost string `json:"listen_host"`
	ListenPort uint16 `json:"listen_port"`
	ServerHost string `json:"server_host"`
	ServerPort uint16 `json:"server_port"`
	ClientKey  string `json:"client_key"`
	ClientCert string `json:"client_cert"`
}

type ServerConfig struct {
	ListenHost string `json:"listen_host"`
	ListenPort uint16 `json:"listen_port"`
	ServerKey  string `json:"server_key"`
	ServerCert string `json:"server_cert"`
}

func LoadClientConfig() (*ClientConfig, error) {
	f, err := os.OpenFile("client.json", os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	configBytes := make([]byte, 4096)
	n, err := f.Read(configBytes)
	if err != nil {
		return nil, err
	}
	clientConfig := ClientConfig{}
	err = json.Unmarshal(configBytes[:n], &clientConfig)
	if err != nil {
		return nil, err
	}
	return &clientConfig, nil
}

func LoadServerConfig() (*ServerConfig, error) {
	f, err := os.OpenFile("server.json", os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	configBytes := make([]byte, 4096)
	n, err := f.Read(configBytes)
	if err != nil {
		return nil, err
	}
	serverConfig := ServerConfig{}
	err = json.Unmarshal(configBytes[:n], &serverConfig)
	if err != nil {
		return nil, err
	}
	return &serverConfig, nil
}
