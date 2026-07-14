package main

import (
	"checkpass/vault"
	"errors"
	"fmt"
	"os"
)

func main() {
	masterKey := os.Getenv("checkpass_key")
	if masterKey == "" {
		fmt.Println("Error: Master password is not set")
		return
	}
	if len(os.Args) < 2 {
		fmt.Println("Usage: checkpass [add/get] [service] [password]")
		return
	}

	Vault, err := vault.NewVault(masterKey)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = Vault.Unlock(masterKey)
	if err != nil {
		fmt.Println(err)
		return
	}
	LoadError := Vault.Load()
	newVault := false
	if LoadError != nil {
		if errors.Is(LoadError, os.ErrNotExist) {
			newVault = true
		} else {
			fmt.Printf("error while loading the vault: %v", LoadError)
			return
		}
	}

	switch os.Args[1] {
	case "add":
		if len(os.Args) < 4 {
			fmt.Println("Usage: checkpass [add/get] [service] [password]")
			return
		}
		ok, err := Vault.AddCredential(os.Args[2], os.Args[3])
		if err != nil {
			fmt.Println(err)
			return
		}
		err = Vault.Save()
		if err != nil {
			fmt.Println(err)
			return
		}
		if newVault {
			fmt.Println("new vault created")
		}
		fmt.Printf("status: %s", ok)

	case "get":
		if len(os.Args) != 3 {
			fmt.Println("Usage: checkpass [add/get] [service] [password]")
			return
		}
		password, status := Vault.GetCredential(os.Args[2])
		if status != nil {
			fmt.Println(status)
			return
		}
		fmt.Printf("saved password for %s is: %s\n", os.Args[2], password)

	default:
		fmt.Println("Usage: checkpass [add/get] [service] [password]")
	}
}
