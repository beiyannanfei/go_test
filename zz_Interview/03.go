package main

import (
	"time"
	"strconv"
	"sort"
	"crypto/rsa"
	"encoding/pem"
	"crypto/x509"
	"crypto/sha256"
	"crypto"
	"errors"
	"crypto/rand"
	"fmt"
	"net/url"
	"encoding/base64"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

type DoAuthTokenRsp struct {
	RtnCode int    `form:"rtnCode"  json:"rtnCode"  binding:"required"` // 0=成功, -1=失败  1=接口鉴权失败  3001=参数错误
	Ts      string `form:"ts"       json:"ts"       binding:"required"` // 时间戳，接口返回的当前时间戳
	RtnSign string `form:"rtnSign"  json:"rtnSign"  binding:"required"` // 返回参数签名值，请根据rtnCode和ts校验签名是否正确
	ErrMsg  string `form:"errMsg"  json:"errMsg"`
}

//authUrl: https://gss-cn.game.hicloud.com/gameservice/api/gbClientApi
//privateKey: MIIEvgIBADANBgkqhkiG9w0BAQEFAASC953TpmniTM0Hcet1xaXpq4uSMOxv*********************maeWak6wvr3MjkGKsFrnF0uyEO94bZNgClpPp8XcQ6Y0m3O/7f7p+1dLhXkWXo+LNCUm5m4XSGv3V8ox3t+qR44MYnQFwia1yiCMr80W/CJAgMBAAECggEAG5EFQNSm8JHa8+va5UEiKDetOA1QHpmaGSj2EOE7urUKrId0HkNGMf9TeeXgGmY9DuVqnzpu3YDA/8agEPlsklnkxGSDtfedEF3inopkohg03EQX49lR07+VQbJXb88OH8Qn3Chol1CL1tlH6L4gjZ9C+AwdGpgyV8Sxvkjq59ckmABcF0C4qIKVZxpkITqyqD5cLrHw6D10bHv0y0l9dQDN0vUkt7jN8kvGZxJQ+u+cmJKPr0AnArDQbJlFJO7avMV6mS/GkZPNrzMBcoW3igMvmy7fNon1uPaz/Lkjp67b4wA7y/wohpfE3HRWQ06Z+MQlft8RUSFpp+9M0KXanQKBgQDeuCuxaGr9yUGOJC2WRgS3DA8lLGm+6BNh42ewBEwkH3jiJD/hrzOb9yFOBVa8ni2K6xUwDnofrpDCSqiYYPQuiyPTzpeCooYDK3/XyWUXTr+TroqPFEUM+K7XeWvW1EKKtEg1YuU4a7LnE7JpJ0LDoqtLtUuJ7O+Wa13HumlFJQKBgQDXk1ldxOJOXvvuEn3iez6BzOXtss6HRB5ppfhQ5rwiyfN5pUFTqvxPwa+GSEIuoeP3U+iZkvqGBM6InSB8ryfjfgL4t46CN5gcNWwA8NqGvZLYeCHpyBTtxo7OLQeGzzPaeRulwhelXwbCQAK1X+18RBsMR+SH9Knp0hS4+JJKlQKBgQDHxEu0ieMFlaIOO5cEJfOOt+tRvX9v87uG3rEfKQuejvgfZsJBzKMu7sBZueIttndFFkzf2OxjRHGlQ8/rNXNv1++fyLsPOnWXnEnEJGlfOYwOi8zOPzEcTGaO8OLwQ10YClKGSBkvvTIvn/Qz6zowPdUFSCzkHrhbpBvuzN4lXQKBgQCebDjGgkPVWFRH9urwH6Yl+ZAXiMnh+htnhILh4U7tOgBlqx5BAGz/p9T4F+4bGvnO7qkHA058Ytfs6ZvQRWBI/Hfuk+Z0p6pvQIsofdf6ISLjVhWGGnXW745O0iSv22G98jZxMBv0ecsbwbK7281I/zvpYIP/rbuYi7yS2omXqQKBgBxU1HVnu49AM+anFR+PZVo+BdYPtR2TfYIiFCdjuhOHMKHizHeSrwU/TusQBPJ5Y8MO+DXb1P5pUgmWbZklzmAzme+AVEO0TK1BRmZXpsLkTdX62CEbRC292cxWRw4kjpCmRvWWWU3TJg6mFmE8yKFB341cAGrKqx4x4vHa7CIq
//appId: 10****975
//cpId: 201d0780****82c352b5ba0e804
func main() {
	ts := strconv.FormatInt(int64(time.Now().UnixNano()/1000000)-int64(1000), 10)

	var queryMap = make(map[string]string)
	queryMap["method"] = "external.hms.gs.checkPlayerSign"
	queryMap["appId"] = "10***975"
	queryMap["cpId"] = "201d0780***0e804"
	queryMap["ts"] = ts
	queryMap["playerId"] = "100"
	queryMap["playerLevel"] = "10"
	queryMap["playerSSign"] = "VUOoWexHeQC98OFHyWapgKSACDwBgEHWb6IvPutKO0Z/wSVU3SDoK7********************BjaYFD03RWb2XBRKlnF7m455DeU2bvPZOsi7BhTDNPD0bTxY7PWlASLCSX7C7WqHN4/AWxDiU+ki2pPBstuSDecoUQQATBU35bQE2V7DtOsoGAhseuKXZe7yExMqszyZHLKaaqsbqq1rCua6FvJtwlwO82eY7N5kyW29r3MQ/uW1XGh4aPDods9UfD90BSLoPPmLjV9tREX/HFIdxkZ3FVWbkcWR4YQ=="

	var keyList []string
	for key, _ := range queryMap {
		if "cpSign" == key || "rtnSign" == key {
			continue
		}
		keyList = append(keyList, key)
	}
	sort.Strings(keyList)

	p := url.Values{}
	for _, key := range keyList {
		p.Add(key, queryMap[key])
	}

	pKey := "MIIEvgIBADANBgkqhkiG9wC7jN4ysPNBjQhmEsVla+OMOxvLvPJr61crxWt4eWak6wvr3*************************NCUm5m4XSGv3V8ox3t+qR4ECggEAG5EFQNSm8JHa8+va5UEiKDetOA1QHpmaGSj2EOE7urUKrId0HkNGMf9TeeXgGmY9DuVqnzpu3YDA/8agEPlsklnkxGSDtfedEF3inopkohg03EQX49lR07+VQbJXb88OH8Qn3Chol1CL1tlH6L4gjZ9C+AwdGpgyV8Sxvkjq59ckmABcF0C4qIKVZxpkITqyqD5cLrHw6D10bHv0y0l9dQDN0vUkt7jN8kvGZxJQ+u+cmJKPr0AnArDQbJlFJO7avMV6mS/GkZPNrzMBcoW3igMvmy7fNon1uPaz/Lkjp67b4wA7y/wohpfE3HRWQ06Z+MQlft8RUSFpp+9M0KXanQKBgQDeuCuxaGr9yUGOJC2WRgS3DA8lLGm+6BNh42ewBEwkH3jiJD/hrzOb9yFOBVa8ni2K6xUwDnofrpDCSqiYYPQuiyPTzpeCooYDK3/XyWUXTr+TroqPFEUM+K7XeWvW1EKKtEg1YuU4a7LnE7JpJ0LDoqtLtUuJ7O+Wa13HumlFJQKBgQDXk1ldxOJOXvvuEn3iez6BzOXtss6HRB5ppfhQ5rwiyfN5pUFTqvxPwa+GSEIuoeP3U+iZkvqGBM6InSB8ryfjfgL4t46CN5gcNWwA8NqGvZLYeCHpyBTtxo7OLQeGzzPaeRulwhelXwbCQAK1X+18RBsMR+SH9Knp0hS4+JJKlQKBgQDHxEu0ieMFlaIOO5cEJfOOt+tRvX9v87uG3rEfKQuejvgfZsJBzKMu7sBZueIttndFFkzf2OxjRHGlQ8/rNXNv1++fyLsPOnWXnEnEJGlfOYwOi8zOPzEcTGaO8OLwQ10YClKGSBkvvTIvn/Qz6zowPdUFSCzkHrhbpBvuzN4lXQKBgQCebDjGgkPVWFRH9urwH6Yl+ZAXiMnh+htnhILh4U7tOgBlqx5BAGz/p9T4F+4bGvnO7qkHA058Ytfs6ZvQRWBI/Hfuk+Z0p6pvQIsofdf6ISLjVhWGGnXW745O0iSv22G98jZxMBv0ecsbwbK7281I/zvpYIP/rbuYi7yS2omXqQKBgBxU1HVnu49AM+anFR+PZVo+BdYPtR2TfYIiFCdjuhOHMKHizHeSrwU/TusQBPJ5Y8MO+DXb1P5pUgmWbZklzmAzme+AVEO0TK1BRmZXpsLkTdX62CEbRC292cxWRw4kjpCmRvWWWU3TJg6mFmE8yKFB341cAGrKqx4x4vHa7CIq"
	privateKey, err := GetPrivateKey(pKey)
	if err != nil {
		fmt.Println("GetPrivateKey err: ", err)
		return
	}

	signature, err := SignSHA256WithRSA(p.Encode(), privateKey)
	if err != nil {
		fmt.Println("SignSHA256WithRSA err: ", err)
		return
	}

	p.Add("cpSign", base64.StdEncoding.EncodeToString(signature))


	resp, err := http.PostForm("https://gss-cn.game.hicloud.com/gameservice/api/gbClientApi", p)
	if err != nil {
		fmt.Println("Post err:", err)
		return
	}

	/*authUrl := fmt.Sprintf("https://gss-cn.game.hicloud.com/gameservice/api/gbClientApi?%s", p.Encode())
	fmt.Println("authUrl: ", authUrl)
	resp, err := http.Post(authUrl, "application/x-www-form-urlencoded;charset=UTF-8", bytes.NewBuffer([]byte("")))
	if err != nil {
		fmt.Println("Post err:", err)
		return
	}*/
	/*timeout := time.Duration(50 * time.Second) //超时时间50s
	client := &http.Client{
		Timeout: timeout,
	}
	request, err := http.NewRequest("POST", authUrl, strings.NewReader(""))
	if err != nil {
		fmt.Println("NewRequest err: ", err)
		return
	}
	//request.Header.Set("Content-Type", "text/html;charset=UTF-8")
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	resp, err := client.Do(request)*/

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ReadAll err: ", err)
		return
	}

	if http.StatusOK != resp.StatusCode {
		fmt.Println("StatusCode err code: ", resp.StatusCode)
		return
	}

	var data *DoAuthTokenRsp
	if err := json.Unmarshal(respBody, &data); err != nil {
		fmt.Println("Unmarshal err: ", err)
		return
	}

	fmt.Printf("data: %#v\n", data)

}

// 生成私钥
func GetPrivateKey(key string) (*rsa.PrivateKey, error) {
	privateKey := "-----BEGIN PRIVATE KEY-----\n" + key + "\n-----END PRIVATE KEY-----\n"
	block, _ := pem.Decode([]byte(privateKey))
	if nil == block {
		fmt.Println("GetPrivateKey error pem decode")
	}
	result, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if nil != err {
		fmt.Println("ParsePKCS8PrivateKey err: ", err)
		return nil, err
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

// sha256withRSA 签名
func SignSHA256WithRSA(data string, key *rsa.PrivateKey) ([]byte, error) {
	hash := sha256.Sum256([]byte(data))
	rng := rand.Reader
	signature, err := rsa.SignPKCS1v15(rng, key, crypto.SHA256, hash[:])
	if nil != err {
		fmt.Println("SignPKCS1v15 err: ", err)
		return nil, errors.New("error from signing: " + data)
	}
	return signature, nil
}
