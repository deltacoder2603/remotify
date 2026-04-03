package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/grandcat/zeroconf"
	"golang.org/x/term"
)

const AUTH_TOKEN = "secret123"

type Agent struct {
	ID   string
	Addr string
}

func main() {
	clear()

	fmt.Println("══════════════════════════════")
	fmt.Println("       REMOTIFY SSH ⚡")
	fmt.Println("══════════════════════════════\n")

	agents := discoverAgents()

	if len(agents) == 0 {
		fmt.Println("No agents found")
		return
	}

	fmt.Println("Available Agents:")
	for i, a := range agents {
		fmt.Printf("[%d] %s (%s)\n", i, a.ID[:8], a.Addr)
	}

	fmt.Print("\nSelect agent: ")
	var choice int
	fmt.Scanln(&choice)

	if choice >= len(agents) {
		fmt.Println("Invalid choice")
		return
	}

	connect(agents[choice].Addr)
}

func discoverAgents() []Agent {
	resolver, _ := zeroconf.NewResolver(nil)
	entries := make(chan *zeroconf.ServiceEntry)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var agents []Agent

	go func() {
		for e := range entries {
			if e.Instance != "remotify" {
				continue
			}

			if len(e.AddrIPv4) == 0 {
				continue
			}

			id := "unknown"
			if len(e.Text) > 0 {
				id = e.Text[0]
			}

			addr := e.AddrIPv4[0].String() + ":2222"

			agents = append(agents, Agent{
				ID:   id,
				Addr: addr,
			})
		}
	}()

	resolver.Browse(ctx, "_tcp", "local.", entries)
	<-ctx.Done()

	return agents
}

func connect(addr string) {
	fmt.Println("\nConnecting to:", addr)

	conn, err := tls.Dial("tcp", addr, &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		fmt.Println("Connection error:", err)
		return
	}
	defer conn.Close()

	fmt.Fprintf(conn, AUTH_TOKEN+"\n")

	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)

	if string(buf[:n]) != "AUTH OK\n" {
		fmt.Println("Auth failed")
		return
	}

	fmt.Println("Connected\n")

	oldState, _ := term.MakeRaw(int(os.Stdin.Fd()))
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	go io.Copy(conn, os.Stdin)
	io.Copy(os.Stdout, conn)
}

func clear() {
	fmt.Print("\033[H\033[2J")
}
