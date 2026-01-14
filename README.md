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

### Option A: Build from Source
```bash
git clone https://github.com/jdiet/memevault.git
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

### 3. Run Your App
Memevault injects the secrets into the environment of the command you run.
```bash
# Node.js
memevault run -- node server.js

# Go
memevault run -- go run main.go

# Python
memevault run -- python app.py
```

## Team Workflow

### Granting Access
When a new team member joins:
1. They install Memevault and run `memevault init`.
2. They send you their public key (run `memevault keys show`).
3. You authorize them:
   ```bash
   memevault grant age1...<THEIR_KEY>...
   ```
4. Commit the updated `secrets.jpg`.

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
