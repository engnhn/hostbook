# hostbook

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go](https://img.shields.io/badge/language-Go-00ADD8.svg)
![Status](https://img.shields.io/badge/status-stable-green.svg)
![Storage](https://img.shields.io/badge/storage-local--json-orange.svg)

> "SSH keys are easy. Remembering 50 IP addresses is not."

**hostbook** is a "boring solution" to a common problem: Managing a growing list of SSH connections without losing your sanity or your security. It acts as a local source of truth for your servers, generating standard OpenSSH configurations while providing a sleek interactive CLI.

## The Problem

SSH is the backbone of infrastructure, but managing connection details is often stuck in the dark ages. You have a text file (`~/.ssh/config`) that grows indefinitely, or worse, a sticky note with IP addresses. Aliases are forgotten, users are mixed up, and finding the right server becomes a memory game.

We don't want "smart" cloud tools that sync our private infrastructure details to someone else's server. We want a tool that respects our local environment but organizes the chaos.

## The Solution

Hostbook is built on honest constraints:

1.  **Zero Magic, Just Config**: It doesn't hijack your connection. It generates a standard `~/.ssh/config` compliant file (`~/.hostbook/ssh_config`) and calls the native `ssh` binary.
2.  **Local Sovereignty**: Your server data lives effectively in `~/.hostbook/hosts.json` on your machine. No cloud sync, no tracking, no "account required".
3.  **Interactive by Default**: Don't remember the server name? Just type `hostbook connect` and pick from the list.

## Usage

### 1. Simple Management

Forget editing text files. Manage your hosts with simple commands.

```bash
# Add a new host (Interactive or Flags)
hostbook add

# List what you have
hostbook list
```

### 2. Connect in Flow

You can connect directly if you know the name, or use the interactive picker if you don't.

```bash
# I know where I'm going
hostbook connect prod-db

# I need a reminder (Opens Interactive Menu)
hostbook connect
```

### 3. Edit & Maintain

Servers change. Your config should too.

```bash
# Update existing host details
hostbook edit prod-db

# Remove a decommissioned server
hostbook delete legacy-app
```

### 4. Advanced Features

#### Security & Auto-Connect
Hostbook securely stores passwords in your OS keyring (Keychain/SecretService).
-   **Add**: Prompts to save password securely.
-   **Connect**: Auto-fills password if `sshpass` is installed.

#### Shell Completion
Get dynamic tab completion for your host names.
```bash
# Add to your .bashrc or .zshrc
source <(hostbook completion bash)
# or
source <(hostbook completion zsh)
```

#### Backup & Export
Don't get locked in. Export your data anytime.
```bash
hostbook export --format yaml > backup.yaml
hostbook export --format json > backup.json
```

#### Port Forwarding
Quickly forward ports without memorizing flags.
```bash
# Forward local 8080 to remote 80
hostbook forward my-server 8080:80
```

#### Health Check
See which servers are online.
```bash
hostbook ping
```

## Installation

```bash
# Install directly (Go required)
go install github.com/engnhn/hostbook@latest

# Or build from source
git clone https://github.com/engnhn/hostbook.git
cd hostbook
go build -o hostbook main.go
```

## Internal Architecture

We favor transparency over abstraction.

-   **Storage**: `~/.hostbook/hosts.json` - Single source of truth.
-   **Config Generation**: `~/.hostbook/ssh_config` - Read-only, generated file passed to SSH.
-   **Execution**: Hostbook replaces itself with the `ssh` process (`syscall.Exec`), ensuring your terminal session handles signals and I/O natively.

## Development

I write code, break it, and then fix it. Here's how to do the same:

```bash
# Install dependencies
go mod download

# Run functionality
go run main.go list

# Build the binary
go build -o hostbook main.go
```

## License

MIT
