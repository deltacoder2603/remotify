# ⚡ Remotify

A lightweight **SSH-like remote terminal system** built in Go, enabling real-time command execution across devices on the same network using **mDNS discovery** and **PTY-based shell streaming**.

---

## 🚀 Features

* 🔍 **mDNS Auto Discovery**
  Automatically detects available agents on the same network (no IP needed)

* 💻 **Real-Time Terminal (PTY)**
  Full interactive shell experience (like SSH)

* 👥 **Multiple Controllers Support**
  Multiple clients can connect to the same agent simultaneously

* 🧠 **Cross-Platform Support**

  * macOS / Linux → Full PTY terminal
  * Windows → PowerShell fallback

* 🔐 **Encrypted Communication (TLS)**
  Secure communication without manual certificate setup

* ⚡ **Minimal & Fast**
  No external dependencies, simple architecture

---

## 🏗️ Architecture

```
Controller (CLI)
        ↓
   mDNS Discovery
        ↓
     Agent (Server)
        ↓
   PTY Shell (bash / powershell)
```

---

## 📁 Project Structure

```
REMOTIFY/
├── agent/
│   └── main.go
├── controller/
│   └── main.go
└── builds/
```

---

## 🛠️ Installation & Setup

### 1. Clone the repository

```bash
git clone <your-repo-url>
cd remotify
```

---

### 2. Run Agent (on target machine)

```bash
cd agent
go run .
```

Output:

```
Agent running on :2222
mDNS: remotify.local
```

---

### 3. Run Controller

```bash
cd controller
go run .
```

---

## 💻 Usage

Once connected:

```bash
ls
pwd
cd ..
```

Output behaves like a real terminal:

```
bash-3.2$ ls
main.go go.mod
```

---

## 📦 Build Binaries

```bash
# macOS
GOOS=darwin GOARCH=amd64 go build -o builds/agent/agent-mac ./agent
GOOS=darwin GOARCH=amd64 go build -o builds/controller/controller-mac ./controller

# Linux
GOOS=linux GOARCH=amd64 go build -o builds/agent/agent-linux ./agent
GOOS=linux GOARCH=amd64 go build -o builds/controller/controller-linux ./controller

# Windows
GOOS=windows GOARCH=amd64 go build -o builds/agent/agent.exe ./agent
GOOS=windows GOARCH=amd64 go build -o builds/controller/controller.exe ./controller
```

---

## ⚠️ Limitations

* Works only on **same LAN / WiFi**
* mDNS may not work on:

  * public networks
  * college / office WiFi
* Windows does not support full PTY (limited terminal features)

---

## 🔮 Future Improvements

* 🌍 Internet-based connectivity (relay server)
* 🔐 SSH-style key authentication
* 📜 Command history & autocomplete
* 🎨 Enhanced CLI UI
* 🌐 Web-based terminal interface

---

## 🧠 Learnings

This project demonstrates:

* TCP networking in Go
* PTY (pseudo-terminal) handling
* Real-time bidirectional streaming
* mDNS (service discovery)
* Cross-platform system design

---

## 📌 Resume Description

> Built a cross-platform SSH-like remote terminal system in Go using TCP sockets, PTY-based shell execution, and mDNS service discovery. Implemented real-time bidirectional communication with support for multiple clients and encrypted connections.

---

## 📜 License

MIT License

---

## 👨‍💻 Author

Srikant Pandey
