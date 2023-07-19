package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

func encryptData(data []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 生成随机盐值
	salt := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}

	// 创建一个新的GCM模式的AES加密器
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 使用随机盐值对数据进行加密
	encryptedData := aesgcm.Seal(nil, salt, data, nil)

	// 将加密后的数据和盐值合并，并转换为Base64字符串
	encodedData := base64.StdEncoding.EncodeToString(append(salt, encryptedData...))

	return encodedData, nil
}

func decryptData(encodedData string, key []byte) ([]byte, error) {
	// 将Base64字符串转换为字节数组
	decodedData, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 提取随机盐值
	salt := decodedData[:12]
	// 提取加密后的数据
	encryptedData := decodedData[12:]

	// 创建一个新的GCM模式的AES解密器
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 使用提取的盐值对数据进行解密
	decryptedData, err := aesgcm.Open(nil, salt, encryptedData, nil)
	if err != nil {
		return nil, err
	}

	return decryptedData, nil
}

func main() {
	data := []byte("test")
	key := []byte("12345678901234567890123456789012")

	encryptedData, err := encryptData(data, key)
	if err != nil {
		fmt.Println("加密失败：", err)
		return
	}

	fmt.Println("加密后的数据：", encryptedData)

	data, err = decryptData(encryptedData, key)
	if err != nil {
		fmt.Println("解密失败：", err)
		return
	}

	fmt.Println(string(data))
}
