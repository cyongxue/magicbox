package xhcrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"runtime"
)

// 配置加密对象
var ConfigAes *AesConfig

func Init(latestVersionFlag string, key []byte, oldKeys map[string][]byte) error {
	if len(latestVersionFlag) != versionSize {
		return fmt.Errorf("latest version flag is invalid, %s", latestVersionFlag)
	}
	for k, _ := range oldKeys {
		if len(k) != versionSize {
			return fmt.Errorf("the version flag is invalid, %s", k)
		}
	}
	ConfigAes = &AesConfig{
		Key:                key,
		OldKeys:            oldKeys,
		CurrentVersionFlag: latestVersionFlag,
	}
	return nil
}

// AesConfig 配置类加密处理
type AesConfig struct {
	CurrentVersionFlag string
	Key                []byte
	OldKeys            map[string][]byte
}

// version 表明加解密的版本区分，采用密文的前6个字节代表版本号
const (
	versionSize = 6

	// 第一次没有版本标签，所以定义一个ZeroVersionFlag
	ZeroVersionFlag = "ZeroVe"
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

// AesConfig.Encrypt 加密
// plainText 明文
// key 32字节秘钥
func (a *AesConfig) Encrypt(plainText []byte) (string, error) {
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("AesConfig.Encrypt panic||panicMsg=%s\n", string(buf))
		}
	}()

	cipherText, err := a.cBCEncrypt(plainText, a.Key)
	if err != nil {
		return "", err
	}
	return a.CurrentVersionFlag + cipherText, nil
}

// AesConfig.AesCBCEncrypt aes cbc方式加密
func (a *AesConfig) cBCEncrypt(plaintext []byte, originKey []byte) (string, error) {
	if len(originKey) < 32 {
		return "", fmt.Errorf("key length must then 32 byte: %d", len(originKey))
	}

	// 生成随机salt，不能小于128bit，这里取16字节
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

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

// AesConfig.Decrypt 解密，支持可选指定老的秘钥，提供版本升级时进行兼容老版本解密处理
// cipherTextStr 密文
// key 最新版解密秘钥
// v   可选参考，兼容版本以及对应的解密秘钥
func (a *AesConfig) Decrypt(cipherTextStr string) ([]byte, error) {
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("AesConfig.Decrypt panic||panicMsg=%s\n", string(buf))
		}
	}()

	versionFlag := cipherTextStr[:versionSize]
	cipherTextNoVersion := cipherTextStr[versionSize:]
	if versionFlag == a.CurrentVersionFlag {
		return a.cBCDecrypt(cipherTextNoVersion, a.Key)
	}

	return nil, fmt.Errorf("no found decrypt function")
}

// AesConfig.CBCDecrypt 解密处理
func (a *AesConfig) cBCDecrypt(cipherTextStr string, originKey []byte) ([]byte, error) {
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
		return nil, fmt.Errorf("aes decrypt failed")
	}
	plaintext := PKCS5UnPadding(cipherText)

	return plaintext, nil
}
