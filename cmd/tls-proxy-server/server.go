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

func serverHandler(conn *tls.Conn) {
	defer conn.Close()

	// read handshake data
	body := make([]byte, 4096)
	n, err := conn.Read(body)
	if err != nil {
		log.Printf("Read: %v", err)
		return
	}

	_, err = tls_proxy_go.ParseHandshakeBody(body[:n])
	if err != nil {
		log.Printf("ParseHandshakeBody: %v", err)
		return
	}

	_, err = conn.Write([]byte{0x05, 0x00})
	if err != nil {
		log.Printf("Write: %v", err)
		return
	}

	// read request data
	n, err = conn.Read(body)
	if err != nil {
		log.Printf("Read: %v", err)
		return
	}

	tuple, err := tls_proxy_go.ParseRequestBody(body[:n])
	if err != nil {
		log.Printf("ParseRequestBody: %v", err)
		return
	}

	if (*tuple)[0].(byte) != 0x01 {
		log.Printf("CMD: 0x%x not support", (*tuple)[0].(byte))
		return
	}

	if (*tuple)[1].(byte) == 0x04 {
		log.Println("ip v6 not supported")
		return
	}

	log.Printf("proxy to: %v", fmt.Sprintf("%s:%d", (*tuple)[2], (*tuple)[3]))

	tcpAddr := fmt.Sprintf("%s:%d", (*tuple)[2], (*tuple)[3])
	clientConn, err := net.DialTimeout(
		"tcp",
		tcpAddr,
		3*time.Second)
	if err != nil {
		log.Printf("net.DialTimeout: %v", err)
		return
	}

	defer clientConn.Close()

	_, err = conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	if err != nil {
		log.Printf("Write: %v", err)
		return
	}

	go func() { _, _ = io.Copy(conn, clientConn) }()
	go func() { _, _ = io.Copy(clientConn, conn) }()

	select {}

}

func main() {
	log.SetOutput(os.Stdout)

	// set rlimit nofile value
	tls_proxy_go.SetRLimit(100000)

	config, err := tls_proxy_go.LoadServerConfig()
	if err != nil {
		log.Fatalf("LoadServerConfig: %v", err)
		return
	}

	cert, err := tls.LoadX509KeyPair(config.ServerCert, config.ServerKey)
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

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}

	bindAddr := fmt.Sprintf("%s:%d", config.ListenHost, config.ListenPort)
	listener, err := reuseport.Listen("tcp", bindAddr)
	if err != nil {
		log.Fatalf("tls.Listen: %v", err)
		return
	}

	tlsListener := tls.NewListener(listener, tlsConfig)
	log.Printf("tls server bind on: %v", bindAddr)

	defer tlsListener.Close()

	for {
		conn, err := tlsListener.Accept()
		if err != nil {
			log.Fatalf("Accept: %v", err)
			return
		}
		log.Printf("accept connection from: %v", conn.RemoteAddr())

		go serverHandler(conn.(*tls.Conn))
	}
}
