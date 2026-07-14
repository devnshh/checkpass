<div align="center">

# checkpass

[![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev)

A simple, local command-line password manager with AES-GCM encryption.

Passwords are stored in an encrypted vault file on your machine. Nothing leaves your system.

</div>

---

## Requirements

- Go 1.18 or later

## Setup

Set the master password as an environment variable before using the tool.

```bash
export checkpass_key="your_master_password"
```

Add this to your shell profile to make it persistent.

## Build

```bash
go build -o checkpass .
```

## Usage

### Add a password

```bash
./checkpass add <service> <password>
```

Example:

```bash
./checkpass add github mypassword123
```

### Get a password

```bash
./checkpass get <service>
```

Example:

```bash
./checkpass get github
```

## Notes

- The vault is stored as `vault.json` in the current directory.
- A new vault is created automatically on the first `add`.
- The master password is never stored anywhere, keep it safe.
- Made while learning Go, will add more features as I progress.
