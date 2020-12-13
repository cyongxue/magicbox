package xhcrypt

import (
	"fmt"
	"log"
	"testing"
)

func TestGenerateKey(t *testing.T) {
	key, err := GenerateKey([]byte("1111"))
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println(string(key))
}

func TestAesCBCEncrypt(t *testing.T) {
	key, err := GenerateKey([]byte("1111"))
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println(string(key))

	aes, err := AesCBCEncrypt([]byte("LftuD3eBuJhRtEkAajk="), []byte("uCd85n9kEOmQf11s+9SShdxzMSDnGqt7ojZPGo0w3nY="))
	if err != nil {
		log.Println(err.Error())
		return
	}
	fmt.Println("===AesCBCEncrypt")
	fmt.Println("LftuD3eBuJhRtEkAajk=")
	fmt.Println(aes)

	aes = "Z1RcJU1cJk1mQlBia2VqRea2B268fG1JeF4LPeLzZXmxnrg0Hox3Yqbc8oGkCc0N"
	plain, err := AesCBCDecrypt(aes, []byte("uCd85n9kEOmQf11s+9SShdxzMSDnGqt7ojZPGo0w3nY="))
	if err != nil {
		log.Println(err.Error())
		return
	}
	fmt.Println("***AesCBCDecrypt")
	fmt.Println(string(plain))
	return
}

//func TestAesCBCEncrypt2(t *testing.T) {
//	aes, err := AesCBCEncrypt2([]byte("q_1dY=Khec2nMNxVq_1dY=Khec2nMNxVq_"))
//	if err != nil {
//		log.Println(err.Error())
//		return
//	}
//
//	log.Println(base64.StdEncoding.EncodeToString(aes))
//	return
//}
