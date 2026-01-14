# Memevault 🔒

**Secure, Portable, and Developer-Friendly Secret Management (with Memes).**

Memevault allows you to store your project secrets (API keys, database URLs, etc.) encrypted *inside* your repository, always hidden within a meme image (`secrets.jpg`). No more sharing `.env` files over Slack!

## Features

- **Store Secrets in Git**: Encrypted vault committed alongside code.
- **Steganography**: Hide secrets inside a meme or image (`secrets.jpg`).
- **Cross-Platform**: Built-in `printenv` polyfill for Windows/Linux consistency.
- **Multi-User**: Easily `grant` access to teammates.
- **Key Rotation**: Securely rotate your identity if compromised.
- **Secret Scanning**: `memevault scan` finds missing secrets in your code.
- **Zero Dependencies**: Single binary.

## Installation

### Option A: Automatic Script (Recommended)
**Linux/macOS**:
```bash
curl -fsSL https://github.com/thoughtlesslabs/memevault/releases/latest/download/install.sh | bash
```

**Windows (PowerShell)**:
```powershell
iwr -useb https://github.com/thoughtlesslabs/memevault/releases/latest/download/install.ps1 | iex
```

### Option B: Build from Source
```bash
git clone https://github.com/thoughtlesslabs/memevault.git
cd memevault
go install
```

## Quick Start

### 1. Initialize
Run this in your project root. It generates your "identity" key (if you don't have one) and creates the vault.
```bash
memevault init
# Defaults to fetching a random meme!
# Or use your own image: memevault init --image ./cool-background.png
```

### 2. Set Secrets
```bash
memevault set DB_PASSWORD "s3cr3t_p@ssw0rd"
memevault set API_KEY "12345-abcde"
```

### 3. View Secrets
You can inspect what's inside the vault:
```bash
# List all secrets
memevault get

# Get a specific secret
memevault get API_KEY
```
*Note: This respects access control. You must be an authorized user to decrypt and view secrets.*

### 4. Run Your App
Memevault injects the secrets into the environment of the command you run.
```bash
# Node.js
memevault run -- node server.js

# Go
memevault run -- go run main.go

# Python
memevault run -- python app.py
```

### Advanced Usage
**Multiple Vaults**: You can specify a different vault file (image) using the `--vault` flag with any command.
```bash
memevault get --vault ./production_secrets.jpg
memevault run --vault ./staging.jpg -- node app.js
```

## Team Workflow

## Managing Access (Multi-User)

Memevault allows you to share secrets with your team securely using **Named Access Keys**.

### 1. View Access List
To see who has access to the vault:
```bash
memevault access list
```

### 2. Grant Access
To add a team member, ask for their public key (`memevault keys show`) and add them:
```bash
# Usage: memevault grant [NAME] [PUBLIC_KEY]
memevault grant bob age1...
```
This will re-encrypt the vault so that both you and Bob can access it.

### 3. Revoke Access
To remove a team member:
```bash
memevault access remove bob
```
This immediately re-encrypts the vault with the remaining keys, locking Bob out.

### Rotating Keys
If your machine is compromised:
```bash
memevault keys rotate
```
This generates a new keypair, re-encrypts the vault (locking out the old key), and backs up the old key.

## Secret Scanning
Check if you've used any variables in your code that aren't in the vault:
```bash
memevault scan
```

## Security Model
**Offline Attack Warning**: If an attacker gets a copy of your `secrets.jpg` AND your private key file, they can decrypt that specific version of the file forever. Key rotation only protects future versions and prevents the compromised key from receiving new updates.
