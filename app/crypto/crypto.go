package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"golang.org/x/crypto/scrypt"
)

const (
	SaltLen            = 32
	KeyLen             = 32
	NumberOfIterations = 1048576
	RelativeMemoryCost = 8
	RelativeCPUCost    = 1
)

// Encrypts data (aka passwords) with key (shared secret) returns encrypted text
func Encrypt(plainText, sharedKey []byte) ([]byte, error) {
	sharedKey, err := deriveKey(sharedKey)
	if err != nil {
		return nil, err
	}

	blockCipher, err := aes.NewCipher(sharedKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}

	encrypted := gcm.Seal(nil, nonce, plainText, nil)
	encrypted = append(nonce, encrypted...)

	return encrypted, nil
}

// Encrypts data (aka passwords) with key (shared secret) returns base64-encoded string
func EncryptString(plainTextData, sharedKey string) (string, error) {
	bPlainData := []byte(plainTextData)
	bSharedKey := []byte(sharedKey)

	encrypted, err := Encrypt(bPlainData, bSharedKey)
	if err != nil {
		return "", err
	}

	encryptedText := base64.StdEncoding.EncodeToString(encrypted)
	return encryptedText, nil
}

// Decrypts data (passwords etc) encrypted with Encrypt function using the same key (shared secret)
func Decrypt(encryptedData, sharedKey []byte) ([]byte, error) {
	sharedKey, err := deriveKey(sharedKey)
	if err != nil {
		return nil, err
	}

	blockCipher, err := aes.NewCipher(sharedKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	nonce, encrypted := encryptedData[:gcm.NonceSize()], encryptedData[gcm.NonceSize():]

	plainText, err := gcm.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return nil, err
	}
	return plainText, nil
}

// Decrypts data (passwords etc) encrypted with Encrypt function using the same key (shared secret). Returns decrypted string
func DecryptString(encryptedData, sharedKey string) (string, error) {
	var bEncryptedData []byte

	if isBase64String(encryptedData) {
		bEncryptedData, _ = base64.StdEncoding.DecodeString(encryptedData)
	} else {
		bEncryptedData = []byte(encryptedData)
	}

	bSharedKey := []byte(sharedKey)

	plainText, err := Decrypt(bEncryptedData, bSharedKey)
	if err != nil {
		return "", err
	}
	return string(plainText), nil
}

func deriveKey(password []byte) ([]byte, error) {
	staticSalt := []byte("My Static Salt") //FIXME static salt
	key, err := scrypt.Key(password, staticSalt, NumberOfIterations, RelativeMemoryCost, RelativeCPUCost, KeyLen)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func isBase64String(s string) bool {
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}
