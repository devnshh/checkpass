package vault

import (
	"crypto/aes"
	"crypto/sha256"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

type Vault struct {
	isLocked       bool
	masterPassword string
	credentials    map[string]string
}

func NewVault(masterPassword string) (*Vault, error) {
	if masterPassword == "" {
		return nil, errors.New("vault password cannot be empty")
	}
	return &Vault{
		masterPassword: masterPassword,
		isLocked:       true,
		credentials:    make(map[string]string),
	}, nil
}

func (v *Vault) Lock() {
	v.isLocked = true
}

func (v *Vault) Unlock(password string) error {
	if password == v.masterPassword {
		v.isLocked = false
		return nil
	}
	return errors.New("wrong password, access denied")
}
 
func encrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err!= nil {
		return nil, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err!= nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return aesGCM.Seal(nonce, nonce, data, nil), nil

}

func decrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return aesGCM.Open(nil, nonce, ciphertext, nil)
}

func (v *Vault) AddCredential(service, password string) (string, error) {
	if v.isLocked {
		return "", errors.New("cannot add credential, the vault is locked")
	} else if password == "" {
		return "", errors.New("credential cannot be empty")
	}
	v.credentials[service] = password
	return fmt.Sprintf("credential added for %s\n",service), nil
}

func (v *Vault) GetCredential(service string) (string, error) {
	if v.isLocked {
		return "", errors.New("cannot get credential, the vault is locked")
	}
	value, ok := v.credentials[service]
	if !ok {
		return "", errors.New("credential not found")
	}
	return value, nil
}

func deriveKey(password string) ([]byte) {
	hash := sha256.Sum256([]byte(password))
	return hash[:]
}

func (v *Vault) Save() error {
	if v.isLocked {
		return errors.New("cannot save vault, the vault is locked")
	}
	bytes, err := json.Marshal(v.credentials)
	if err != nil {
		return fmt.Errorf("failed to marshal json:%w", err)
	}
	key := deriveKey(v.masterPassword)
	bytes_encrypted, err := encrypt(bytes, key)
	if err != nil {
		return fmt.Errorf("failed to encrypt data:%w",err)
	}

	err = os.WriteFile("vault.json", bytes_encrypted, 0600)
	if err != nil {
		return fmt.Errorf("could not save file:%w", err)
	}
	return nil
}

func (v *Vault) Load() error {
	if v.isLocked {
		return errors.New("cannot load vault, the vault is locked")
	}
	bytes, err := os.ReadFile("vault.json")
	if err != nil {
		return fmt.Errorf("reading vault: %w", err)
	}
	key := deriveKey(v.masterPassword)
	bytes_decrypted,err := decrypt(bytes, key)
	if err != nil {
		return fmt.Errorf("failed to decrypt data:\n%w", err)
	}
	err = json.Unmarshal(bytes_decrypted, &v.credentials)
	if err != nil {
		return fmt.Errorf("could not unmarshal:%w", err)
	}
	return nil
}
