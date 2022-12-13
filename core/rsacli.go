// rsa encrypt decrypt
package core

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
)

type rsaClient struct {
	cli *rsa.PrivateKey
}

var rsaCli *rsaClient

func init() {
	cli, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}

	rsaCli = &rsaClient{
		cli: cli,
	}
}

func GetRsaClient() *rsaClient {
	return rsaCli
}

func (r *rsaClient) GetPublicKeyStr() string {
	pubKey := r.cli.PublicKey
	pubBytes := x509.MarshalPKCS1PublicKey(&pubKey)
	pubBlock := pem.Block{Type: "RSA PUBLIC KEY", Bytes: pubBytes}
	return string(pem.EncodeToMemory(&pubBlock))
}

func (r *rsaClient) Decrypt(cipherStr string) (string, error) {
	cipherBytes, err := base64.StdEncoding.DecodeString(cipherStr)
	if err != nil {
		return "", err
	}
	plainBytes, err := rsa.DecryptPKCS1v15(rand.Reader, r.cli, cipherBytes)
	if err != nil {
		return "", err
	}
	return string(plainBytes), nil
}

func (r *rsaClient) Encrypt(plainStr string) (string, error) {
	cipherBytes, err := rsa.EncryptPKCS1v15(rand.Reader, &r.cli.PublicKey, []byte(plainStr))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(cipherBytes), nil
}

func (r *rsaClient) EncryptByPubkey(pubkey string, plainStr string) (string, error) {
	pubBlock, _ := pem.Decode([]byte(pubkey))
	pubObj, err := x509.ParsePKCS1PublicKey(pubBlock.Bytes)
	if err != nil {
		return "", err
	}
	cipherBytes, err := rsa.EncryptPKCS1v15(rand.Reader, pubObj, []byte(plainStr))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(cipherBytes), nil
}
