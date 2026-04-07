package wecom

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"math/rand"
	"sort"
	"strings"
)

// CallbackCrypto 处理企业微信回调消息的加解密。
type CallbackCrypto struct {
	token  string
	aesKey []byte
	corpID string
}

// NewCallbackCrypto 创建 CallbackCrypto 实例。
// encodingAESKey 是企业微信管理后台配置的 43 位 Base64 编码密钥。
func NewCallbackCrypto(token, encodingAESKey, corpID string) (*CallbackCrypto, error) {
	aesKey, err := base64.StdEncoding.DecodeString(encodingAESKey + "=")
	if err != nil {
		return nil, fmt.Errorf("wecom: invalid encoding aes key: %w", err)
	}
	if len(aesKey) != 32 {
		return nil, fmt.Errorf("wecom: aes key length must be 32, got %d", len(aesKey))
	}
	return &CallbackCrypto{token: token, aesKey: aesKey, corpID: corpID}, nil
}

// VerifyURL 处理企业微信回调 URL 验证的 GET 请求。
// 返回解密后的 echostr，需作为 HTTP 响应体原样返回。
func (c *CallbackCrypto) VerifyURL(msgSignature, timestamp, nonce, echostr string) ([]byte, error) {
	if !c.verifySignature(msgSignature, timestamp, nonce, echostr) {
		return nil, fmt.Errorf("wecom: signature verification failed")
	}
	plaintext, err := c.decrypt(echostr)
	if err != nil {
		return nil, fmt.Errorf("wecom: decrypt echostr failed: %w", err)
	}
	return plaintext, nil
}

// CallbackMessage 企业微信回调 POST 请求体的 XML 结构。
type CallbackMessage struct {
	XMLName    xml.Name `xml:"xml"`
	ToUserName string   `xml:"ToUserName"`
	Encrypt    string   `xml:"Encrypt"`
	AgentID    string   `xml:"AgentID"`
}

// DecryptMsg 解密回调 POST 消息。
func (c *CallbackCrypto) DecryptMsg(msgSignature, timestamp, nonce string, postBody []byte) ([]byte, error) {
	var msg CallbackMessage
	if err := xml.Unmarshal(postBody, &msg); err != nil {
		return nil, fmt.Errorf("wecom: unmarshal callback xml failed: %w", err)
	}
	if !c.verifySignature(msgSignature, timestamp, nonce, msg.Encrypt) {
		return nil, fmt.Errorf("wecom: signature verification failed")
	}
	plaintext, err := c.decrypt(msg.Encrypt)
	if err != nil {
		return nil, fmt.Errorf("wecom: decrypt message failed: %w", err)
	}
	return plaintext, nil
}

// EncryptMsg 加密回复消息。
func (c *CallbackCrypto) EncryptMsg(replyMsg, timestamp, nonce string) ([]byte, error) {
	encrypted, err := c.encrypt([]byte(replyMsg))
	if err != nil {
		return nil, fmt.Errorf("wecom: encrypt message failed: %w", err)
	}
	signature := c.signature(timestamp, nonce, encrypted)
	resp := fmt.Sprintf(`<xml><Encrypt><![CDATA[%s]]></Encrypt><MsgSignature><![CDATA[%s]]></MsgSignature><TimeStamp>%s</TimeStamp><Nonce><![CDATA[%s]]></Nonce></xml>`,
		encrypted, signature, timestamp, nonce)
	return []byte(resp), nil
}

// verifySignature 校验消息签名。
func (c *CallbackCrypto) verifySignature(msgSignature, timestamp, nonce, encrypt string) bool {
	return c.signature(timestamp, nonce, encrypt) == msgSignature
}

// signature 计算签名：SHA1(sort(token, timestamp, nonce, encrypt))。
func (c *CallbackCrypto) signature(timestamp, nonce, encrypt string) string {
	sl := []string{c.token, timestamp, nonce, encrypt}
	sort.Strings(sl)
	h := sha1.New()
	h.Write([]byte(strings.Join(sl, "")))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// decrypt 解密 Base64 编码的 AES-CBC 密文。
func (c *CallbackCrypto) decrypt(base64Msg string) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(base64Msg)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(c.aesKey)
	if err != nil {
		return nil, err
	}
	if len(ciphertext) < aes.BlockSize || len(ciphertext)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("wecom: invalid ciphertext length")
	}

	iv := c.aesKey[:aes.BlockSize]
	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	// 去除 PKCS#7 填充
	padLen := int(plaintext[len(plaintext)-1])
	if padLen < 1 || padLen > aes.BlockSize || padLen > len(plaintext) {
		return nil, fmt.Errorf("wecom: invalid pkcs7 padding")
	}
	plaintext = plaintext[:len(plaintext)-padLen]

	// 格式：16 字节随机数 + 4 字节消息长度（大端序）+ 消息体 + CorpID
	if len(plaintext) < 20 {
		return nil, fmt.Errorf("wecom: plaintext too short")
	}
	msgLen := binary.BigEndian.Uint32(plaintext[16:20])
	if uint32(len(plaintext)) < 20+msgLen {
		return nil, fmt.Errorf("wecom: message length mismatch")
	}
	msg := plaintext[20 : 20+msgLen]
	corpID := string(plaintext[20+msgLen:])
	if corpID != c.corpID {
		return nil, fmt.Errorf("wecom: corp_id mismatch: got %s, want %s", corpID, c.corpID)
	}
	return msg, nil
}

// encrypt 使用 AES-CBC 加密明文。
func (c *CallbackCrypto) encrypt(plaintext []byte) (string, error) {
	// 构造：16 字节随机数 + 4 字节消息长度 + 消息体 + CorpID
	randomBytes := make([]byte, 16)
	for i := range randomBytes {
		randomBytes[i] = byte(rand.Intn(256))
	}
	msgLenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(msgLenBuf, uint32(len(plaintext)))

	buf := make([]byte, 0, 16+4+len(plaintext)+len(c.corpID))
	buf = append(buf, randomBytes...)
	buf = append(buf, msgLenBuf...)
	buf = append(buf, plaintext...)
	buf = append(buf, []byte(c.corpID)...)

	// PKCS#7 填充
	padLen := aes.BlockSize - len(buf)%aes.BlockSize
	padding := make([]byte, padLen)
	for i := range padding {
		padding[i] = byte(padLen)
	}
	buf = append(buf, padding...)

	block, err := aes.NewCipher(c.aesKey)
	if err != nil {
		return "", err
	}
	iv := c.aesKey[:aes.BlockSize]
	mode := cipher.NewCBCEncrypter(block, iv)
	ciphertext := make([]byte, len(buf))
	mode.CryptBlocks(ciphertext, buf)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}
