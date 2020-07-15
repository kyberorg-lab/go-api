package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"golang.org/x/crypto/scrypt"
)

const (
	SaltLen            = 32
	KeyLen             = 32
	NumberOfIterations = 1048576
	RelativeMemoryCost = 8
	RelativeCPUCost    = 1
)

// Encrypts data (aka passwords) with key (shared secret) returns hex-encoded encrypted text
func Encrypt(key, data []byte) ([]byte, error) {
	key, salt, err := deriveKey(key, nil)
	if err != nil {
		return nil, err
	}

	blockCipher, err := aes.NewCipher(key)
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

	cipherText := gcm.Seal(nonce, nonce, data, nil)
	cipherText = append(cipherText, salt...)

	return cipherText, nil
}

// Encrypts data (aka passwords) with key (shared secret) returns strings
func EncryptString(key, data string) (string, error) {
	bKey := []byte(key)
	bData := []byte(data)

	encryptedText, err := Encrypt(bKey, bData)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(encryptedText), nil
}

// Decrypts data (passwords etc) encrypted with Encrypt function using the same key (shared secret)
func Decrypt(key, data []byte) ([]byte, error) {
	salt, data := data[len(data)-32:], data[:len(data)-32]

	key, _, err := deriveKey(key, salt)
	if err != nil {
		return nil, err
	}

	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	nonce, cipherText := data[:gcm.NonceSize()], data[gcm.NonceSize():]

	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, err
	}
	return plainText, nil
}

// Decrypts data (passwords etc) encrypted with Encrypt function using the same key (shared secret). Returns decrypted string
func DecryptString(key, data string) (string, error) {
	bKey := []byte(key)
	bData := []byte(data)

	plainHexText, err := Decrypt(bKey, bData)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(plainHexText), nil
}

func deriveKey(password, salt []byte) ([]byte, []byte, error) {
	if salt == nil {
		salt = make([]byte, SaltLen)
		if _, err := rand.Read(salt); err != nil {
			return nil, nil, err
		}
	}
	key, err := scrypt.Key(password, salt, NumberOfIterations, RelativeMemoryCost, RelativeCPUCost, KeyLen)
	if err != nil {
		return nil, nil, err
	}
	return key, salt, nil
}
