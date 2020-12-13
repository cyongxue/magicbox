package xhcrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
)

//AES256-CBC(data,key[],iv[])
//Key=PDKDF2(MasterKey,随机数)
//说明：
// 		data: 加密数据
//		采用PKCS5Padding的方式进行padding处理
//	KDF算法中的password是MasterKey
//	KDF算法中的salt是“随机数”，随机数要求不小于128bits
//	iv可以等于随机数s
//	MasterKey为用安全随机函数生成的256bits及以上固定字符串

// AesCBCEncrypt AES CBC加密
func AesCBCEncrypt(plaintext []byte, originKey []byte) (string, error) {
	if len(originKey) < 32 {
		return "", fmt.Errorf("key length must then 32 byte: %d", len(originKey))
	}
	log.Println(originKey)

	// 生成随机salt，不能小于128bit，这里取16字节
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	log.Println(iv)

	// 利用PBKDF2算法计算key
	key := PBKDF2(originKey, iv)

	// 创建一个block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 明文和盐进行block size处理
	plaintext = PKCS5Padding(plaintext, aes.BlockSize)

	// 加密
	cipherText := make([]byte, aes.BlockSize+len(plaintext))
	blockMode := cipher.NewCBCEncrypter(block, iv)
	blockMode.CryptBlocks(cipherText[aes.BlockSize:], plaintext)
	copy(cipherText[:aes.BlockSize], iv)

	cipherTextStr := base64.StdEncoding.EncodeToString(cipherText)
	return cipherTextStr, nil
}

// AesCBCDecrypt AES CBC解密
func AesCBCDecrypt(cipherTextStr string, originKey []byte) ([]byte, error) {
	if len(originKey) < 32 {
		return nil, fmt.Errorf("key length must then 32 byte: %d", len(originKey))
	}
	// 取出salt
	cipherText, err := base64.StdEncoding.DecodeString(cipherTextStr)
	if err != nil {
		return nil, fmt.Errorf("cipherTextStr DecodeString error: %s", err.Error())
	}
	iv := make([]byte, aes.BlockSize)
	copy(iv, cipherText[:aes.BlockSize])

	// 利用PBKDF2算法计算key
	key := PBKDF2(originKey, iv)

	// 创建一个block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	cipherText = cipherText[aes.BlockSize:]
	blockMode := cipher.NewCBCDecrypter(block, iv)
	blockMode.CryptBlocks(cipherText, cipherText)

	if int(cipherText[len(cipherText)-1]) > len(cipherText) {
		return nil, errors.New("aes decrypt failed")
	}
	plaintext := PKCS5UnPadding(cipherText)

	return plaintext, nil
}
