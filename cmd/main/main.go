package main

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
)

func main() {

	someText := "booobooo"
	hash, err := hashTextTo32Bytes(someText)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(hash)

}

func hashTextTo32Bytes(hashThis string) (hashed string, err error) {
	if len(hashThis) == 0 {
		return "", errors.New("no input supplied")
	}

	hasher := sha256.New()
	hasher.Write([]byte(hashThis))

	stringToSHA256 := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	return stringToSHA256[:32], nil
}

