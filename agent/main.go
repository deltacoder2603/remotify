package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"math/big"
	"net"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/creack/pty"
	"github.com/google/uuid"
	"github.com/grandcat/zeroconf"
)

const AUTH_TOKEN = "secret123"

type Message struct {
	Type string
	Data string
}

// 🔐 TLS (auto)
func generateTLSConfig() *tls.Config {
	priv, _ := rsa.GenerateKey(rand.Reader, 2048)

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}

	derBytes, _ := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)

	cert := tls.Certificate{
		Certificate: [][]byte{derBytes},
		PrivateKey:  priv,
	}

	return &tls.Config{Certificates: []tls.Certificate{cert}}
}

func main() {
	config := generateTLSConfig()

	ln, err := tls.Listen("tcp", "0.0.0.0:2222", config)
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	id := uuid.New().String()

	fmt.Println("Agent running:", id)

	mdns, _ := zeroconf.Register("remotify", "_tcp", "local.", 2222, []string{id}, nil)
	defer mdns.Shutdown()

	for {
		conn, _ := ln.Accept()
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)

	if strings.TrimSpace(string(buf[:n])) != AUTH_TOKEN {
		conn.Write([]byte("AUTH FAILED\n"))
		return
	}

	conn.Write([]byte("AUTH OK\n"))

	switch runtime.GOOS {

	case "windows":
		cmd := exec.Command("powershell")

		stdin, _ := cmd.StdinPipe()
		stdout, _ := cmd.StdoutPipe()
		cmd.Stderr = cmd.Stdout

		cmd.Start()

		go io.Copy(stdin, conn)
		io.Copy(conn, stdout)

	case "linux", "darwin":
		cmd := exec.Command("bash")

		ptmx, _ := pty.Start(cmd)
		defer ptmx.Close()

		go io.Copy(ptmx, conn)
		io.Copy(conn, ptmx)
	}
}
