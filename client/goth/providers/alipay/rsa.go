package alipay

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"strings"
)

const (
	kPublicKeyPrefix = "-----BEGIN PUBLIC KEY-----"
	kPublicKeySuffix = "-----END PUBLIC KEY-----"

	kPKCS1Prefix = "-----BEGIN RSA PRIVATE KEY-----"
	KPKCS1Suffix = "-----END RSA PRIVATE KEY-----"

	kPKCS8Prefix = "-----BEGIN PRIVATE KEY-----"
	KPKCS8Suffix = "-----END PRIVATE KEY-----"

	kPublicKeyType     = "PUBLIC KEY"
	kPrivateKeyType    = "PRIVATE KEY"
	kRSAPrivateKeyType = "RSA PRIVATE KEY"
)

var (
	ErrPrivateKeyFailedToLoad = errors.New("[oauth2s]alipay: private key failed to load")
	ErrPublicKeyFailedToLoad  = errors.New("[oauth2s]alipay: public key failed to load")

	pkcsPrefixAndSuffixReplacer = strings.NewReplacer(
		kPKCS8Prefix, "",
		KPKCS8Suffix, "",
		kPKCS1Prefix, "",
		KPKCS1Suffix, "",
	)
	spaceReplacer = strings.NewReplacer(
		" ", "",
		"\n", "",
		"\r", "",
		"\t", "",
	)
)

func formatKey(raw, prefix, suffix string, lineCount int) []byte {
	if raw == "" {
		return nil
	}
	if len(prefix) > 0 {
		raw = strings.Replace(raw, prefix, "", 1)
	}
	if len(suffix) > 0 {
		raw = strings.Replace(raw, suffix, "", 1)
	}
	raw = spaceReplacer.Replace(raw)

	var sl = len(raw)
	var c = sl / lineCount
	if sl%lineCount > 0 {
		c = c + 1
	}

	var buf bytes.Buffer
	buf.WriteString(prefix + "\n")
	for i := 0; i < c; i++ {
		var b = i * lineCount
		var e = b + lineCount
		if e > sl {
			buf.WriteString(raw[b:])
		} else {
			buf.WriteString(raw[b:e])
		}
		buf.WriteString("\n")
	}
	buf.WriteString(suffix)
	return buf.Bytes()
}

func FormatPublicKey(raw string) []byte {
	return formatKey(raw, kPublicKeyPrefix, kPublicKeySuffix, 64)
}

func FormatPKCS1PrivateKey(raw string) []byte {
	raw = pkcsPrefixAndSuffixReplacer.Replace(raw)
	return formatKey(raw, ``, ``, 64)
}

func FormatPKCS8PrivateKey(raw string) []byte {
	raw = pkcsPrefixAndSuffixReplacer.Replace(raw)
	return formatKey(raw, ``, ``, 64)
}

func ParsePKCS1PrivateKey(data []byte) (key *rsa.PrivateKey, err error) {
	var block *pem.Block
	block, _ = pem.Decode(data)
	if block == nil {
		return nil, ErrPrivateKeyFailedToLoad
	}

	key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return key, err
}

func ParsePKCS8PrivateKey(data []byte) (key *rsa.PrivateKey, err error) {
	var block *pem.Block
	block, _ = pem.Decode(data)
	if block == nil {
		return nil, ErrPrivateKeyFailedToLoad
	}

	rawKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	key, ok := rawKey.(*rsa.PrivateKey)
	if !ok {
		return nil, ErrPrivateKeyFailedToLoad
	}

	return key, err
}

func ParsePublicKey(data []byte) (key *rsa.PublicKey, err error) {
	var block *pem.Block
	block, _ = pem.Decode(data)
	if block == nil {
		return nil, ErrPublicKeyFailedToLoad
	}

	var pubInterface interface{}
	pubInterface, err = x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	key, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, ErrPublicKeyFailedToLoad
	}

	return key, err
}

// RSAEncryptWithKey 使用公钥 key 对数据 data 进行 RSA 加密
func RSAEncryptWithKey(plaintext []byte, key *rsa.PublicKey) ([]byte, error) {
	var pData = packageData(plaintext, key.N.BitLen()/8-11)
	var ciphertext = make([]byte, 0, len(pData))

	for _, d := range pData {
		var c, e = rsa.EncryptPKCS1v15(rand.Reader, key, d)
		if e != nil {
			return nil, e
		}
		ciphertext = append(ciphertext, c...)
	}

	return ciphertext, nil
}

// RSADecryptWithPKCS1 使用私钥 key 对数据 data 进行 RSA 解密，key 的格式为 pkcs1
func RSADecryptWithPKCS1(ciphertext, key []byte) ([]byte, error) {
	priKey, err := ParsePKCS1PrivateKey(key)
	if err != nil {
		return nil, err
	}

	return RSADecryptWithKey(ciphertext, priKey)
}

// RSADecryptWithPKCS8 使用私钥 key 对数据 data 进行 RSA 解密，key 的格式为 pkcs8
func RSADecryptWithPKCS8(ciphertext, key []byte) ([]byte, error) {
	priKey, err := ParsePKCS8PrivateKey(key)
	if err != nil {
		return nil, err
	}

	return RSADecryptWithKey(ciphertext, priKey)
}

// RSADecryptWithKey 使用私钥 key 对数据 data 进行 RSA 解密
func RSADecryptWithKey(ciphertext []byte, key *rsa.PrivateKey) ([]byte, error) {
	var pData = packageData(ciphertext, key.PublicKey.N.BitLen()/8)
	var plaintext = make([]byte, 0, len(pData))

	for _, d := range pData {
		var p, e = rsa.DecryptPKCS1v15(rand.Reader, key, d)
		if e != nil {
			return nil, e
		}
		plaintext = append(plaintext, p...)
	}
	return plaintext, nil
}

func RSASignWithPKCS1(plaintext, key []byte, hash crypto.Hash) ([]byte, error) {
	priKey, err := ParsePKCS1PrivateKey(key)
	if err != nil {
		return nil, err
	}
	return RSASignWithKey(plaintext, priKey, hash)
}

func RSASignWithPKCS8(plaintext, key []byte, hash crypto.Hash) ([]byte, error) {
	priKey, err := ParsePKCS8PrivateKey(key)
	if err != nil {
		return nil, err
	}
	return RSASignWithKey(plaintext, priKey, hash)
}

func RSASignWithKey(plaintext []byte, key *rsa.PrivateKey, hash crypto.Hash) ([]byte, error) {
	var h = hash.New()
	h.Write(plaintext)
	var hashed = h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, key, hash, hashed)
}

func RSAVerify(ciphertext, sign, key []byte, hash crypto.Hash) error {
	pubKey, err := ParsePublicKey(key)
	if err != nil {
		return err
	}
	return RSAVerifyWithKey(ciphertext, sign, pubKey, hash)
}

func RSAVerifyWithKey(ciphertext, sign []byte, key *rsa.PublicKey, hash crypto.Hash) error {
	var h = hash.New()
	h.Write(ciphertext)
	var hashed = h.Sum(nil)
	return rsa.VerifyPKCS1v15(key, hash, hashed, sign)
}

func getPublicKeyBytes(publicKey *rsa.PublicKey) ([]byte, error) {
	pubDer, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, err
	}

	pubBlock := &pem.Block{Type: kPublicKeyType, Bytes: pubDer}

	var pubBuf bytes.Buffer
	if err = pem.Encode(&pubBuf, pubBlock); err != nil {
		return nil, err
	}
	return pubBuf.Bytes(), nil
}

func GenRSAKeyWithPKCS1(bits int) (privateKey, publicKey []byte, err error) {
	priKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	priDer := x509.MarshalPKCS1PrivateKey(priKey)
	priBlock := &pem.Block{Type: kRSAPrivateKeyType, Bytes: priDer}

	var priBuf bytes.Buffer
	if err = pem.Encode(&priBuf, priBlock); err != nil {
		return nil, nil, err
	}

	publicKey, err = getPublicKeyBytes(&priKey.PublicKey)
	if err != nil {
		return nil, nil, err
	}
	privateKey = priBuf.Bytes()
	return privateKey, publicKey, err
}

func GenRSAKeyWithPKCS8(bits int) (privateKey, publicKey []byte, err error) {
	priKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	priDer, err := x509.MarshalPKCS8PrivateKey(priKey)
	if err != nil {
		return nil, nil, err
	}
	priBlock := &pem.Block{Type: kPrivateKeyType, Bytes: priDer}

	var priBuf bytes.Buffer
	if err = pem.Encode(&priBuf, priBlock); err != nil {
		return nil, nil, err
	}

	publicKey, err = getPublicKeyBytes(&priKey.PublicKey)
	if err != nil {
		return nil, nil, err
	}
	privateKey = priBuf.Bytes()

	return privateKey, publicKey, err
}
