package main

import (
	"crypto/tls"
	"fmt"
	"github.com/mutalisk999/tls-proxy-go"
	"io"
	"log"
	"net"
	"os"
)

func clientHandler(conn *net.TCPConn, config *tls_proxy_go.ClientConfig) {
	defer conn.Close()

	cert, err := tls.LoadX509KeyPair(config.ClientCert, config.ClientKey)
	if err != nil {
		log.Fatalf("LoadX509KeyPair: %v", err)
		return
	}

	tlsConfig := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
	clientConn, err := tls.Dial("tcp",
		fmt.Sprintf("%s:%d", config.ServerHost, config.ServerPort),
		&tlsConfig)
	if err != nil {
		log.Printf("tls.Dial: %v", err)
		return
	}

	defer clientConn.Close()

	go func() { _, _ = io.Copy(conn, clientConn) }()
	go func() { _, _ = io.Copy(clientConn, conn) }()

	select {}
}

func main() {
	log.SetOutput(os.Stdout)

	config, err := tls_proxy_go.LoadClientConfig()
	if err != nil {
		log.Fatalf("LoadClientConfig: %v", err)
		return
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", config.ListenHost, config.ListenPort))
	if err != nil {
		log.Fatalf("ResolveTCPAddr: %v", err)
		return
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatalf("ListenTCP: %v", err)
		return
	}
	log.Printf("bind on: %v", tcpAddr)

	defer listener.Close()

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Fatalf("AcceptTCP: %v", err)
		}
		log.Printf("accept connection from: %v", conn.RemoteAddr())

		go clientHandler(conn, config)
	}
}
