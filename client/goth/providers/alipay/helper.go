package alipay

import (
	"crypto"
	"encoding/base64"
	"net/url"
)

func SignRSA2(param url.Values, privateKey []byte) string {
	if param == nil {
		param = make(url.Values, 0)
	}

	src := param.Encode()
	sig, err := SignPKCS1v15([]byte(src), privateKey, crypto.SHA256)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(sig)
}

func SignRSA(param url.Values, privateKey []byte) string {
	if param == nil {
		param = make(url.Values, 0)
	}

	src := param.Encode()
	sig, err := SignPKCS1v15([]byte(src), privateKey, crypto.SHA1)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(sig)
}

func VerifySign(val url.Values, key []byte) error {
	sign, err := base64.StdEncoding.DecodeString(val.Get("sign"))
	if err != nil {
		return err
	}
	signType := val.Get("sign_type")
	if _, ok := val[`sign`]; ok {
		val.Del("sign")
	}
	if _, ok := val[`sign_type`]; ok {
		val.Del("sign_type")
	}
	s := val.Encode()

	if signType == `RSA` {
		err = VerifyPKCS1v15([]byte(s), sign, key, crypto.SHA1)
	} else {
		err = VerifyPKCS1v15([]byte(s), sign, key, crypto.SHA256)
	}
	return err
}

func VerifyResponseData(data []byte, signType, sign string, key []byte) error {
	signBytes, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return err
	}

	if signType == `RSA` {
		err = VerifyPKCS1v15(data, signBytes, key, crypto.SHA1)
	} else {
		err = VerifyPKCS1v15(data, signBytes, key, crypto.SHA256)
	}
	return err
}

func FormatNormalKey(b []byte, isPublicKey bool) (r []byte) {
	var name string
	if isPublicKey {
		name = `PUBLIC`
	} else {
		name = `PRIVATE`
	}
	r = append(r, []byte(`-----BEGIN `+name+` KEY-----`)...)
	for p, j := 1, len(b); ; p++ {
		offset := (p - 1) * 64
		end := offset + 64
		r = append(r, '\n')
		if end < j {
			r = append(r, b[offset:end]...)
			continue
		}
		r = append(r, b[offset:]...)
		break
	}
	r = append(r, '\n')
	r = append(r, []byte(`-----END `+name+` KEY-----`)...)
	return
}
