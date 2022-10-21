package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/libp2p/go-reuseport"
	"github.com/mutalisk999/tls-proxy-go"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"
)

func clientHandler(conn net.Conn, config *tls_proxy_go.ClientConfig) {
	defer conn.Close()

	cert, err := tls.LoadX509KeyPair(
		config.ClientCert,
		config.ClientKey,
	)
	if err != nil {
		log.Fatalf("LoadX509KeyPair: %v", err)
		return
	}

	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(config.CACert)
	if err != nil {
		log.Fatalf("ioutil.ReadFile CACert: %v", err)
		return
	}

	ok := certPool.AppendCertsFromPEM(ca)
	if !ok {
		log.Fatalf("AppendCertsFromPEM")
		return
	}

	tlsConfig := tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
		RootCAs:            certPool,
	}

	dialer := net.Dialer{Timeout: 5 * time.Second}
	tcpAddr := fmt.Sprintf("%s:%d", config.ServerHost, config.ServerPort)
	clientConn, err := tls.DialWithDialer(
		&dialer,
		"tcp",
		tcpAddr,
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

	// set rlimit nofile value
	tls_proxy_go.SetRLimit(100000)

	config, err := tls_proxy_go.LoadClientConfig()
	if err != nil {
		log.Fatalf("LoadClientConfig: %v", err)
		return
	}

	bindAddr := fmt.Sprintf("%s:%d", config.ListenHost, config.ListenPort)
	listener, err := reuseport.Listen("tcp", bindAddr)
	if err != nil {
		log.Fatalf("reuseport.Listen: %v", err)
		return
	}
	log.Printf("server bind on: %v", bindAddr)

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Accept: %v", err)
			return
		}
		log.Printf("accept connection from: %v", conn.RemoteAddr())

		go clientHandler(conn, config)
	}
}
