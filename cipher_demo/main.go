package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	mathrand "math/rand"
	"strconv"
	"time"
)

var PrivateKey = `LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1JSUV2UUlCQURBTkJna3Foa2lHOXcwQkFRRUZBQVNDQktjd2dnU2pBZ0VBQW9JQkFRRHZ6ZlBXNFlpZnZ5ZVoKR2VWK2cvbEd5ZHVaMkVINjBPakxFQW9BOVZLdFo2Y3gyQWtXK2ZaRnhxSkFGNGZZRnNvMUJ3Ly9LQmg0NFFMWgptRnY0UmtjRERKUFJtamZpOWxBQ0JWWTZMbk5VTTNERzdoVTk5ejdibjVnMEN2NW5MKzNHb21WWVhrUnMvS3FOCkprWkJBai9IQmxLK0JidnE1elBhK1BZc3FqSjEwUktsdVpYUkJCaFhHOURXQWVnRTRNZGU1enAwUEJCa3oxdzkKYUg4MytjZjBqamxQTFJsWUpNNlAyb0tULzNZS0V0b2orRmpqbm5iSXNNeXZPeW9TcXN2UXF3SkFxcWRQSUo1UwpDcThjWEVZclRObTNmbWh2UUx1blVWZ25nY2NkOU1kWHBQUnJ4VFVKaTVyazBrMkxtWS80MXZvb2pLNVlJZXRBClo0cDQzdmtMQWdNQkFBRUNnZ0VBR1lEejMvU1lqVFROUjFFandUTFQvaDFWeDZUaVQ0U01YWnhWRkFrRFJBSDAKSEM3M3VJUFpGMDZxenRveHpsK09IZG1rYitTWnhiRllsai9IK0QyeEs3Zll1TUFJRlpGUXJRWllSMVBOWERVawpWMlBXeW9KVEl4UjBJWVRSemJPaFBERGxTbUtzTUZNbTZHQWJ0S3BDa2k0ditwbXRoS3dXTGNHd25NdDFGbVp6Ck9tVXNWUkg3a01rL1doVzNvb1B3Ny9ZVGhNSW1zZkJoNm1QT3VpRm9IRHRHNnNsWmJ5YTRhdjcrRlQ5eldOWHIKUzNkTmtVZFR6eENtcC9TV1paSWV6UlR3R0wvaXRvdTVtN2owSjVPbkJDYVNxL0lWR1U1UXoydm1ydzl4LzQyagpIZXVzWXhKZGlSQUJuR3N3Tm1xL0x4VXYyNVBuckptUDQ3MTY3VHhYc1FLQmdRRDYwZVBqTDJNVGx0REF2U3crClcvMkJDcFZTYlI4bk1xUXo0U3k1bzVxZGN2Z04rcXk4OHZVbVoxVGFWemZRekRyT3N1OEhDMXBuTzVVL0lIdk4KdzlWZFdjNUxQY0JqRkVWYlltRFAvdTlDVlN1bWhjb3Z6bjN2enlvMHBQa0k4U29Yby9PdjU1TDIxKzdFWnFoawpXejJlY0ZHRlliYU1rRnh0NjBhQVV6bTlEd0tCZ1FEMHdkS294M0szaE5ZS0srVHp4Unh4NUxOVWtVeUlGT1NSCjlzeHhRS0FGTS9iZTdvaEpSeUxrSERQdXYzcHJhM3p6ejQzeVFVY0lWb0w5MmlTU2c0NnRtQUJNYjQ5RnJ1ZE4KaU0xYktseUo1MmhHM0p5LzlhQkVVeGhMNEZXNWI5ODYyTlZBSjY5M1Nwcm5PYjdKVXFEdlZmZ0QvYmZTUjlZUwpzcGNpbFFpOFJRS0JnREd3R053dHpBdmFhYnAwLzJuUElYWkoyWEQ5eXhraDBDT3k3UUJOcDlpZktRTGo4UXB4CmV4MURoU3pIOEhlOXJieTY5OTFHWTM5bDcxZ1ZJRkdRQlJtOEs4RCtGN25KRDBCZVNkMktuRzFnb0FnYUl3YTcKZW5saWFmTUo1NGZjOXNDK0t3MWk2OXRZeGFWRXpRRXNqaFZ3SE1ZMnFFcEtZdlVua2N6a0wrRURBb0dCQU94UApNb0FJU1NjTTFzYnRTRmxmeHkzam5JMGEyQ1JPMzd4WjUxdTFCSXJoZUFvWGpYZ0tlWko1OUY0ZmV5ZTVtT09oClVqQkNmRDE5b1cxTXI2RFI2ZkNLNEVic013MFphSE5Ba056alVvTkc3RFAyamxUNzV1Znd2bldMdTlpVlBaY0kKZ1NRMjdMK2xSVmZZTmU4VW14TlpFbU53RkltdkYrM25oZW82c0R0dEFvR0FTS3Vkam85YUZ3NjZmUk1Wb0YrTQpuM1lsa0Q0MmF4RDRpUHkzcVNpZDVSYXkvbGxvUmNja3VFL2NmNG9lSnY5bVlYMWloRFRKYTgva1Zva0dzNEJ1CkI0MnRQQUFLSW5MK3BVamFnR0grdzFZanBXQnNocDljdDltM3pyb2JiSGhMenpVUVB1Zk5MN0ltOW1Xdyt6UDQKcWl3WlBpME9BV0VHVlpxeTVjM0FpOGc9Ci0tLS0tRU5EIFBSSVZBVEUgS0VZLS0tLS0=`
var PublicKey = `LS0tLS1CRUdJTiBQdWJsaWMgS2V5LS0tLS0KTUlJQklqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEFNSUlCQ2dLQ0FRRUE3ODN6MXVHSW43OG5tUm5sZm9QNQpSc25ibWRoQit0RG95eEFLQVBWU3JXZW5NZGdKRnZuMlJjYWlRQmVIMkJiS05RY1AveWdZZU9FQzJaaGIrRVpICkF3eVQwWm8zNHZaUUFnVldPaTV6VkROd3h1NFZQZmMrMjUrWU5BcitaeS90eHFKbFdGNUViUHlxalNaR1FRSS8KeHdaU3ZnVzc2dWN6MnZqMkxLb3lkZEVTcGJtVjBRUVlWeHZRMWdIb0JPREhYdWM2ZER3UVpNOWNQV2gvTi9uSAo5STQ1VHkwWldDVE9qOXFDay85MkNoTGFJL2hZNDU1MnlMRE1yenNxRXFyTDBLc0NRS3FuVHlDZVVncXZIRnhHCkswelp0MzVvYjBDN3AxRllKNEhISGZUSFY2VDBhOFUxQ1l1YTVOSk5pNW1QK05iNktJeXVXQ0hyUUdlS2VONzUKQ3dJREFRQUIKLS0tLS1FTkQgUHVibGljIEtleS0tLS0t`

func RsaDecode(strPlainText string) (string, error) {
	plainText, err := base64.StdEncoding.DecodeString(strPlainText)
	if err != nil {
		return "", fmt.Errorf("base64 decode error: %v", err)
	}
	bytePrivateKey, err := base64.StdEncoding.DecodeString(PrivateKey)
	if err != nil {
		return "", err
	}

	priBlock, _ := pem.Decode(bytePrivateKey)
	priKey, err := x509.ParsePKCS8PrivateKey(priBlock.Bytes)
	if err != nil {
		return "", err
	}
	decryptText, err := rsa.DecryptPKCS1v15(rand.Reader, priKey.(*rsa.PrivateKey), plainText)
	if err != nil {
		return "", fmt.Errorf("rsa decrypt error: %v", err)
	}
	return string(decryptText), nil
}

func RsaEncode(strPlainText string) (string, error) {
	plainText := []byte(strPlainText)
	bytePublicKey, err := base64.StdEncoding.DecodeString(PublicKey)
	if err != nil {
		return "", err
	}

	pubBlock, _ := pem.Decode(bytePublicKey)
	pubKey, err := x509.ParsePKIXPublicKey(pubBlock.Bytes)
	if err != nil {
		return "", err
	}

	rsaPubKey := pubKey.(*rsa.PublicKey)

	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPubKey, plainText)
	if err != nil {
		return "", err
	}

	base64CipherText := base64.StdEncoding.EncodeToString(cipherText)
	return base64CipherText, nil
}

type Data struct {
	Colume []string
	Value  []DataItem
}

type DataItem struct {
	Name   string `json:"name"`
	Value  int    `json:"value"`
	Id     int    `json:"id"`
	Method string `json:"method"`
	Path   string `json:"path"`
}

// GenerateDataItem 生成10条DataItem的数据
func GenerateDataItem() []DataItem {
	var dataItem []DataItem
	for i := 0; i < 100; i++ {
		dataItem = append(dataItem, DataItem{
			Name:   "name" + strconv.Itoa(i),
			Value:  i,
			Id:     i,
			Method: "GET",
			Path:   "中国",
		})
	}
	return dataItem
}

func generateKey(length int) (string, error) {
	key := make([]byte, length)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}

	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	for i := 0; i < length; i++ {
		key[i] = charset[int(key[i])%len(charset)]
	}

	return string(key), nil
}

func Cipher() {
	test := Data{
		Colume: []string{"id",
			"path",
			"description",
			"api_group",
			"method",
			"action",
			"created_at",
			"updated_at",
			"deleted_at"},
		Value: GenerateDataItem(),
	}

	bytes, err := json.Marshal(test)
	if err != nil {
		log.Fatal(err)
	}
	base64Data := base64.RawStdEncoding.EncodeToString(bytes)
	//	 随机生成一个key 16位
	key, err := generateKey(10)
	if err != nil {
		log.Fatal(err)
	}
	//加密key
	encryptKey, err := RsaEncode(key)
	fmt.Println("key", key)
	fmt.Println()
	fmt.Println("加密之后的key", encryptKey)

	// 按照index=3进行分割，将base64Data后面的放到前面
	fmt.Println("base64Data", base64Data)

	// 生成一个随机数
	mathrand.Seed(time.Now().UnixNano())
	index := mathrand.Intn(len(base64Data))
	fmt.Println("随机数", index)

	base64Data = base64Data[index:] + base64Data[:index]
	fmt.Println()
	fmt.Println("分割之后的数据", base64Data)

	fmt.Println("加密key", encryptKey)
	fmt.Println()
	split := reverseString(base64Data)
	encryptData := key + split
	fmt.Println("加密之后的数据", encryptData)
}

func main() {
	Cipher()
}

func reverseString(s string) string {
	// 将字符串转换为字节数组
	bytes := []byte(s)

	// 使用 bytes 包的 Reverse 函数进行字节数组反转
	bytes = reverseBytes(bytes)

	// 将字节数组转换回字符串并返回
	return string(bytes)
}

func reverseBytes(bytes []byte) []byte {
	// 使用 bytes 包的 Reverse 函数对字节数组进行反转
	reversed := bytes[:]
	bytesLen := len(reversed)
	for i := 0; i < bytesLen/2; i++ {
		j := bytesLen - i - 1
		reversed[i], reversed[j] = reversed[j], reversed[i]
	}

	return reversed
}
