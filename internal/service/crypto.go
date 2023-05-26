package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"strconv"
	"strings"
)

func (s *Service) encryptProductAES(chatID int64, countBought int64) (string, error) {
	log.WithFields(log.Fields{
		"chatID":      chatID,
		"countBought": countBought,
	}).Info("encryptProductAES")

	message := fmt.Sprintf("%d;%d", chatID, countBought)

	//Create byte array from the input string
	plainText := []byte(message)

	//Create a new AES cipher using the key
	block, err := aes.NewCipher([]byte(s.payment.KeyPayment))

	//IF NewCipher failed, exit:
	if err != nil {
		return "", err
	}

	//Make the cipher text a byte array of size BlockSize + the length of the message
	cipherText := make([]byte, aes.BlockSize+len(plainText))

	//iv is the ciphertext up to the blocksize (16)
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	//Encrypt the data:
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	//Return string encoded in base64
	return base64.RawURLEncoding.EncodeToString(cipherText), nil
}

func (s *Service) decryptProductAES(uEnc string) (int64, int64, error) {
	//Remove base64 encoding:
	cipherText, err := base64.RawURLEncoding.DecodeString(uEnc)

	//IF DecodeString failed, exit:
	if err != nil {
		return 0, 0, err
	}

	//Create a new AES cipher with the key and encrypted message
	block, err := aes.NewCipher([]byte(s.payment.KeyPayment))

	//IF NewCipher failed, exit:
	if err != nil {
		return 0, 0, err
	}

	//IF the length of the cipherText is less than 16 Bytes:
	if len(cipherText) < aes.BlockSize {
		return 0, 0, blockSizeError
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	//Decrypt the message
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	stringRes := string(cipherText)

	arrParams := strings.Split(stringRes, ";")

	if len(arrParams) != 2 {
		return 0, 0, decodingError
	}

	chatID, err := strconv.ParseInt(arrParams[0], 10, 64)
	if err != nil {
		return 0, 0, decodingError
	}

	countBought, err := strconv.ParseInt(arrParams[1], 10, 64)
	if err != nil {
		return 0, 0, decodingError
	}

	log.WithFields(log.Fields{
		"chatID":      chatID,
		"countBought": countBought,
	}).Info("decryptProductAES")

	return chatID, countBought, nil
}
