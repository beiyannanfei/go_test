package encoding

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/errors"
	"io"
	"os"
	"regexp"
	"strings"
)

// 生成私钥
func GetRSAPrivateKey(key string) (*rsa.PrivateKey, error) {
	//privateKey := "-----BEGIN PRIVATE KEY-----\n"+key+"\n-----END PRIVATE KEY-----\n"
	block, _ := pem.Decode([]byte(key))
	if nil == block {
		logrus.WithFields(logrus.Fields{
			"error": "error pem decode",
		}).Error("GetRSAPrivateKey")
		return nil, errors.New("GetRSAPrivateKey failed")
	}
	result, err1 := x509.ParsePKCS8PrivateKey(block.Bytes)
	if nil != err1 {
		logrus.WithFields(logrus.Fields{
			"result": result,
			"err":    err1,
		}).Error("error x509 ParsePKCS8PrivateKey")
		return nil, err1
	}

	switch key := result.(type) {
	case *rsa.PrivateKey:
		{
			return key, nil
		}
	default:
		return nil, errors.New("not private key")
	}
}

func GetRSAPublicKey(key []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(key)
	if nil == block {
		logrus.WithFields(logrus.Fields{
			"error": "error pem decode",
		}).Error("GetRSAPublicKey")
		return nil, errors.New("GetRSAPublicKey failed")
	}

	cert, err := x509.ParsePKIXPublicKey(block.Bytes)
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("GetRSAPublicKey")
		return nil, err
	}

	switch key := cert.(type) {
	case *rsa.PublicKey:
		{
			return key, nil
		}
	default:
		return nil, errors.New("not public key")
	}
}

// sha256withRSA 签名
func SignSHA256WithRSA(data string, key *rsa.PrivateKey) (string, error) {
	hash := sha256.Sum256([]byte(data))
	rng := rand.Reader
	signature, err := rsa.SignPKCS1v15(rng, key, crypto.SHA256, hash[:])
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("error rsa SignPKCS1v15 ")
		return "", errors.New("error from signing: " + data)
	}
	return base64.StdEncoding.EncodeToString(signature), nil
	//return signature, nil
}

func VerifySHA1WithRSA(data string, signature []byte, key *rsa.PublicKey) error {
	hash := sha1.Sum([]byte(data))
	err := rsa.VerifyPKCS1v15(key, crypto.SHA1, hash[:], signature)
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"err":       err,
			"data":      data,
			"signature": string(signature),
		}).Error("failed VerifySHA1WithRSA")
		return nil
	}
	return nil
}

func VerifySHA256WithRSA(data string, signature []byte, key *rsa.PublicKey) error {
	hash := sha256.Sum256([]byte(data))
	err := rsa.VerifyPKCS1v15(key, crypto.SHA256, hash[:], signature)
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"err":       err,
			"data":      data,
			"signature": string(signature),
		}).Error("failed VerifySHA256WithRSA")
		return err
	}

	return nil
}

func GenPemFile(key string, length int, publicKey bool) error {
	fileName := ""
	if true == publicKey {
		fileName = "public"
	} else {
		fileName = "private"
	}

	f, err := os.OpenFile(fileName+".pem", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0600)
	defer f.Close()
	if nil != err {
		return err
	}

	_, err = io.WriteString(f, "-----BEGIN "+strings.ToUpper(fileName)+" KEY-----\n")
	if nil != err {
		return err
	}
	reg := regexp.MustCompile("\\s+")
	if length < 0 {
		key = reg.ReplaceAllString(key, "")
		_, err := io.WriteString(f, key+"\n")
		if nil != err {
			return err
		}
	} else {
		for index := 0; index < len(key); index += 64 {
			line := ""
			if index+length <= len(key) {
				line = string([]byte(key))[index:index+length]
			} else {
				line = string([]byte(key))[index: len(key)]
			}
			line = reg.ReplaceAllString(line, "")
			if len(line) > 0 {
				_, err = io.WriteString(f, line+"\n")
				if nil != err {
					return err
				}
			} else {
				break
			}
		}
	}

	_, err = io.WriteString(f, "-----END "+strings.ToUpper(fileName)+" KEY-----\n")
	return err
}
