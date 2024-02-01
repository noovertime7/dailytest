package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log"
)

var PrivateKey = `LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1JSUV2UUlCQURBTkJna3Foa2lHOXcwQkFRRUZBQVNDQktjd2dnU2pBZ0VBQW9JQkFRRHZ6ZlBXNFlpZnZ5ZVoKR2VWK2cvbEd5ZHVaMkVINjBPakxFQW9BOVZLdFo2Y3gyQWtXK2ZaRnhxSkFGNGZZRnNvMUJ3Ly9LQmg0NFFMWgptRnY0UmtjRERKUFJtamZpOWxBQ0JWWTZMbk5VTTNERzdoVTk5ejdibjVnMEN2NW5MKzNHb21WWVhrUnMvS3FOCkprWkJBai9IQmxLK0JidnE1elBhK1BZc3FqSjEwUktsdVpYUkJCaFhHOURXQWVnRTRNZGU1enAwUEJCa3oxdzkKYUg4MytjZjBqamxQTFJsWUpNNlAyb0tULzNZS0V0b2orRmpqbm5iSXNNeXZPeW9TcXN2UXF3SkFxcWRQSUo1UwpDcThjWEVZclRObTNmbWh2UUx1blVWZ25nY2NkOU1kWHBQUnJ4VFVKaTVyazBrMkxtWS80MXZvb2pLNVlJZXRBClo0cDQzdmtMQWdNQkFBRUNnZ0VBR1lEejMvU1lqVFROUjFFandUTFQvaDFWeDZUaVQ0U01YWnhWRkFrRFJBSDAKSEM3M3VJUFpGMDZxenRveHpsK09IZG1rYitTWnhiRllsai9IK0QyeEs3Zll1TUFJRlpGUXJRWllSMVBOWERVawpWMlBXeW9KVEl4UjBJWVRSemJPaFBERGxTbUtzTUZNbTZHQWJ0S3BDa2k0ditwbXRoS3dXTGNHd25NdDFGbVp6Ck9tVXNWUkg3a01rL1doVzNvb1B3Ny9ZVGhNSW1zZkJoNm1QT3VpRm9IRHRHNnNsWmJ5YTRhdjcrRlQ5eldOWHIKUzNkTmtVZFR6eENtcC9TV1paSWV6UlR3R0wvaXRvdTVtN2owSjVPbkJDYVNxL0lWR1U1UXoydm1ydzl4LzQyagpIZXVzWXhKZGlSQUJuR3N3Tm1xL0x4VXYyNVBuckptUDQ3MTY3VHhYc1FLQmdRRDYwZVBqTDJNVGx0REF2U3crClcvMkJDcFZTYlI4bk1xUXo0U3k1bzVxZGN2Z04rcXk4OHZVbVoxVGFWemZRekRyT3N1OEhDMXBuTzVVL0lIdk4KdzlWZFdjNUxQY0JqRkVWYlltRFAvdTlDVlN1bWhjb3Z6bjN2enlvMHBQa0k4U29Yby9PdjU1TDIxKzdFWnFoawpXejJlY0ZHRlliYU1rRnh0NjBhQVV6bTlEd0tCZ1FEMHdkS294M0szaE5ZS0srVHp4Unh4NUxOVWtVeUlGT1NSCjlzeHhRS0FGTS9iZTdvaEpSeUxrSERQdXYzcHJhM3p6ejQzeVFVY0lWb0w5MmlTU2c0NnRtQUJNYjQ5RnJ1ZE4KaU0xYktseUo1MmhHM0p5LzlhQkVVeGhMNEZXNWI5ODYyTlZBSjY5M1Nwcm5PYjdKVXFEdlZmZ0QvYmZTUjlZUwpzcGNpbFFpOFJRS0JnREd3R053dHpBdmFhYnAwLzJuUElYWkoyWEQ5eXhraDBDT3k3UUJOcDlpZktRTGo4UXB4CmV4MURoU3pIOEhlOXJieTY5OTFHWTM5bDcxZ1ZJRkdRQlJtOEs4RCtGN25KRDBCZVNkMktuRzFnb0FnYUl3YTcKZW5saWFmTUo1NGZjOXNDK0t3MWk2OXRZeGFWRXpRRXNqaFZ3SE1ZMnFFcEtZdlVua2N6a0wrRURBb0dCQU94UApNb0FJU1NjTTFzYnRTRmxmeHkzam5JMGEyQ1JPMzd4WjUxdTFCSXJoZUFvWGpYZ0tlWko1OUY0ZmV5ZTVtT09oClVqQkNmRDE5b1cxTXI2RFI2ZkNLNEVic013MFphSE5Ba056alVvTkc3RFAyamxUNzV1Znd2bldMdTlpVlBaY0kKZ1NRMjdMK2xSVmZZTmU4VW14TlpFbU53RkltdkYrM25oZW82c0R0dEFvR0FTS3Vkam85YUZ3NjZmUk1Wb0YrTQpuM1lsa0Q0MmF4RDRpUHkzcVNpZDVSYXkvbGxvUmNja3VFL2NmNG9lSnY5bVlYMWloRFRKYTgva1Zva0dzNEJ1CkI0MnRQQUFLSW5MK3BVamFnR0grdzFZanBXQnNocDljdDltM3pyb2JiSGhMenpVUVB1Zk5MN0ltOW1Xdyt6UDQKcWl3WlBpME9BV0VHVlpxeTVjM0FpOGc9Ci0tLS0tRU5EIFBSSVZBVEUgS0VZLS0tLS0=`
var PublicKey = `LS0tLS1CRUdJTiBQdWJsaWMgS2V5LS0tLS0KTUlJQklqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEFNSUlCQ2dLQ0FRRUE3ODN6MXVHSW43OG5tUm5sZm9QNQpSc25ibWRoQit0RG95eEFLQVBWU3JXZW5NZGdKRnZuMlJjYWlRQmVIMkJiS05RY1AveWdZZU9FQzJaaGIrRVpICkF3eVQwWm8zNHZaUUFnVldPaTV6VkROd3h1NFZQZmMrMjUrWU5BcitaeS90eHFKbFdGNUViUHlxalNaR1FRSS8KeHdaU3ZnVzc2dWN6MnZqMkxLb3lkZEVTcGJtVjBRUVlWeHZRMWdIb0JPREhYdWM2ZER3UVpNOWNQV2gvTi9uSAo5STQ1VHkwWldDVE9qOXFDay85MkNoTGFJL2hZNDU1MnlMRE1yenNxRXFyTDBLc0NRS3FuVHlDZVVncXZIRnhHCkswelp0MzVvYjBDN3AxRllKNEhISGZUSFY2VDBhOFUxQ1l1YTVOSk5pNW1QK05iNktJeXVXQ0hyUUdlS2VONzUKQ3dJREFRQUIKLS0tLS1FTkQgUHVibGljIEtleS0tLS0t`

// 使用RSA公钥加密数据
func encryptRSA(publicKey *rsa.PublicKey, data []byte) ([]byte, error) {
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, data)
	if err != nil {
		return nil, err
	}
	return ciphertext, nil
}

// 使用RSA私钥解密数据
func decryptRSA(privateKey *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertext)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

func main() {
	// 生成RSA密钥对
	//privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	//if err != nil {
	//	fmt.Println("Failed to generate RSA private key:", err)
	//	return
	//}
	//
	//publicKey := &privateKey.PublicKey
	bytePrivateKey, err := base64.StdEncoding.DecodeString(PrivateKey)
	if err != nil {
		log.Fatal(err)
	}
	priBlock, _ := pem.Decode(bytePrivateKey)
	priKey, err := x509.ParsePKCS8PrivateKey(priBlock.Bytes)
	if err != nil {
		log.Fatal(err)
	}

	bytePublicKey, err := base64.StdEncoding.DecodeString(PublicKey)
	if err != nil {
		log.Fatal(err)
	}

	pubBlock, _ := pem.Decode(bytePublicKey)
	pubKey, err := x509.ParsePKIXPublicKey(pubBlock.Bytes)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := pubKey.(*rsa.PublicKey)
	privateKey := priKey.(*rsa.PrivateKey)

	// 原始数据
	originalData := []byte("mGUiOQIOXI")

	// 使用公钥加密数据
	encryptedData, err := encryptRSA(publicKey, originalData)
	if err != nil {
		fmt.Println("Failed to encrypt data:", err)
		return
	}
	fmt.Println(base64.StdEncoding.EncodeToString(encryptedData))

	// 使用私钥解密数据
	decryptedData, err := decryptRSA(privateKey, encryptedData)
	if err != nil {
		fmt.Println("Failed to decrypt data:", err)
		return
	}

	fmt.Println("Original data:", string(originalData))
	fmt.Println("Decrypted data:", string(decryptedData))
}
