package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	r "math/rand"
	"sort"
	"strings"
	"time"
)

type DingTalkCrypto struct {
	Token          string
	EncodingAESKey string
	SuiteKey       string
	BKey           []byte
	Block          cipher.Block
}

func NewDingTalkCrypto(token, encodingAESKey, suiteKey string) *DingTalkCrypto {
	fmt.Println(len(encodingAESKey))
	if len(encodingAESKey) != int(43) {
		panic("不合法的EncodingAESKey")
	}
	bkey, err := base64.StdEncoding.DecodeString(encodingAESKey + "=")
	if err != nil {
		panic(err.Error())
	}
	block, err := aes.NewCipher(bkey)
	if err != nil {
		panic(err.Error())
	}
	c := &DingTalkCrypto{
		Token:          token,
		EncodingAESKey: encodingAESKey,
		SuiteKey:       suiteKey,
		BKey:           bkey,
		Block:          block,
	}
	return c
}

func (c *DingTalkCrypto) GetDecryptMsg(signature, timestamp, nonce, secretMsg string) (string, error) {
	if !c.VerificationSignature(c.Token, timestamp, nonce, secretMsg, signature) {
		return "", errors.New("ERROR: 签名不匹配")
	}
	decode, err := base64.StdEncoding.DecodeString(secretMsg)
	if err != nil {
		return "", err
	}
	if len(decode) < aes.BlockSize {
		return "", errors.New("ERROR: 密文太短")
	}
	blockMode := cipher.NewCBCDecrypter(c.Block, c.BKey[:c.Block.BlockSize()])
	plantText := make([]byte, len(decode))
	blockMode.CryptBlocks(plantText, decode)
	plantText = pkCS7UnPadding(plantText)
	size := binary.BigEndian.Uint32(plantText[16:20])
	plantText = plantText[20:]
	corpID := plantText[size:]
	if string(corpID) != c.SuiteKey {
		return "", errors.New("ERROR: CorpID匹配不正确")
	}
	return string(plantText[:size]), nil
}

func (c *DingTalkCrypto) GetEncryptMsg(msg string) (map[string]string, error) {
	var timestamp = time.Now().Second()
	var nonce = randomString(12)
	str, sign, err := c.GetEncryptMsgDetail(msg, fmt.Sprint(timestamp), nonce)

	return map[string]string{"nonce": nonce, "timeStamp": fmt.Sprint(timestamp), "encrypt": str, "msg_signature": sign}, err
}
func (c *DingTalkCrypto) GetEncryptMsgDetail(msg, timestamp, nonce string) (string, string, error) {
	size := make([]byte, 4)
	binary.BigEndian.PutUint32(size, uint32(len(msg)))
	msg = randomString(16) + string(size) + msg + c.SuiteKey
	plantText := pkCS7Padding([]byte(msg), c.Block.BlockSize())
	if len(plantText)%aes.BlockSize != 0 {
		return "", "", errors.New("ERROR: 消息体size不为16的倍数")
	}
	blockMode := cipher.NewCBCEncrypter(c.Block, c.BKey[:c.Block.BlockSize()])
	chipherText := make([]byte, len(plantText))
	blockMode.CryptBlocks(chipherText, plantText)
	outMsg := base64.StdEncoding.EncodeToString(chipherText)
	signature := c.CreateSignature(c.Token, timestamp, nonce, string(outMsg))
	return string(outMsg), signature, nil
}

func sha1Sign(s string) string {
	// The pattern for generating a hash is `sha1.New()`,
	// `sha1.Write(bytes)`, then `sha1.Sum([]byte{})`.
	// Here we start with a new hash.
	h := sha1.New()

	// `Write` expects bytes. If you have a string `s`,
	// use `[]byte(s)` to coerce it to bytes.
	h.Write([]byte(s))

	// This gets the finalized hash result as a byte
	// slice. The argument to `Sum` can be used to append
	// to an existing byte slice: it usually isn't needed.
	bs := h.Sum(nil)

	// SHA1 values are often printed in hex, for example
	// in git commits. Use the `%x` format verb to convert
	// a hash results to a hex string.
	return fmt.Sprintf("%x", bs)
}

// 数据签名
func (c *DingTalkCrypto) CreateSignature(token, timestamp, nonce, msg string) string {
	params := make([]string, 0)
	params = append(params, token)
	params = append(params, timestamp)
	params = append(params, nonce)
	params = append(params, msg)
	sort.Strings(params)
	return sha1Sign(strings.Join(params, ""))
}

// 验证数据签名
func (c *DingTalkCrypto) VerificationSignature(token, timestamp, nonce, msg, sigture string) bool {
	return c.CreateSignature(token, timestamp, nonce, msg) == sigture
}

// 解密补位
func pkCS7UnPadding(plantText []byte) []byte {
	length := len(plantText)
	unpadding := int(plantText[length-1])
	return plantText[:(length - unpadding)]
}

// 加密补位
func pkCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// 随机字符串
func randomString(n int, alphabets ...byte) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	var randby bool
	if num, err := rand.Read(bytes); num != n || err != nil {
		r.Seed(time.Now().UnixNano())
		randby = true
	}
	for i, b := range bytes {
		if len(alphabets) == 0 {
			if randby {
				bytes[i] = alphanum[r.Intn(len(alphanum))]
			} else {
				bytes[i] = alphanum[b%byte(len(alphanum))]
			}
		} else {
			if randby {
				bytes[i] = alphabets[r.Intn(len(alphabets))]
			} else {
				bytes[i] = alphabets[b%byte(len(alphabets))]
			}
		}
	}
	return string(bytes)
}

var data = `
{
    "processInstanceId": "EyUhKga1SQO_xTHp1T2euA02721704350609",
    "eventId": "7948c32bb9d545afbdf901c23d8ca007",
    "corpId": "ding83662b550e82d42abc961a6cb783455b",
    "resource": "/v1.0/event/bpms_instance_change/processCode/PROC-149294E5-56B3-45A9-B459-EFE633D227DF/type/start",
    "EventType": "bpms_instance_change",
    "businessId": "202401021528000280874",
    "title": "沈立兴提交的应用自动部署",
    "type": "start",
    "url": "https://aflow.dingtalk.com/dingtalk/mobile/homepage.htm?corpid=ding83662b550e82d42abc961a6cb783455b&dd_share=false&showmenu=false&dd_progress=false&back=native&procInstId=M1795P3hQoKO94q9YoRjdg02721704180508&taskId=&swfrom=isv&dinghash=approval&dtaction=os&dd_from=corp#approval",
    "createTime": 1704180508000,
    "processCode": "PROC-149294E5-56B3-45A9-B459-EFE633D227DF",
    "bizCategoryId": "",
    "businessType": "",
    "staffId": "083114326827695761"
}
`

func main() {
	var ding = NewDingTalkCrypto("j2vYfxwzAhAzV3uayCDoz5RdFjfyPB", "POhRQTAqjKY9pciuNRMdNRkt6fiPq68tpeNROJTpxRX", "dingayb2nax1cqwishln")
	msg, _ := ding.GetEncryptMsg(data)
	fmt.Println(msg)

	success, _ := ding.GetDecryptMsg(msg["msg_signature"], msg["timeStamp"], msg["nonce"], msg["encrypt"])
	//success, _ := ding.GetDecryptMsg("c334e9ff41b7700b78cb2ae5f1be43fcb6b98f7b", "1704165524173", "xepiE4D3", "Gc/+BX+LAV8K5WFAm6A9pXVsZDIcHelh0jW0wv1wXUG0Ll4rs7TuecyxYzhd32fiLP00GZ6pRUwGBZgaCLtkX9tItzFN8IR26G3hct+wOUdkB5ZGN1BTUDPEIUXqw6vS")

	fmt.Println(success)
}
