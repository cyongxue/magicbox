package xhcrypt

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestAesConfig_Encrypt(t *testing.T) {
	plaintext := "LftuD3eBuJhRtEkAajk="
	key := "uCd85n9kEOmQf11s+9SShdxzMSDnGqt7ojZPGo0w3nY="

	oldKeys := make(map[string][]byte)
	if err := Init("=hehui", []byte(key), oldKeys); err != nil {
		fmt.Println(err.Error())
		return
	}

	cipherText, err := ConfigAes.Encrypt([]byte(plaintext))
	if err != nil {
		fmt.Println("Encrypt error: " + err.Error())
		return
	}
	fmt.Println("cipherText = " + cipherText)

	newPlain, err := ConfigAes.Decrypt(cipherText)
	if err != nil {
		fmt.Println("Decrypt error: " + err.Error())
		return
	}
	fmt.Println("newPlain  = " + string(newPlain))
	fmt.Println("plaintext = " + plaintext)
	return
}

func TestAesConfig_Decrypt(t *testing.T) {
	key := "uCd85n9kEOmQf11s+9SShdxzMSDnGqt7ojZPGo0w3nY="

	oldKeys := make(map[string][]byte)
	if err := Init("=DhaV=", []byte(key), oldKeys); err != nil {
		fmt.Println(err.Error())
		return
	}

	cipherTextNew := "=DhaV=BiIHehUlKWJqQEUKTg4UMGePRBju+iTlznZLf1jtsvw1BSzOsVsDuX3iM9nc+Rac"
	decodePlain, err := ConfigAes.Decrypt(cipherTextNew)
	if err != nil {
		fmt.Println("Decrypt error: " + err.Error())
		return
	}
	fmt.Println("newPlain  = " + string(decodePlain))
	plaintext := "LftuD3eBuJhRtEkAajk="
	fmt.Println("plaintext = " + plaintext)

	fmt.Println("========================================")
	cipherTextOld := "pM5fTs4HG/AoOy90GCkHZQAAAAc="
	decodePlain, err = ConfigAes.Decrypt(cipherTextOld)
	if err != nil {
		fmt.Println("Decrypt error: " + err.Error())
		return
	}
	fmt.Println("newPlain  = " + string(decodePlain))

	return
}

func BenchmarkAesConfig_Encrypt(b *testing.B) {
	plaintext := "LftuD3eBuJhRtEkAajk="
	key := "uCd85n9kEOmQf11s+9SShdxzMSDnGqt7ojZPGo0w3nY="

	oldKeys := make(map[string][]byte)
	if err := Init("=DhaV=", []byte(key), oldKeys); err != nil {
		fmt.Println(err.Error())
		return
	}

	for i := 0; i < b.N; i++ {
		_, err := ConfigAes.Encrypt([]byte(plaintext))
		if err != nil {
			return
		}
	}
}

func TestPBKDF2(t *testing.T) {
	data, err := base64.StdEncoding.DecodeString("b19mUDdXVighL2NSfhtYOgl0Vv/mwRaY5FI7ApP6MLQ=")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(data)
	//f, err := os.Create("Z:\\Goland\\src\\gitee.com\\yongxue\\magicbox\\main\\logs\\test.log")
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return
	//}
	//f.Write(data)
	//f.Close()
	return
}
