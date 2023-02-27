package secret

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	encrypt "github.com/radoslavboychev/gophercises-secret/cipher"
)

// File creates a new instance of a Vault object
func File(encodingKey, filepath string) *Vault {
	return &Vault{
		encodingKey: encodingKey,
		filepath:    filepath,
	}
}

// Vault represents a data structure
type Vault struct {
	encodingKey string
	filepath    string
	mutex       sync.Mutex
	keyValues   map[string]string
}

// load opens
func (v *Vault) load() error {
	f, err := os.Open(v.filepath)
	if err != nil {
		v.keyValues = make(map[string]string)
		return nil
	}
	defer f.Close()
	r, err := encrypt.DecryptReader(v.encodingKey, f)
	if err != nil {
		return err
	}
	return v.readKeyValues(r)

}

func (v *Vault) readKeyValues(r io.Reader) error {
	dec := json.NewDecoder(r)
	return dec.Decode(&v.keyValues)
}

// loadKeyValues opens the file with the values, decrypts them and decodes
func (v *Vault) loadKeyValues() error {
	f, err := os.Open(v.filepath)
	if err != nil {
		v.keyValues = make(map[string]string)
		return nil
	}
	defer f.Close()
	var sb strings.Builder
	_, err = io.Copy(&sb, f)
	if err != nil {
		return err
	}
	decryptedJSON, err := encrypt.Decrypt(v.encodingKey, sb.String())
	if err != nil {
		return err
	}
	r := strings.NewReader(decryptedJSON)
	dec := json.NewDecoder(r)
	err = dec.Decode(&v.keyValues)
	if err != nil {
		return err
	}
	return nil
}

// saveKeyValues
func (v *Vault) saveKeyValues() error {
	var sb strings.Builder
	enc := json.NewEncoder(&sb)
	err := enc.Encode(v.keyValues)
	if err != nil {
		return err
	}
	encryptedJSON, err := encrypt.Encrypt(v.encodingKey, sb.String())
	if err != nil {
		return err
	}
	f, err := os.OpenFile(v.filepath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = fmt.Fprint(f, encryptedJSON)
	if err != nil {
		return err
	}
	return nil
}

// Get decrypts and returns a key from the vault
func (v *Vault) Get(key string) (string, error) {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	err := v.loadKeyValues()
	if err != nil {
		return "", err
	}
	value, ok := v.keyValues[key]
	if !ok {
		return "", errors.New("secret: no value for that key")
	}
	return value, nil
}

// Set sets the value of a key in the Vault, returns error if any
func (v *Vault) Set(key, value string) error {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	err := v.loadKeyValues()
	if err != nil {
		return err
	}
	v.keyValues[key] = value
	err = v.saveKeyValues()
	return err

}
